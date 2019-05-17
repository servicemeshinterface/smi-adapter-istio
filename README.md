# smi-adapter-istio

## Contributing

Please refer to [CONTRIBUTING.md](./CONTRIBUTING.md) for more information on contributing to the specification.

## How to build

- Install the [Operator SDK CLI](https://github.com/operator-framework/operator-sdk/blob/master/doc/user-guide.md#install-the-operator-sdk-cli)
- Choose the container image name for the operator and build:

```
export OPERATOR_IMAGE=docker.io/servicemeshinterface/smi-adapter-istio:latest
make
```

- Push on your container registry:

```
make push
```

## How to install

After installing Istio you can deploy the adapter in the `istio-system` namespace with:

```bash
kubectl apply -f deploy/crds/split_v1alpha1_trafficsplit_crd.yaml
kubectl -n istio-system apply -f deploy/rbac.yaml
cat deploy/operator.yaml | sed "s,OPERATOR_IMAGE,$OPERATOR_IMAGE,g" | kubectl apply -f -
```

## Documentation

- SMI [`TrafficSplit`](https://github.com/deislabs/smi-spec/blob/master/traffic-split.md) spec [usage guide](docs/smi-trafficsplit).
- Using [Flagger](https://docs.flagger.app/) with SMI `TrafficSplit` [usage guide](docs/smi-flagger).
