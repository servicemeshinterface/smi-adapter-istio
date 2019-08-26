# Testing SMI Istio Adapter

Make sure you have setup Istio on a cluster.

#### Initial Setup

Install the bookinfo applications except the reviews application:

```bash
kubectl label namespace default istio-injection=enabled
kubectl apply -f 00-bookinfo-setup.yaml
```

**Note:** Change namespace to suit your needs.

#### Install reviews v1 application

```bash
kubectl apply -f 01-deploy-v1.yaml
```

#### Verify the application works for version v1

```bash
curl `minikube ip`:31380/productpage
```

OR visit above URL in the browser and see that you won't see any `ratings` since the `reviews v1` does not call `ratings` service at all.


#### Create TrafficSplit for v1

```bash
kubectl apply -f 02-traffic-split-1000-0.yaml
```

If you don't have SMI Istio Adapter installed run following:

```
kubectl apply -f 02-traffic-split-1000-0-output.yaml
```

#### Install reviews v2 application

```
kubectl apply -f 03-deploy-v2.yaml
```

You can try following step to verify that the traffic is still only sent to `v1` of the application even though there is `v2` version of the application deployed.

```bash
curl `minikube ip`:31380/productpage
```

Until you perform the next step no traffic will be sent to the `v2`.

#### Create TrafficSplit to send 33% of traffic to v2

```bash
kubectl apply -f 04-traffic-split-1000-500.yaml
```

If you don't have SMI Istio Adapter installed run following:

```
kubectl apply -f 04-traffic-split-1000-500-output.yaml
```

Refresh browser multiple times to see that the traffic is sent to `v1` `67%` of the time and to `v2` `33%` of the time.

#### Verify the v2 to send all traffic

If v2 of the reviews application looks good then start sending entire traffic to v2.

```bash
kubectl apply -f 05-traffic-split-0-1000.yaml
```

If you don't have SMI Istio Adapter installed run following:

```
kubectl apply -f 05-traffic-split-0-1000-output.yaml
```

#### Now cleanup unwanted resources

General Kubernetes cleanup

```bash
kubectl delete deployment reviews-v1
kubectl delete service reviews-v1
kubectl delete service reviews-v2
```

SMI Istio Adapter cleanup

```bash
kubectl delete trafficsplit reviews-rollout
```

If you don't have SMI Istio Adapter installed run following:

```bash
kubectl delete destinationrule reviews-rollout
kubectl delete virtualservice reviews-rollout
```

To clean everything from this test run following command:

```bash
kubectl delete -f .
```
