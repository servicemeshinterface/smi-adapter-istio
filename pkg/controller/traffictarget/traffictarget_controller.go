package traffictarget

import (
	"context"
	"fmt"

	accessv1alpha1 "github.com/deislabs/smi-sdk-go/pkg/apis/access/v1alpha1"
	specsv1alpha1 "github.com/deislabs/smi-sdk-go/pkg/apis/specs/v1alpha1"
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

	rbacv1alpha1 "github.com/deislabs/smi-adapter-istio/pkg/apis/rbac/v1alpha1"
)

var log = logf.Log.WithName("controller_traffictarget")

// Add creates a new TrafficTarget Controller and adds it to the Manager.
// The Manager will set fields on the Controller and Start it when the Manager
// is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileTrafficTarget{client: mgr.GetClient(),
		scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("traffictarget-controller", mgr,
		controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource TrafficTarget
	err = c.Watch(&source.Kind{Type: &accessv1alpha1.TrafficTarget{}},
		&handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource ServiceRoleBindings and requeue
	// the owner TrafficTarget
	return c.Watch(&source.Kind{Type: &rbacv1alpha1.ServiceRoleBinding{}},
		&handler.EnqueueRequestForOwner{
			IsController: true,
			OwnerType:    &accessv1alpha1.TrafficTarget{},
		})
}

var _ reconcile.Reconciler = &ReconcileTrafficTarget{}

// ReconcileTrafficTarget reconciles a TrafficTarget object
type ReconcileTrafficTarget struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a TrafficTarget object and
// makes changes based on the state read and what is in the TrafficTarget.Spec
//
// The Controller will requeue the Request to be processed again if the returned
// error is non-nil or Result.Requeue is true, otherwise upon completion it will
// remove the work from the queue.
func (r *ReconcileTrafficTarget) Reconcile(
	request reconcile.Request,
) (reconcile.Result, error) {
	reqLogger := log.WithValues(
		"Request.Namespace", request.Namespace,
		"Request.Name", request.Name)
	reqLogger.Info("Reconciling TrafficTarget")

	// Fetch the TrafficTarget instance
	trafficTarget := &accessv1alpha1.TrafficTarget{}
	err := r.client.Get(context.TODO(), request.NamespacedName, trafficTarget)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile
			// request. Owned objects are automatically garbage collected. For
			// additional cleanup logic use finalizers. Return and don't requeue
			reqLogger.Info("TrafficTarget object not found.")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		reqLogger.Error(err,
			"Failed to get TrafficTarget. Request will be requeued.")
		return reconcile.Result{}, err
	}

	return r.reconcileTrafficTarget(trafficTarget, reqLogger)
}

