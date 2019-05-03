# smi-adapter-istio

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

```
kubectl create -f deploy/service_account.yaml
kubectl create -f deploy/role.yaml
kubectl create -f deploy/role_binding.yaml
cat deploy/operator.yaml | sed "s,OPERATOR_IMAGE,$OPERATOR_IMAGE,g" | kubectl create -f -
```
