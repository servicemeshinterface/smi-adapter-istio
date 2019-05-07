package v1alpha3

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DestinationRuleSpec defines the desired state of DestinationRule
// +k8s:openapi-gen=true
type DestinationRuleSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html

	// REQUIRED. The name of a service from the service registry. Service
	// names are looked up from the platform's service registry (e.g.,
	// Kubernetes services, Consul services, etc.) and from the hosts
	// declared by [ServiceEntries](/docs/reference/config/networking/v1alpha3/service-entry/#ServiceEntry). Rules defined for
	// services that do not exist in the service registry will be ignored.
	//
	// *Note for Kubernetes users*: When short names are used (e.g. "reviews"
	// instead of "reviews.default.svc.cluster.local"), Istio will interpret
	// the short name based on the namespace of the rule, not the service. A
	// rule in the "default" namespace containing a host "reviews will be
	// interpreted as "reviews.default.svc.cluster.local", irrespective of
	// the actual namespace associated with the reviews service. _To avoid
	// potential misconfigurations, it is recommended to always use fully
	// qualified domain names over short names._
	//
	// Note that the host field applies to both HTTP and TCP services.
	Host string `protobuf:"bytes,1,opt,name=host,proto3" json:"host,omitempty"`
	// Traffic policies to apply (load balancing policy, connection pool
	// sizes, outlier detection).
	TrafficPolicy *TrafficPolicy `protobuf:"bytes,2,opt,name=traffic_policy,json=trafficPolicy,proto3" json:"traffic_policy,omitempty"`
	// One or more named sets that represent individual versions of a
	// service. Traffic policies can be overridden at subset level.
	Subsets []*Subset `protobuf:"bytes,3,rep,name=subsets,proto3" json:"subsets,omitempty"`
}

// Traffic policies to apply for a specific destination, across all
// destination ports. See DestinationRule for examples.
// +k8s:openapi-gen=true
type TrafficPolicy struct {
	// TLS related settings for connections to the upstream service.
	Tls *TLSSettings `protobuf:"bytes,4,opt,name=tls,proto3" json:"tls,omitempty"`
}

// A subset of endpoints of a service. Subsets can be used for scenarios
// like A/B testing, or routing to a specific version of a service. Refer
// to [VirtualService](/docs/reference/config/networking/v1alpha3/virtual-service/#VirtualService) documentation for examples of using
// subsets in these scenarios. In addition, traffic policies defined at the
// service-level can be overridden at a subset-level. The following rule
// uses a round robin load balancing policy for all traffic going to a
// subset named testversion that is composed of endpoints (e.g., pods) with
// labels (version:v3).
//
// ```yaml
// apiVersion: networking.istio.io/v1alpha3
// kind: DestinationRule
// metadata:
//   name: bookinfo-ratings
// spec:
//   host: ratings.prod.svc.cluster.local
//   trafficPolicy:
//     loadBalancer:
//       simple: LEAST_CONN
//   subsets:
//   - name: testversion
//     labels:
//       version: v3
//     trafficPolicy:
//       loadBalancer:
//         simple: ROUND_ROBIN
// ```
//
// **Note:** Policies specified for subsets will not take effect until
// a route rule explicitly sends traffic to this subset.
//
// One or more labels are typically required to identify the subset destination,
// however, when the corresponding DestinationRule represents a host that
// supports multiple SNI hosts (e.g., an egress gateway), a subset without labels
// may be meaningful. In this case a traffic policy with [TLSSettings](#TLSSettings)
// can be used to identify a specific SNI host corresponding to the named subset.
// +k8s:openapi-gen=true
type Subset struct {
	// REQUIRED. Name of the subset. The service name and the subset name can
	// be used for traffic splitting in a route rule.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Labels apply a filter over the endpoints of a service in the
	// service registry. See route rules for examples of usage.
	Labels map[string]string `protobuf:"bytes,2,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Traffic policies that apply to this subset. Subsets inherit the
	// traffic policies specified at the DestinationRule level. Settings
	// specified at the subset level will override the corresponding settings
	// specified at the DestinationRule level.
	TrafficPolicy *TrafficPolicy `protobuf:"bytes,3,opt,name=traffic_policy,json=trafficPolicy,proto3" json:"traffic_policy,omitempty"`
}

// SSL/TLS related settings for upstream connections. See Envoy's [TLS
// context](https://www.envoyproxy.io/docs/envoy/latest/api-v2/api/v2/auth/cert.proto.html)
// for more details. These settings are common to both HTTP and TCP upstreams.
//
// For example, the following rule configures a client to use mutual TLS
// for connections to upstream database cluster.
//
// ```yaml
// apiVersion: networking.istio.io/v1alpha3
// kind: DestinationRule
// metadata:
//   name: db-mtls
// spec:
//   host: mydbserver.prod.svc.cluster.local
//   trafficPolicy:
//     tls:
//       mode: MUTUAL
//       clientCertificate: /etc/certs/myclientcert.pem
//       privateKey: /etc/certs/client_private_key.pem
//       caCertificates: /etc/certs/rootcacerts.pem
// ```
//
// The following rule configures a client to use TLS when talking to a
// foreign service whose domain matches *.foo.com.
//
// ```yaml
// apiVersion: networking.istio.io/v1alpha3
// kind: DestinationRule
// metadata:
//   name: tls-foo
// spec:
//   host: "*.foo.com"
//   trafficPolicy:
//     tls:
//       mode: SIMPLE
// ```
//
// The following rule configures a client to use Istio mutual TLS when talking
// to rating services.
//
// ```yaml
// apiVersion: networking.istio.io/v1alpha3
// kind: DestinationRule
// metadata:
//   name: ratings-istio-mtls
// spec:
//   host: ratings.prod.svc.cluster.local
//   trafficPolicy:
//     tls:
//       mode: ISTIO_MUTUAL
// ```
// +k8s:openapi-gen=true
type TLSSettings struct {
	// REQUIRED: Indicates whether connections to this port should be secured
	// using TLS. The value of this field determines how TLS is enforced.
	Mode string `protobuf:"bytes,1,opt,name=mode,proto3" json:"mode,omitempty"`
}

// DestinationRuleStatus defines the observed state of DestinationRule
// +k8s:openapi-gen=true
// +k8s:openapi-gen=true
type DestinationRuleStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DestinationRule is the Schema for the destinationrules API
// +k8s:openapi-gen=true
type DestinationRule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DestinationRuleSpec   `json:"spec,omitempty"`
	Status DestinationRuleStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DestinationRuleList contains a list of DestinationRule
type DestinationRuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DestinationRule `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DestinationRule{}, &DestinationRuleList{})
}
