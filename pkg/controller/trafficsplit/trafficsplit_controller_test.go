package trafficsplit

import (
	"reflect"
	"testing"

	splitv1alpha1 "github.com/deislabs/smi-sdk-go/pkg/apis/split/v1alpha1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func quantityOrDie(str string) resource.Quantity {
	val, err := resource.ParseQuantity(str)
	if err != nil {
		panic(err)
	}
	return val
}

func TestNewVS(t *testing.T) {
	split := splitv1alpha1.TrafficSplit{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo",
			Namespace: "namespace",
		},
		Spec: splitv1alpha1.TrafficSplitSpec{
			Service: "foobar",
			Backends: []splitv1alpha1.TrafficSplitBackend{
				splitv1alpha1.TrafficSplitBackend{
					Service: "baz",
					Weight:  quantityOrDie("100"),
				},
				splitv1alpha1.TrafficSplitBackend{
					Service: "blah",
					Weight:  quantityOrDie("900"),
				},
			},
		},
	}
	svc := newVSForCR(&split)

	if svc == nil {
		t.Errorf("unexpected nil")
	}

	if svc.Namespace != split.Namespace {
		t.Errorf("expected %s, saw %s", split.Namespace, svc.Namespace)
	}

	if svc.Name != split.Name+"-vs" {
		t.Errorf("expected %s-vs, saw %s", split.Name, svc.Name)
	}

	if svc.Labels["traffic-split"] != split.Name {
		t.Errorf("expected %s, saw %s", split.Name, svc.Labels["traffic-split"])
	}

	if !reflect.DeepEqual(svc.Spec.Hosts, []string{split.Spec.Service}) {
		t.Errorf("expected %#v, saw %#v", split.Spec.Service, svc.Spec.Hosts)
	}

	routes := svc.Spec.Http[0].Route
	if len(routes) != len(split.Spec.Backends) {
		t.Errorf("expected %#v, saw %#v", split.Spec.Backends, svc.Spec.Http[0].Route)
	}
	for ix := range routes {
		if routes[ix].Destination.Host != split.Spec.Backends[ix].Service {
			t.Errorf("expected %s, saw %s", split.Spec.Backends[ix].Service, routes[ix].Destination.Host)
		}
		expectedWeight := int32(getIstioBackendPercentage(&split, ix))
		if routes[ix].Weight != expectedWeight {
			t.Errorf("expected %d, saw %d", expectedWeight, routes[ix].Weight)
		}
	}
}
