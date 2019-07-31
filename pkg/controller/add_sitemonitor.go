package controller

import (
	"wen/site-monitor-operator/pkg/controller/sitemonitor"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, sitemonitor.Add)
}
