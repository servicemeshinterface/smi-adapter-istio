package trafficsplit

import (
	"context"
	"testing"

	splitv1alpha2 "github.com/deislabs/smi-sdk-go/pkg/apis/split/v1alpha2"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	apitypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"

	networkingv1alpha3 "github.com/deislabs/smi-adapter-istio/pkg/apis/networking/v1alpha3"
)

func TestNewReconciler(t *testing.T) {
	mgr := FakeManager{}
	r := newReconciler(mgr)
	var _ reconcile.Reconciler = r // test r is reconcile.Reconciler
}

func TestReconcile_ErrorIsNotFound(t *testing.T) {
	trafficSplitObj := &splitv1alpha2.TrafficSplit{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "traffic-split-name",
			Namespace: "default",
		},
		Spec: splitv1alpha2.TrafficSplitSpec{},
	}
	objs := []runtime.Object{}
	cl := fake.NewFakeClient(objs...)
	s := scheme.Scheme
	s.AddKnownTypes(splitv1alpha2.SchemeGroupVersion, trafficSplitObj)
	reconcileTrafficSplit := &ReconcileTrafficSplit{client: cl, scheme: s}
	req := reconcile.Request{NamespacedName: apitypes.NamespacedName{
		Namespace: "default",
		Name:      "traffic-split-name"},
	}
	_, err := reconcileTrafficSplit.Reconcile(req)
	if err != nil {
		t.Errorf("Expected no err, got %v", err)
	}
}

func TestReconcile(t *testing.T) {
	trafficSplitObj := &splitv1alpha2.TrafficSplit{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "traffic-split-name",
			Namespace: "default",
		},
		Spec: splitv1alpha2.TrafficSplitSpec{},
	}
	virtualServiceObj := &networkingv1alpha3.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "traffic-split-name-vs",
			Namespace: "default",
			Labels:    map[string]string{"traffic-split": "traffic-split-name"},
		},
		Spec: networkingv1alpha3.VirtualServiceSpec{},
	}
	// Register the object in the fake client.
	objs := []runtime.Object{
		trafficSplitObj,
	}
	s := scheme.Scheme
	s.AddKnownTypes(splitv1alpha2.SchemeGroupVersion, trafficSplitObj)
	s.AddKnownTypes(networkingv1alpha3.SchemeGroupVersion, virtualServiceObj)

	cl := fake.NewFakeClient(objs...)
	err := cl.Get(context.TODO(), apitypes.NamespacedName{
		Namespace: "default",
		Name:      "traffic-split-name-vs"},
		virtualServiceObj)
	if !apierrors.IsNotFound(err) {
		t.Fatalf("expected virtual service not to exist, got err: %s", err)
	}

	reconcileTrafficSplit := &ReconcileTrafficSplit{client: cl, scheme: s}
	req := reconcile.Request{NamespacedName: apitypes.NamespacedName{
		Namespace: "default",
		Name:      "traffic-split-name"},
	}

	_, err = reconcileTrafficSplit.Reconcile(req)
	if err != nil {
		t.Errorf("Expected no err, got %s", err)
	}

	err = cl.Get(context.TODO(), apitypes.NamespacedName{
		Namespace: "default",
		Name:      "traffic-split-name-vs"},
		virtualServiceObj)
	if err != nil {
		t.Errorf("Expected virtual service object to be created successfully, but was not: %s", err)
	}
}

func TestWeightToPercent(t *testing.T) {
	backends := []splitv1alpha2.TrafficSplitBackend{
		splitv1alpha2.TrafficSplitBackend{Service: "a", Weight: 1000},
		splitv1alpha2.TrafficSplitBackend{Service: "b", Weight: 2000},
	}
	weights := weightToPercent(backends)
	if weights["a"] != 33 {
		t.Errorf("Expected Service a to have percent weight of 33 but got %v", weights["a"])
	}
	if weights["b"] != 67 {
		t.Errorf("Expected Service b to have percent weight of 67 but got %v", weights["b"])
	}

	backends = []splitv1alpha2.TrafficSplitBackend{
		splitv1alpha2.TrafficSplitBackend{Service: "a", Weight: 1000},
		splitv1alpha2.TrafficSplitBackend{Service: "b", Weight: 1000},
		splitv1alpha2.TrafficSplitBackend{Service: "c", Weight: 1000},
	}
	weights = weightToPercent(backends)
	if weights["a"] != 33 {
		t.Errorf("Expected Service a to have percent weight of 33 but got %v", weights["a"])
	}
	if weights["b"] != 33 {
		t.Errorf("Expected Service b to have percent weight of 33 but got %v", weights["b"])
	}
	if weights["c"] != 34 {
		t.Errorf("Expected Service b to have percent weight of 34 but got %v", weights["c"])
	}

	backends = []splitv1alpha2.TrafficSplitBackend{
		splitv1alpha2.TrafficSplitBackend{Service: "a", Weight: 20},
		splitv1alpha2.TrafficSplitBackend{Service: "b", Weight: 30},
		splitv1alpha2.TrafficSplitBackend{Service: "c", Weight: 50},
	}
	weights = weightToPercent(backends)
	if weights["a"] != 20 {
		t.Errorf("Expected Service a to have percent weight of 20 but got %v", weights["a"])
	}
	if weights["b"] != 30 {
		t.Errorf("Expected Service b to have percent weight of 30 but got %v", weights["b"])
	}
	if weights["c"] != 50 {
		t.Errorf("Expected Service b to have percent weight of 50 but got %v", weights["c"])
	}

	backends = []splitv1alpha2.TrafficSplitBackend{
		splitv1alpha2.TrafficSplitBackend{Service: "a", Weight: 5},
		splitv1alpha2.TrafficSplitBackend{Service: "b", Weight: 10},
		splitv1alpha2.TrafficSplitBackend{Service: "c", Weight: 20},
	}
	weights = weightToPercent(backends)
	if weights["a"] != 14 {
		t.Errorf("Expected Service a to have percent weight of 14 but got %v", weights["a"])
	}
	if weights["b"] != 29 {
		t.Errorf("Expected Service b to have percent weight of 29 but got %v", weights["b"])
	}
	if weights["c"] != 57 {
		t.Errorf("Expected Service b to have percent weight of 57 but got %v", weights["c"])
	}
}

type FakeManager struct{}

func (fm FakeManager) Add(manager.Runnable) error                   { return nil }
func (fm FakeManager) SetFields(interface{}) error                  { return nil }
func (fm FakeManager) Start(<-chan struct{}) error                  { return nil }
func (fm FakeManager) GetConfig() *rest.Config                      { return &rest.Config{} }
func (fm FakeManager) GetScheme() *runtime.Scheme                   { return &runtime.Scheme{} }
func (fm FakeManager) GetAdmissionDecoder() types.Decoder           { return nil }
func (fm FakeManager) GetClient() client.Client                     { return nil }
func (fm FakeManager) GetFieldIndexer() client.FieldIndexer         { return nil }
func (fm FakeManager) GetCache() cache.Cache                        { return nil }
func (fm FakeManager) GetRecorder(name string) record.EventRecorder { return nil }
func (fm FakeManager) GetRESTMapper() meta.RESTMapper               { return nil }
func (fm FakeManager) GetAPIReader() client.Reader                  { return nil }
func (fm FakeManager) GetWebhookServer() *webhook.Server            { return &webhook.Server{} }
