package trafficsplit

import (
	"context"
	"encoding/json"
	"math"

	splitv1alpha2 "github.com/deislabs/smi-sdk-go/pkg/apis/split/v1alpha2"
	"github.com/go-logr/logr"
	"github.com/google/go-cmp/cmp"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	networkingv1alpha3 "github.com/deislabs/smi-adapter-istio/pkg/apis/networking/v1alpha3"
)

var log = logf.Log.WithName("controller_trafficsplit")

// Add creates a new TrafficSplit Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileTrafficSplit{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("trafficsplit-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource TrafficSplit
	err = c.Watch(&source.Kind{Type: &splitv1alpha2.TrafficSplit{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource VirtualService
	return c.Watch(&source.Kind{Type: &networkingv1alpha3.VirtualService{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &splitv1alpha2.TrafficSplit{},
	})
}

var _ reconcile.Reconciler = &ReconcileTrafficSplit{}

// ReconcileTrafficSplit reconciles a TrafficSplit object
type ReconcileTrafficSplit struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a TrafficSplit object and makes changes based on the state read
// and what is in the TrafficSplit.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileTrafficSplit) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling TrafficSplit")

	// Fetch the TrafficSplit instance
	trafficSplit := &splitv1alpha2.TrafficSplit{}
	err := r.client.Get(context.TODO(), request.NamespacedName, trafficSplit)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			reqLogger.Info("TrafficSplit object not found.")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		reqLogger.Error(err, "Failed to get TrafficSplit. Request will be requeued.")
		return reconcile.Result{}, err
	}

	return r.reconcileVirtualService(trafficSplit, reqLogger)
}

func (r *ReconcileTrafficSplit) reconcileVirtualService(trafficSplit *splitv1alpha2.TrafficSplit,
	reqLogger logr.Logger) (reconcile.Result, error) {
	// Define a new VirtualService object
	vs := newVSForCR(trafficSplit)

	// Set TrafficSplit instance as the owner and controller
	if err := controllerutil.SetControllerReference(trafficSplit, vs, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this VS already exists
	found := &networkingv1alpha3.VirtualService{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: vs.Name, Namespace: vs.Namespace}, found)

	// Create VS
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new VirtualService", "VirtualService.Namespace", vs.Namespace,
			"VirtualService.Name", vs.Name)
		err = r.client.Create(context.TODO(), vs)
		if err != nil {
			return reconcile.Result{}, err
		}

		// VirtualService created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get VirtualService.", "VirtualService.Namespace", vs.Namespace,
			"VirtualService.Name", vs.Name)
		return reconcile.Result{}, err
	}

	// Update VS
	if diff := cmp.Diff(vs.Spec, found.Spec); diff != "" {
		reqLogger.Info("Updating VirtualService", "VirtualService.Namespace", vs.Namespace,
			"VirtualService.Name", vs.Name)
		clone := found.DeepCopy()
		clone.Spec = vs.Spec
		err = r.client.Update(context.TODO(), clone)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

//weightToPercent takes TrafficSplitBackends and maps a service to a weight
//  in the form of an integer from 1 - 100. The sum of weights must be 100. The last service
//  will be adjusted so that the total weight of the map returned at the end is equal to 100.
func weightToPercent(backends []splitv1alpha2.TrafficSplitBackend) map[string]int {
	weights := map[string]int{}
	totalWeight := 0
	for _, b := range backends {
		totalWeight = totalWeight + int(b.Weight)
		weights[b.Service] = 0
	}

	if totalWeight == 0 {
		return weights
	}

	totalPercent := 0
	for i, b := range backends {
		w := (float64(b.Weight) / float64(totalWeight)) * 100
		per := math.Round(float64(w))
		percent := int(per)
		if i == len(backends)-1 {
			percent = 100 - totalPercent
		}
		weights[b.Service] = percent
		totalPercent = totalPercent + percent
	}
	return weights
}

// newVSForCR returns a VirtualService with the same name/namespace as the cr
func newVSForCR(cr *splitv1alpha2.TrafficSplit) *networkingv1alpha3.VirtualService {
	labels := map[string]string{
		"traffic-split": cr.Name,
	}

	var backends []*networkingv1alpha3.HTTPRouteDestination
	weights := weightToPercent(cr.Spec.Backends)

	for _, b := range cr.Spec.Backends {
		r := &networkingv1alpha3.HTTPRouteDestination{
			Destination: &networkingv1alpha3.Destination{Host: b.Service},
			Weight:      int32(weights[b.Service]),
		}

		backends = append(backends, r)
	}

	gatewaysStr := cr.ObjectMeta.Annotations["VirtualService.v1alpha3.networking.istio.io/spec.gateways"]
	var gateways []string
	_ = json.Unmarshal([]byte(gatewaysStr), &gateways)

	return &networkingv1alpha3.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-vs",
			Namespace: cr.Namespace,
			Labels:    labels,
		},

		Spec: networkingv1alpha3.VirtualServiceSpec{
			Hosts:    []string{cr.Spec.Service},
			Gateways: gateways,

			Http: []*networkingv1alpha3.HTTPRoute{
				{
					Route: backends,
				},
			},
		},
	}
}
