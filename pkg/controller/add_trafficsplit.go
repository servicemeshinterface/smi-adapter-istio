package controller

import (
	"github.com/deislabs/smi-adapter-istio/pkg/controller/trafficsplit"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, trafficsplit.Add)
}
