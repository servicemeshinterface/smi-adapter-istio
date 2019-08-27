# smi-adapter-istio

This is a Kubernetes operator which implements the Service Mesh Interface(SMI) [Traffic Split](https://github.com/deislabs/smi-spec/blob/master/traffic-split.md), [Traffic Access Control](https://github.com/deislabs/smi-spec/blob/master/traffic-access-control.md) and [Traffic Specs](https://github.com/deislabs/smi-spec/blob/master/traffic-specs.md) APIs to work with Istio.

Tools or humans may set up and use this operator after installing Istio to do things like:
- orchestrate canary releases for new versions of software or more generally manage traffic shifting over time for applications
- define which services are allowed to send traffic to another service (and even a route on a service) or more generally define access control policies for applications

SMI defines a set of CRDs that allow for a common set of interfaces to build on top of when building tooling or working with service mesh implementations like Istio. This project builds logic around those commonly defined CRDs to work specifically with Istio.

## Getting Started

### Prerequesites
- Running Kubernetes cluster
- `kubectl` is installed locally
- [Istio installed](https://istio.io/docs/setup/kubernetes/install/kubernetes/) on Kubernetes cluster

### Install operator
1. Clone this repository in $GOPATH
```console
$ mkdir -p $GOPATH/src/github.com/deislabs/
$ cd $GOPATH/src/github.com/deislabs
$ git clone git@github.com:deislabs/smi-adapter-istio.git
```
2. cd into project directory
```console
$ cd $GOPATH/src/github.com/deislabs/smi-adapter-istio
```
3. Deploy this operator and all necessary CRDs in the `istio-system` namespace with:
```console
$ kubectl apply -R -f deploy/
```
4. Check that the operator has been deployed
```console
$ kubectl get pods -n istio-system -l name=smi-adapter-istio
NAME                                      READY     STATUS      RESTARTS   AGE
smi-adapter-istio-5ffcm8fqm               1/1       Running     0          20s
```

### Deploy SMI custom resources for managing traffic to services
- See how to use this project for traffic splitting [here](docs/smi-trafficsplit)
- See how to use this project to restrict service to service communication [here](docs/smi-traffictarget)
- See how to use Flagger and this project to do canary deploys [here](docs/smi-flagger)

## Contributing

Please refer to [CONTRIBUTING.md](./CONTRIBUTING.md) for more information on contributing to the specification.

## How to build

- Install the [Operator SDK CLI](https://github.com/operator-framework/operator-sdk/blob/master/doc/user-guide.md#install-the-operator-sdk-cli)
- Choose the container image name for the operator and build:

```bash
export OPERATOR_IMAGE=docker.io/<your username>/smi-adapter-istio:latest
make
```

- Push on your container registry:

```bash
make push
```
