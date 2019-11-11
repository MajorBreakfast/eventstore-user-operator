package controller

import (
	"github.com/MajorBreakfast/eventstore-user-operator/pkg/controller/eventstoreuser"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, eventstoreuser.Add)
}