func (r *ReconcileTrafficTarget) reconcileTrafficTarget(trafficTarget *accessv1alpha1.TrafficTarget, reqLogger logr.Logger) (reconcile.Result, error) {

	svcRole, svcRoleBinding, err := r.createRBAC(trafficTarget)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Set TrafficTarget instance as the owner and controller
	if err := controllerutil.SetControllerReference(trafficTarget,
		svcRole, r.scheme); err != nil {
		return reconcile.Result{}, err
	}
	if err := controllerutil.SetControllerReference(trafficTarget,
		svcRoleBinding, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	recSvcRole, errSvcRole := r.createServiceRole(svcRole, reqLogger)
	recSvcRoleBinding, errSvcRoleBinding := r.createServiceRoleBinding(
		svcRoleBinding, reqLogger,
	)
	if errSvcRole != nil {
		return recSvcRole, errSvcRole
	}
	if errSvcRoleBinding != nil {
		return recSvcRoleBinding, errSvcRoleBinding
	}
	return reconcile.Result{}, nil
}

func (r *ReconcileTrafficTarget) createServiceRole(
	svcRole *rbacv1alpha1.ServiceRole,
	reqLogger logr.Logger,
) (reconcile.Result, error) {
	// Check if this ServiceRole already exists
	foundSvcRole := &rbacv1alpha1.ServiceRole{}
	err := r.client.Get(
		context.TODO(),
		types.NamespacedName{Name: svcRole.Name, Namespace: svcRole.Namespace},
		foundSvcRole,
	)

	// Create ServiceRole
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new ServiceRole",
			"ServiceRole.Namespace", svcRole.Namespace,
			"ServiceRole.Name", svcRole.Name)
		err = r.client.Create(context.TODO(), svcRole)
		if err != nil {
			return reconcile.Result{}, err
		}

		// ServiceRole created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get ServiceRole",
			"ServiceRole.Namespace", svcRole.Namespace,
			"ServiceRole.Name", svcRole.Name)
		return reconcile.Result{}, err
	}

	// Update ServiceRole
	if diff := cmp.Diff(svcRole.Spec, foundSvcRole.Spec); diff != "" {
		reqLogger.Info("Updating ServiceRole",
			"ServiceRole.Namespace", svcRole.Namespace,
			"ServiceRole.Name", svcRole.Name)
		clone := foundSvcRole.DeepCopy()
		clone.Spec = svcRole.Spec
		err = r.client.Update(context.TODO(), clone)
		if err != nil {
			return reconcile.Result{}, err
		}
	}
	return reconcile.Result{}, nil
}

func (r *ReconcileTrafficTarget) createServiceRoleBinding(
	svcRoleBinding *rbacv1alpha1.ServiceRoleBinding,
	reqLogger logr.Logger,
) (reconcile.Result, error) {
	// Check if this ServiceRoleBinding already exists
	foundSvcRoleBinding := &rbacv1alpha1.ServiceRoleBinding{}
	err := r.client.Get(
		context.TODO(),
		types.NamespacedName{
			Name: svcRoleBinding.Name, Namespace: svcRoleBinding.Namespace,
		},
		foundSvcRoleBinding,
	)

	// Create ServiceRoleBinding
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new ServiceRoleBinding",
			"ServiceRoleBinding.Namespace", svcRoleBinding.Namespace,
			"ServiceRoleBinding.Name", svcRoleBinding.Name)
		err = r.client.Create(context.TODO(), svcRoleBinding)
		if err != nil {
			return reconcile.Result{}, err
		}

		// ServiceRoleBinding created successfully - don't requeue
		return reconcile.Result{}, err
	} else if err != nil {
		reqLogger.Error(err, "Failed to get ServiceRoleBinding",
			"ServiceRoleBinding.Namespace", svcRoleBinding.Namespace,
			"ServiceRoleBinding.Name", svcRoleBinding.Name)
		return reconcile.Result{}, err
	}

	// Update ServiceRoleBinding
	if diff := cmp.Diff(
		svcRoleBinding.Spec, foundSvcRoleBinding.Spec,
	); diff != "" {
		reqLogger.Info("Updating ServiceRoleBinding",
			"ServiceRoleBinding.Namespace", svcRoleBinding.Namespace,
			"ServiceRoleBinding.Name", svcRoleBinding.Name)
		clone := foundSvcRoleBinding.DeepCopy()
		clone.Spec = svcRoleBinding.Spec
		err = r.client.Update(context.TODO(), clone)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

// createRBAC creates a ServiceRole and ServiceRoleBinding for each
// TrafficTarget. For all the HTTPRouteGroup objects referred in the
// TrafficTarget will also be queried.
func (r *ReconcileTrafficTarget) createRBAC(trafficTarget *accessv1alpha1.TrafficTarget) (*rbacv1alpha1.ServiceRole, *rbacv1alpha1.ServiceRoleBinding, error) {
	var subjects []*rbacv1alpha1.Subject
	for _, src := range trafficTarget.Sources {
		// TODO:
		// Remove the hardcoded value of `cluster.local`
		subjects = append(subjects, &rbacv1alpha1.Subject{
			User: fmt.Sprintf("cluster.local/ns/%s/sa/%s", src.Namespace, src.Name),
		})
	}
	// same set of constraints generated from a TrafficTarget apply to all the
	// AccessRules so build them early itself
	constraints := getConstraints(trafficTarget.Destination)

	var rules []*rbacv1alpha1.AccessRule
	for _, spec := range trafficTarget.Specs {
		matches, err := r.findMatches(spec, trafficTarget.Namespace)
		if err != nil {
			return nil, nil, err
		}
		for _, match := range matches {
			rules = append(rules, &rbacv1alpha1.AccessRule{
				// Apply to all the services, hardcoded because the
				// authorization of traffic is done at constraints level
				Services:    []string{"*"},
				Methods:     match.Methods,
				Paths:       []string{match.PathRegex},
				Constraints: constraints,
			})
		}
	}

	svcRole := &rbacv1alpha1.ServiceRole{
		ObjectMeta: metav1.ObjectMeta{
			Name:      trafficTarget.Name,
			Namespace: trafficTarget.Namespace,
		},
		Spec: rbacv1alpha1.ServiceRoleSpec{
			Rules: rules,
		},
	}
	svcRB := &rbacv1alpha1.ServiceRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      trafficTarget.Name,
			Namespace: trafficTarget.Namespace,
		},
		Spec: rbacv1alpha1.ServiceRoleBindingSpec{
			Subjects: subjects,
			RoleRef: &rbacv1alpha1.RoleRef{
				Kind: "ServiceRole",
				Name: trafficTarget.Name,
			},
		},
	}

	return svcRole, svcRB, nil
}

