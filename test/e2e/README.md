# Running e2e tests for smi-adapter-istio

## Prerequistes
- Running Kubernetes cluster
- [Install Istio](https://github.com/servicemeshinterface/smi-adapter-istio/tree/master/docs/smi-trafficsplit#install-istio) on Kubernetes cluster
- [Install operator-sdk cli](https://github.com/operator-framework/operator-sdk/blob/master/doc/user-guide.md#install-the-operator-sdk-cli)
- Install SMI CRDs on cluster
```bash
$ kubectl apply -f https://raw.githubusercontent.com/servicemeshinterface/smi-adapter-istio/master/deploy/crds/crds.yaml
```

## Run e2e tests

```bash
$ cd $GOPATH/src/github.com/servicemeshinterface/smi-adapter-istio/

$ operator-sdk test local ./test/e2e --namespaced-manifest deploy/operator-and-rbac.yaml --namespace istio-system
```
