# smi-adapter-istio

This is a Kubernetes operator which implements the [Traffic Split](https://github.com/deislabs/smi-spec/blob/master/traffic-split.md), [Traffic Access Control](https://github.com/deislabs/smi-spec/blob/master/traffic-access-control.md) and [Traffic Specs](https://github.com/deislabs/smi-spec/blob/master/traffic-specs.md) APIS from the Service Mesh Interface (SMI) to use with Istio.

Tools or humans may set up and use this operator after installing Istio to do things like:
- orchestrate canary releases for new versions of software or more generally manage traffic shifting over time for applications
- define which services are allowed to send traffic to another service (and even a route on a service) or more generally define access control policies for applications

SMI defines a set of CRDs that allows for a common set of interfaces to build on top of when building tooling or working with service mesh implementations like Istio. This project builds logic around those commonly defined CRDs to work specifically with Istio.

## Contributing

Please refer to [CONTRIBUTING.md](./CONTRIBUTING.md) for more information on contributing to the specification.

## How to build

- Install the [Operator SDK CLI](https://github.com/operator-framework/operator-sdk/blob/master/doc/user-guide.md#install-the-operator-sdk-cli)
- Choose the container image name for the operator and build:

```bash
export OPERATOR_IMAGE=docker.io/servicemeshinterface/smi-adapter-istio:latest
make
```

- Push on your container registry:

```bash
make push
```

## How to install

After [installing Istio](https://istio.io/docs/setup/kubernetes/install/kubernetes/) you can deploy the adapter in the `istio-system` namespace with:

```bash
kubectl apply -R -f deploy/
```

## Documentation

- SMI [`TrafficSplit`](https://github.com/deislabs/smi-spec/blob/master/traffic-split.md) spec [usage guide](docs/smi-trafficsplit).
- Using [Flagger](https://docs.flagger.app/) with SMI `TrafficSplit` [usage guide](docs/smi-flagger).
