package e2e

import (
	goctx "context"
	"fmt"
	"testing"

	accessv1alpha1 "github.com/deislabs/smi-sdk-go/pkg/apis/access/v1alpha1"
	specsv1alpha1 "github.com/deislabs/smi-sdk-go/pkg/apis/specs/v1alpha1"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/deislabs/smi-adapter-istio/pkg/apis"
	rbacv1alpha1 "github.com/deislabs/smi-adapter-istio/pkg/apis/rbac/v1alpha1"
)

// deploy operator, create traffic target, verify servicerole and servicerolebind were created
func TestTrafficTarget(t *testing.T) {

	ttList := &accessv1alpha1.TrafficTarget{}
	err := framework.AddToFrameworkScheme(apis.AddToScheme, ttList)
	if err != nil {
		t.Fatalf("failed to add custom resource scheme to framework: %v", err)
	}

	// run subtests
	t.Run("traffictarget-group", func(t *testing.T) {
		t.Run("Cluster", TrafficTargetCluster)
	})

}

func TrafficTargetCluster(t *testing.T) {
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

	if err = TrafficTargetCreateTest(t, f, ctx); err != nil {
		t.Error(err)
	}

}

func TrafficTargetCreateTest(t *testing.T, f *framework.Framework, ctx *framework.TestCtx) error {

	namespace, err := ctx.GetNamespace()
	if err != nil {
		return fmt.Errorf("could not get namespace: %v", err)
	}
	ttName := "test-traffic-target"
	httpRouteGroup := &specsv1alpha1.HTTPRouteGroup{
		TypeMeta: metav1.TypeMeta{
			Kind:       "HTTPRouteGroup",
			APIVersion: "specs.smi-spec.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "the-routes",
			Namespace: namespace,
		},
		Matches: []specsv1alpha1.HTTPMatch{{
			Name:      "testGet",
			PathRegex: "/test",
			Methods:   []string{"Get"},
		}},
	}
	// create custom resource
	trafficTarget := &accessv1alpha1.TrafficTarget{
		TypeMeta: metav1.TypeMeta{
			Kind:       "TrafficTarget",
			APIVersion: "access.smi-spec.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      ttName,
			Namespace: namespace,
		},
		Destination: accessv1alpha1.IdentityBindingSubject{
			Kind:      "ServiceAccount",
			Name:      "service-a",
			Namespace: namespace,
			Port:      "8080",
		},
		Specs: []accessv1alpha1.TrafficTargetSpec{{
			Kind:    "HTTPRouteGroup",
			Name:    "the-routes",
			Matches: []string{"testGet"},
		}},
		Sources: []accessv1alpha1.IdentityBindingSubject{
			{
				Kind:      "ServiceAccount",
				Name:      "service-b",
				Namespace: namespace,
			},
		},
	}

	t.Log("Ensure servicerole and servicerolebinding not already running")
	serviceRoleBinding := &rbacv1alpha1.ServiceRoleBinding{}
	serviceRole := &rbacv1alpha1.ServiceRole{}
	err = f.Client.Get(goctx.TODO(), client.ObjectKey{Name: ttName, Namespace: namespace}, serviceRole)
	if !apierrors.IsNotFound(err) {
		return fmt.Errorf("Service role %s already exists in namespace: %s", ttName, namespace)
	}
	err = f.Client.Get(goctx.TODO(), client.ObjectKey{Name: ttName, Namespace: namespace}, serviceRoleBinding)
	if !apierrors.IsNotFound(err) {
		return fmt.Errorf("Service role binding %s already exists in namespace: %s", ttName, namespace)
	}

	// use TestCtx's create helper to create the specs and traffic target objects and
	//  add a cleanup function for the new object
	err = f.Client.Create(goctx.TODO(), httpRouteGroup, &framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		return err
	}
	err = f.Client.Create(goctx.TODO(), trafficTarget, &framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		return err
	}

	waitErr := wait.Poll(retryInterval, timeout, func() (bool, error) {
		err := f.Client.Get(goctx.TODO(), client.ObjectKey{Name: ttName, Namespace: namespace}, serviceRole)
		if err != nil && apierrors.IsNotFound(err) {
			return false, nil
		} else if err != nil {
			return false, err
		}
		return true, nil
	})
	if waitErr != nil {
		return fmt.Errorf("Expected service role '%s' to exist in namespace '%s' but does not exist: %s", ttName, namespace, waitErr)

	}
	waitErr = wait.Poll(retryInterval, timeout, func() (bool, error) {
		err = f.Client.Get(goctx.TODO(), client.ObjectKey{Name: ttName, Namespace: namespace}, serviceRoleBinding)
		if err != nil && apierrors.IsNotFound(err) {
			return false, nil
		} else if err != nil {
			return false, err
		}
		return true, nil
	})
	if waitErr != nil {
		return fmt.Errorf("Expected service role binding '%s' to exist in namespace '%s' but does not exist: %s", ttName, namespace, waitErr)

	}

	return nil
}
