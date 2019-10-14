# Overview

The SMI Istio Adapter implements the SMI Traffic Split API as dictated by the [SMI spec](https://github.com/deislabs/smi-spec/blob/master/traffic-split.md).

See example TrafficSplit resource below:
```yaml
apiVersion: split.smi-spec.io/v1alpha1
kind: TrafficSplit
metadata:
  name: my-trafficsplit
spec:
  # The root service that clients use to connect to the destination application.
  service: website
  # Services inside the namespace with their own selectors, endpoints and configuration.
  backends:
  - service: website-v1
    weight: 80
  - service: website-v2
    weight: 20
```

The above TrafficSplit resource describes how traffic should be split between the services `website-v1` and `website-v2`. Under the hood, this adapter watches Traffic Split resources and creates/modifies Istio [Virtual Services](https://istio.io/docs/reference/config/networking/v1alpha3/virtual-service/) which then direct traffic to Kubernetes Services accordingly. Once this example TrafficSplit resource is created in Kubernetes, a corresponding Istio `VirtualService` will be created called `my-trafficsplit-vs`. Note the naming convention between the `TrafficSplit` resource and the `VirtualService` resource. The name of the `VirtualService` will be the same as the `TrafficSplit` resource with `-vs` appended to it.

## Weights
A user can specify a field called `weight` for every service specified in the list of backends on the Traffic Split resource.

The weights above correlate to the weight on a corresponding `VirtualService` object in Istio. There are some rules around weights in Istio Virtual Services though:
1. Weights must be whole numbers 0-100.
2. The sum of the weights must equal 100

To achieve the most intentional behavior, it is a best practice to ensure that the weights specified in a TrafficSplit resource meet these constraints. A `weight` in a TrafficSplit object must be a whole number as is specified in the SMI spec, however there is no constraint/limit on what that number can be. If you specify numbers that do not sum to 100, the SMI Istio adapter will take the total weight and assign whole numbers to each service that correspond to the percentage of how much traffic each weight should receive. If the weights don't cleanly add up to 100, then the last service will be rounded such that the total weight does equal 100 and the corresponding Istio `VirtualService` will be created with those re-calculated weights.

## Demo
Take the Traffic Split functionality for a spin with [this demo](smi-trafficsplit/README.md).
