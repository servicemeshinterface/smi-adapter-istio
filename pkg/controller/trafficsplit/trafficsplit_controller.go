package trafficsplit

import (
	"context"

	networkingv1alpha3 "github.com/deislabs/smi-adapter-istio/pkg/apis/networking/v1alpha3"
	smispecv1beta1 "github.com/deislabs/smi-adapter-istio/pkg/apis/smispec/v1beta1"

	"github.com/go-logr/logr"
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
)

var log = logf.Log.WithName("controller_trafficsplit")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

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
	err = c.Watch(&source.Kind{Type: &smispecv1beta1.TrafficSplit{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource VirtualService and DestinationRule
	err = c.Watch(&source.Kind{Type: &networkingv1alpha3.VirtualService{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &smispecv1beta1.TrafficSplit{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &networkingv1alpha3.DestinationRule{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &smispecv1beta1.TrafficSplit{},
	})
	if err != nil {
		return err
	}

	return nil
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
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileTrafficSplit) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling TrafficSplit")

	// Fetch the TrafficSplit instance
	trafficSplit := &smispecv1beta1.TrafficSplit{}
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

	reconcileResultVS, errVS := r.reconcileVirtualService(trafficSplit, reqLogger)
	reconcileResultDR, errDR := r.reconcileDestinationRule(trafficSplit, reqLogger)
	if errVS != nil {
		return reconcileResultVS, errVS
	}
	if errDR != nil {
		return reconcileResultDR, errDR
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileTrafficSplit) reconcileVirtualService(trafficSplit *smispecv1beta1.TrafficSplit,
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
		reqLogger.Error(err, "Failed to get VirtualService.")
		return reconcile.Result{}, err
	}

	// VS already exists - don't requeue
	reqLogger.Info("Skip reconcile: VirtualService already exists", "VirtualService.Namespace", found.Namespace,
		"VirtualService.Name", found.Name)
	return reconcile.Result{}, nil
}

func (r *ReconcileTrafficSplit) reconcileDestinationRule(trafficSplit *smispecv1beta1.TrafficSplit,
	reqLogger logr.Logger) (reconcile.Result, error) {
	// Define a new DestinationRule object
	dr := newDestinationRuleForCR(trafficSplit)

	// Set TrafficSplit instance as the owner and controller
	if err := controllerutil.SetControllerReference(trafficSplit, dr, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this DR already exists
	found := &networkingv1alpha3.DestinationRule{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: dr.Name, Namespace: dr.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new DestinationRule", "DestinationRule.Namespace", dr.Namespace,
			"DestinationRule.Name", dr.Name)
		err = r.client.Create(context.TODO(), dr)
		if err != nil {
			return reconcile.Result{}, err
		}

		// DestinationRule created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get DestinationRule.")
		return reconcile.Result{}, err
	}
	// DR already exists - don't requeue
	reqLogger.Info("Skip reconcile: DestinationRule already exists", "DestinationRule.Namespace", found.Namespace,
		"DestinationRule.Name", dr.Name)
	return reconcile.Result{}, nil
}

// newVSForCR returns a VirtualService with the same name/namespace as the cr
func newVSForCR(cr *smispecv1beta1.TrafficSplit) *networkingv1alpha3.VirtualService {
	labels := map[string]string{
		"traffic-split": cr.Name,
	}
	return &networkingv1alpha3.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-vs",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: networkingv1alpha3.VirtualServiceSpec{
			Hosts: []string{"myfoobarservice"},
			Http: []*networkingv1alpha3.HTTPRoute{{
				Route: []*networkingv1alpha3.HTTPRouteDestination{{
					Destination: &networkingv1alpha3.Destination{Host: "foo.com"},
					Weight:      42,
				},
				},
			}},
		},
	}
}

// newDestinationRuleForCR returns DestinationRule with the same name & namespace as of the
// Custom Resource
func newDestinationRuleForCR(cr *smispecv1beta1.TrafficSplit) *networkingv1alpha3.DestinationRule {
	labels := map[string]string{
		"traffic-split": cr.Name,
	}
	return &networkingv1alpha3.DestinationRule{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-dr",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: networkingv1alpha3.DestinationRuleSpec{
			Host: cr.Spec.Service,
			TrafficPolicy: &networkingv1alpha3.TrafficPolicy{
				Tls: &networkingv1alpha3.TLSSettings{
					Mode: "ISTIO_MUTUAL",
				},
			},
			Subsets: []*networkingv1alpha3.Subset{
				{
					Name: "v1",
					Labels: map[string]string{
						"version": "v1",
					},
				},
				{
					Name: "v2",
					Labels: map[string]string{
						"version": "v2",
					},
				},
			},
		},
	}
}
