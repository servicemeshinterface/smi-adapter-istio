package controller

import (
	"testing"

	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func TestBasic(t *testing.T) {
	called := false
	AddToManagerFuncs = []func(manager.Manager) error{
		func(m manager.Manager) error {
			called = true
			return nil
		},
	}

	testenv := &envtest.Environment{}
	cfg, err := testenv.Start()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		t.FailNow()
	}
	m, err := manager.New(cfg, manager.Options{
		NewClient: func(cache cache.Cache, config *rest.Config, options client.Options) (client.Client, error) {
			return nil, nil
		},
	})
	if m == nil {
		t.Errorf("Unexpected nil manager")
		t.FailNow()
	}

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	AddToManager(m)

	if !called {
		t.Errorf("Expected manager function to be called and it wasn't")
	}
}
