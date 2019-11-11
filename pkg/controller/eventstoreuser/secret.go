package eventstoreuser

import (
	"context"

	josefbrandlv1 "github.com/MajorBreakfast/eventstore-user-operator/pkg/apis/josefbrandl/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func doesSecretExist(
	r *ReconcileEventStoreUser,
	esUserCR *josefbrandlv1.EventStoreUser,
) (bool, error) {
	searchKey := types.NamespacedName{
		Name:      esUserCR.Name + "-eventstore-user",
		Namespace: esUserCR.Namespace,
	}

	if err := r.client.Get(context.TODO(), searchKey, &corev1.Secret{}); err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func createOrUpdateSecret(
	r *ReconcileEventStoreUser,
	esUserCR *josefbrandlv1.EventStoreUser,
	password string,
) error {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      esUserCR.Name + "-eventstore-user",
			Namespace: esUserCR.Namespace,
		},
		Data: map[string][]byte{
			"username": []byte(esUserCR.Name),
			"password": []byte(password),
		},
	}
	if err := controllerutil.SetControllerReference(esUserCR, secret, r.scheme); err != nil {
		return err
	}

	exists, err := doesSecretExist(r, esUserCR)
	if err != nil {
		return err
	}

	if exists {
		if err := r.client.Update(context.TODO(), secret); err != nil {
			return err
		}
	} else {
		if err := r.client.Create(context.TODO(), secret); err != nil {
			return err
		}
	}

	return nil
}
