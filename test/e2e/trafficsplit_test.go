package e2e

import (
	goctx "context"
	"fmt"
	"testing"
	"time"

	splitv1alpha2 "github.com/deislabs/smi-sdk-go/pkg/apis/split/v1alpha2"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/deislabs/smi-adapter-istio/pkg/apis"
	networkingv1alpha3 "github.com/deislabs/smi-adapter-istio/pkg/apis/networking/v1alpha3"
)

var (
	retryInterval        = time.Second * 5
	timeout              = time.Second * 60
	cleanupRetryInterval = time.Second * 1
	cleanupTimeout       = time.Second * 5
)

// deploy operator, create traffic split, verify virtualservice was created
func TestTrafficSplit(t *testing.T) {

	tsList := &splitv1alpha2.TrafficSplitList{}
	err := framework.AddToFrameworkScheme(apis.AddToScheme, tsList)
	if err != nil {
		t.Fatalf("failed to add custom resource scheme to framework: %v", err)
	}

	// run subtests
	t.Run("trafficsplit-group", func(t *testing.T) {
		t.Run("Cluster", TrafficSplitCluster)
	})

}

func TrafficSplitCluster(t *testing.T) {
	t.Parallel()
	ctx := framework.NewTestCtx(t)
	defer ctx.Cleanup()

	err := ctx.InitializeClusterResources(&framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		t.Fatalf("failed to initialize cluster resources: %v", err)
	}

	t.Log("Initialized cluster resources")

	namespace, err := ctx.GetNamespace()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("namespace: " + namespace)
	// get global framework variables
	f := framework.Global

	// wait for smi-istio-adapter operator to be ready
	err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, "smi-adapter-istio", 1, retryInterval, timeout)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("wait completed")

	if err = TrafficSplitCreateTest(t, f, ctx); err != nil {
		t.Error(err)
	}

}

func TrafficSplitCreateTest(t *testing.T, f *framework.Framework, ctx *framework.TestCtx) error {

	namespace, err := ctx.GetNamespace()
	if err != nil {
		return fmt.Errorf("could not get namespace: %v", err)
	}

	// create custom resource
	trafficSplit := &splitv1alpha2.TrafficSplit{
		TypeMeta: metav1.TypeMeta{
			Kind:       "TrafficSplit",
			APIVersion: "split.smi-spec.io/v1alpha2",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-service-rollout",
			Namespace: namespace,
		},
		Spec: splitv1alpha2.TrafficSplitSpec{
			Service: "test-service",
			Backends: []splitv1alpha2.TrafficSplitBackend{
				splitv1alpha2.TrafficSplitBackend{
					Service: "test-service-v1",
					Weight:  1,
				},
				splitv1alpha2.TrafficSplitBackend{
					Service: "test-service-v2",
					Weight:  0,
				},
			},
		},
	}

	t.Log("Ensure virtual service is not already running")
	virtualService := &networkingv1alpha3.VirtualService{}
	vsName := "test-service-rollout-vs"
	err = f.Client.Get(goctx.TODO(), client.ObjectKey{Name: vsName, Namespace: namespace}, virtualService)
	if !apierrors.IsNotFound(err) {
		return fmt.Errorf("virtual service %s already exists in namespace: %s", vsName, namespace)
	}

	// use TestCtx's create helper to create the traffic split object and
	//  add a cleanup function for the new object
	err = f.Client.Create(goctx.TODO(), trafficSplit, &framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		return err
	}

	waitErr := wait.Poll(retryInterval, timeout, func() (bool, error) {
		err := f.Client.Get(goctx.TODO(), client.ObjectKey{Name: vsName, Namespace: namespace}, virtualService)
		if err != nil && apierrors.IsNotFound(err) {
			return false, nil
		} else if err != nil {
			return false, err
		}
		return true, nil
	})
	if waitErr != nil {
		return fmt.Errorf("Expected virtual service '%s' to exist in namespace '%s' but does not exist: %s", vsName, namespace, waitErr)

	}

	return nil
}
