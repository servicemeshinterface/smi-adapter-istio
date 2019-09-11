# Using Flagger

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

### Install SMI CRDs and Operator to work with Istio

```bash
cd $GOPATH/src/github.com/deislabs/smi-adapter-istio/docs/smi-flagger
kubectl apply -f https://raw.githubusercontent.com/deislabs/smi-adapter-istio/master/deploy/crds/crds.yaml
kubectl apply -f https://raw.githubusercontent.com/deislabs/smi-adapter-istio/master/deploy/operator-and-rbac.yaml
```

### Install Flagger

```bash
kubectl apply -f manifests/00-install-flagger.yaml
```

## Demo

### Deploy `v1` of Bookinfo application

```bash
kubectl create ns smi-demo
kubectl label namespace smi-demo istio-injection=enabled
kubectl -n smi-demo apply -f manifests/00-bookinfo-setup.yaml
kubectl -n smi-demo apply -f manifests/01-deploy-v1.yaml
```

### Create the canary object

```bash
kubectl -n smi-demo apply -f manifests/02-canary.yaml
```

### Verify the application works

Run following command to get the URL to acces the application.

```bash
echo "http://$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.status.loadBalancer.ingress[0].ip}')/productpage"
```

It generally looks like following: `http://52.168.69.51/productpage`.

Visit above URL in the browser and you won't see any `ratings` since the `reviews v1` does not call `ratings` service at all.

### Deploy LoadTester

Deploy the load testing service to generate traffic during the canary analysis:

```bash
export REPO=https://raw.githubusercontent.com/weaveworks/flagger/master
kubectl -n smi-demo apply -f ${REPO}/artifacts/loadtester/deployment.yaml
kubectl -n smi-demo apply -f ${REPO}/artifacts/loadtester/service.yaml
```

### Deploy `v2` of reviews application

```bash
kubectl -n smi-demo set image deploy reviews reviews=istio/examples-bookinfo-reviews-v2:1.10.1
```

Once any changes are detected to the deployment it will be gradually scaled up. You can refresh the page to see that the traffic to new version is gradually increased.

### Verify successful deployment of v2

Run following command to see that the deployment is progressing. Also meanwhile you can refresh browser to see that the traffic sent to version `v2` of the reviews micro-service is increasing.

```bash
kubectl -n smi-demo describe canaries reviews
```

Once the output looks like following that means the newer version is successfully deployed. And all of the traffic will be sent to the newer version.

```console
$ kubectl -n smi-demo get canaries
NAME      STATUS      WEIGHT   LASTTRANSITIONTIME
reviews   Succeeded   0        2019-05-16T14:45:15Z
```

Also you can observe that the TrafficSplit object is created.

```bash
kubectl get trafficsplit
```
