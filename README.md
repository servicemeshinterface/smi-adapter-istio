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

After installing Istio you can deploy the adapter in the `istio-system` namespace with:

```
kubectl apply -f deploy/crds/split_v1alpha1_trafficsplit_crd.yaml
kubectl -n istio-system apply -f deploy/rbac.yaml

cat deploy/operator.yaml | sed "s,OPERATOR_IMAGE,$OPERATOR_IMAGE,g" | kubectl create -f -
```

## Developer setup

#### Prerequisites

- Minikube machine with atleast `4GB` memory.
- `operator-sdk` installed, download from [here](https://github.com/operator-framework/operator-sdk/releases).

#### Download Istio

Follow [instructions in istio documentation](https://istio.io/docs/setup/kubernetes/download/#download-and-prepare-for-the-installation) to download istio.

#### Install Istio

Once you are the root of the downloaded directory. Run following command to install Istio.

```bash
kubectl apply -f install/kubernetes/istio-demo-auth.yaml
```

Verify the istio is running fine as mentioned [here](https://istio.io/docs/setup/kubernetes/install/kubernetes/#verifying-the-installation).

#### Now build the operator image

```bash
eval $(minikube docker-env)
operator-sdk build devimage
```

By exporting all the minikube docker environment variables locally the build happens in the virtual machine directly.

#### Deploy the operator and related configs

```bash
sed -i 's|OPERATOR_IMAGE|devimage|g' deploy/operator.yaml
sed -i 's|imagePullPolicy: Always|imagePullPolicy: Never|g' deploy/operator.yaml
kubectl apply -R -f deploy/
```

#### To rebuild and redeploy

After making changes to the code run following commands

```bash
eval $(minikube docker-env)
operator-sdk build devimage
kubectl delete pod -l 'name=smi-adapter-istio'
```
