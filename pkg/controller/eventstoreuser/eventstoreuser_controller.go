package eventstoreuser

import (
	"context"

	josefbrandlv1 "github.com/MajorBreakfast/eventstore-user-operator/pkg/apis/josefbrandl/v1"
	"github.com/MajorBreakfast/eventstore-user-operator/pkg/controller/eventstoreuser/config"
	"github.com/MajorBreakfast/eventstore-user-operator/pkg/controller/eventstoreuser/eventstore"
	"github.com/dchest/uniuri"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/util/slice"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_eventstoreuser")

const eventStoreUserFinalizer = "finalizer.eventstoreuser.josefbrandl.com"

// Add creates a new EventStoreUser Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileEventStoreUser{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("eventstoreuser-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource EventStoreUser
	err = c.Watch(&source.Kind{Type: &josefbrandlv1.EventStoreUser{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resources
	err = c.Watch(&source.Kind{Type: &corev1.Secret{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &josefbrandlv1.EventStoreUser{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileEventStoreUser implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileEventStoreUser{}

// ReconcileEventStoreUser reconciles a EventStoreUser object
type ReconcileEventStoreUser struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a EventStoreUser object and makes changes based on the state read
// and what is in the EventStoreUser.Spec
func (r *ReconcileEventStoreUser) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues(
		"Request.Namespace", request.Namespace,
		"Request.Name", request.Name,
	)
	reqLogger.Info("Reconciling EventStoreUser")

	// Fetch the EventStoreUser custom resource
	esUserCR := &josefbrandlv1.EventStoreUser{}
	err := r.client.Get(context.TODO(), request.NamespacedName, esUserCR)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted. Return and don't requeue
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// Add the Event Store name to the log messages
	reqLogger = reqLogger.WithValues("eventStore", esUserCR.Spec.EventStore)

	desiredUser := eventstore.UserWithPassword{
		User: eventstore.User{
			LoginName: request.Name,
			FullName:  request.Name,
			Groups:    esUserCR.Spec.Groups,
		},
		Password: "",
	}

	log.Info("Loading connection information for the specified Event Store")
	esConnOptions, err := config.LoadEventStoreConnectionOptions(esUserCR.Spec.EventStore)
	if err != nil {
		return reconcile.Result{}, err
	}
	es := eventstore.NewConnection(esConnOptions)

	// Finalizer setup
	isMarkedToBeDeleted := esUserCR.GetDeletionTimestamp() != nil
	if isMarkedToBeDeleted {
		if hasFinalizer(esUserCR) {
			if err := es.DeleteUser(desiredUser.LoginName); err != nil {
				return reconcile.Result{}, err
			}
			if err := r.removeFinalizer(reqLogger, esUserCR); err != nil {
				return reconcile.Result{}, err
			}
		}
		return reconcile.Result{}, nil // Finalized: Return and don't requeue
	}
	if err := r.addFinalizerIfNeeded(reqLogger, esUserCR); err != nil {
		return reconcile.Result{}, err
	}

	log.Info("Checking for existing Event Store user")
	userExistsOnEventStore := true
	existingUser, err := es.GetUser(desiredUser.LoginName)
	if err != nil {
		if err == eventstore.ErrUserNotFound {
			userExistsOnEventStore = false
		} else {
			return reconcile.Result{}, err
		}
	}

	log.Info("Checking whether secret exists")
	secretExists, err := doesSecretExist(r, esUserCR)
	if err != nil {
		return reconcile.Result{}, err
	}
	userGetsNewPassword := !secretExists

	if userGetsNewPassword {
		log.Info("Generating new Event Store user password")
		desiredUser.Password = uniuri.NewLen(50)
	}

	if userExistsOnEventStore {
		if existingUser.IsEqual(desiredUser.User) {
			log.Info("User already exists on the Event Store with correct fields")
		} else {
			log.Info("Updating user on the Event Store")
			if err := es.UpdateUser(desiredUser.User); err != nil {
				return reconcile.Result{}, err
			}
		}

		if userGetsNewPassword {
			log.Info("Resetting password of user on the Event Store")
			if err := es.ResetPassword(desiredUser.LoginName, desiredUser.Password); err != nil {
				return reconcile.Result{}, err
			}
		} else {
			log.Info("Keep the current password of the user on the Event Store")
		}
	} else {
		log.Info("Creating new user on the Event Store")
		if err := es.CreateUser(desiredUser); err != nil {
			return reconcile.Result{}, err
		}
	}

	if userGetsNewPassword {
		log.Info("Creating/updating secret")
		err := createOrUpdateSecret(r, esUserCR, desiredUser.Password)
		if err != nil {
			return reconcile.Result{}, err
		}
	} else {
		log.Info("Keep the current secret")
	}

	reqLogger.Info("Finished reconcile")
	return reconcile.Result{}, nil
}

func hasFinalizer(esUserCR *josefbrandlv1.EventStoreUser) bool {
	return slice.ContainsString(esUserCR.GetFinalizers(), eventStoreUserFinalizer, nil)
}

func (r *ReconcileEventStoreUser) addFinalizerIfNeeded(reqLogger logr.Logger, esUserCR *josefbrandlv1.EventStoreUser) error {
	if hasFinalizer(esUserCR) {
		return nil
	}

	reqLogger.Info("Adding Finalizer to the EventStoreUser custom resource")
	esUserCR.SetFinalizers(append(esUserCR.GetFinalizers(), eventStoreUserFinalizer))

	err := r.client.Update(context.TODO(), esUserCR)
	if err != nil {
		return err
	}

	return nil
}

func (r *ReconcileEventStoreUser) removeFinalizer(reqLogger logr.Logger, esUserCR *josefbrandlv1.EventStoreUser) error {
	esUserCR.SetFinalizers(slice.RemoveString(esUserCR.GetFinalizers(), eventStoreUserFinalizer, nil))
	return r.client.Update(context.TODO(), esUserCR)
}
