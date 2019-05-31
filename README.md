# smi-adapter-istio

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