// getConstraints reads the information from the destination object of
// TrafficTarget and creates correspoding constraints to be applied in the
// Istio ServiceRole object
func getConstraints(
	dst accessv1alpha1.IdentityBindingSubject,
) []*rbacv1alpha1.AccessRule_Constraint {
	constraints := []*rbacv1alpha1.AccessRule_Constraint{
		{Key: "destination.user", Values: []string{dst.Name}},
		{Key: "destination.namespace", Values: []string{dst.Namespace}},
	}
	if dst.Port != "" {
		constraints = append(constraints, &rbacv1alpha1.AccessRule_Constraint{
			Key:    "destination.port",
			Values: []string{dst.Port},
		})
	}
	return constraints
}

// findMatches finds and returns a "match" object for a HTTPRouteGroup referred
// in the TrafficTarget, if the referred HTTPRouteGroup does not exists then it
// returns an error. If the specific match name referred in the
// `TrafficTarget.specs.matches` does not exist in the given HTTPRouteGroup
// object, then it returns an error.
func (r *ReconcileTrafficTarget) findMatches(
	spec accessv1alpha1.TrafficTargetSpec, ns string,
) ([]specsv1alpha1.HTTPMatch, error) {
	httpRouteGroup := &specsv1alpha1.HTTPRouteGroup{}
	if err := r.client.Get(
		context.TODO(),
		types.NamespacedName{Namespace: ns, Name: spec.Name},
		httpRouteGroup,
	); err != nil {
		if errors.IsNotFound(err) {
			return nil, fmt.Errorf("HTTPRouteGroup not found: %v", err)
		}
		return nil, fmt.Errorf("Failed to get HTTPRouteGroup: %v", err)
	}

	// Create a map to make it easier to search later
	matches := make(map[string]specsv1alpha1.HTTPMatch)
	for _, match := range httpRouteGroup.Matches {
		matches[match.Name] = match
	}

	var ret []specsv1alpha1.HTTPMatch
	for _, matchName := range spec.Matches {
		if _, ok := matches[matchName]; !ok {
			return nil, fmt.Errorf(
				"Match with name %s not found in HTTPRouteGroup %s",
				matchName, spec.Name)
		}
		ret = append(ret, matches[matchName])
	}
	return ret, nil
}
