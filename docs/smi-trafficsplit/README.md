# Using SMI spec TrafficSplit

## Setup

### Prerequisite

* `kubectl` installed.
* Working Kubernetes cluster.

### Install Istio

Download manifests as mentioned in upstream [docs](https://istio.io/docs/setup/kubernetes/download/#download-and-prepare-for-the-installation).


```bash
curl -L https://git.io/getLatestIstio | ISTIO_VERSION=1.1.6 sh -
```

Install Istio on Kubernetes

```bash
cd istio-1.1.6
kubectl apply -f install/kubernetes/istio-demo-auth.yaml
```

**NOTE**: Above apply might sometimes give errors like `unable to recognize "install/kubernetes/istio-demo-auth.yaml": no matches for kind "DestinationRule" in version "networking.istio.io/v1alpha3"`, this is most likely because the CRDs are not registered yet and apiserver will reconcile it. Try running the above `kubectl apply ...` again.

### Install SMI Operator for Istio

```bash
cd $GOPATH/src/github.com/deislabs/smi-adapter-istio/docs/smi-trafficsplit
kubectl apply -f deploy/crds/split_v1alpha1_trafficsplit_crd.yaml
kubectl -n istio-system apply -f deploy/rbac.yaml
export OPERATOR_IMAGE=docker.io/<your username>/smi-adapter-istio:latest
cat deploy/operator.yaml | sed "s,OPERATOR_IMAGE,$OPERATOR_IMAGE,g" | kubectl apply -f -
```

## Demo

### Deploy `v1` of Bookinfo application

```bash
kubectl create ns smi-demo
kubectl label namespace smi-demo istio-injection=enabled
kubectl -n smi-demo apply -f manifests/00-bookinfo-setup.yaml
kubectl -n smi-demo apply -f manifests/01-deploy-v1.yaml
```

### Verify the application works

Run following command to get the URL to acces the application.

```bash
echo "http://$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.status.loadBalancer.ingress[0].ip}')/productpage"
```

It generally looks like following: `http://52.168.69.51/productpage`.

Visit above URL in the browser and you won't see any `ratings` since the `reviews v1` does not call `ratings` service at all.

### Create TrafficSplit for `v1`

```bash
kubectl -n smi-demo apply -f manifests/02-traffic-split-1000-0.yaml
```

### Deploy `v2` of `reviews` application

```
kubectl -n smi-demo apply -f manifests/03-deploy-v2.yaml
```

Verify that the traffic is still only sent to version `v1` of the application even though there is `v2` version deployed. Until you perform the next step no traffic will be sent to the `v2`.

### Update TrafficSplit to send `10%` of traffic

```bash
kubectl -n smi-demo apply -f manifests/04-traffic-split-900-100.yaml
```

Refresh browser multiple times to see that the traffic is sent to version `v1` of reviews micro-service `90%` of the time and to `v2`, `10%` of the time. Verify that application works fine and increase the traffic for `v2`.

### Update TrafficSplit to send `25%` of traffic

```bash
kubectl -n smi-demo apply -f manifests/05-traffic-split-750-250.yaml
```

Refresh browser multiple times to see that the traffic is sent to version `v1` of reviews micro-service `75%` of the time and to `v2`, `25%` of the time. Verify that application works fine and increase the traffic for `v2`.

### Update TrafficSplit to send `80%` of traffic

```bash
kubectl -n smi-demo apply -f manifests/06-traffic-split-200-800.yaml
```

Refresh browser multiple times to see that the lesser traffic is sent to version `v1` of reviews micro-service which is `20%` of the time and to `v2`, `80%` of the time. Verify that application still works fine and now is the right time to send entire traffic to `v2`.

### Update TrafficSplit to send `100%` of traffic

```bash
kubectl -n smi-demo apply -f manifests/07-traffic-split-0-1000.yaml
```

Refresh browser and you will see that the traffic is sent only to version `v2` of the reviews application.

### Delete old version of application

```bash
kubectl -n smi-demo delete -f manifests/01-deploy-v1.yaml
```

### Get rid of the TrafficSplit object

```bash
kubectl -n smi-demo delete trafficsplit reviews-rollout
```
