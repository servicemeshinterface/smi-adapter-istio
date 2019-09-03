# Contributing

## Support Channels

Whether you are a user or contributor, official support channels include:

- [Issues](https://github.com/deislabs/smi-spec/issues)
- #general Slack channel in the [SMI Slack](https://smi-spec.slack.com)

## CLA Requirement

This project welcomes contributions and suggestions. Most contributions require you to agree to a Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us the rights to use your contribution. For details, visit https://cla.microsoft.com.
When you submit a pull request, a CLA-bot will automatically determine whether you need to provide a CLA and decorate the PR appropriately (e.g., label, comment). Simply follow the instructions provided by the bot. You will only need to do this once across all repositories using our CLA.

## Code of Conduct

This project has adopted the [Microsoft Open Source Code of conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or contact opencode@microsoft.com (mailto:opencode@microsoft.com) with any additional questions or comments.

## Project Governance

Project maintainership is outlined in the [GOVERNANCE](GOVERNANCE.md) file.

## Developer setup

### Prerequisites
- Running Kubernetes cluster (If using Minikube, Minikube machine with atleast `4GB` memory).
- `operator-sdk` installed, download from [here](https://github.com/operator-framework/operator-sdk/releases).
- Follow [instructions in istio documentation](https://istio.io/docs/setup/kubernetes/download/#download-and-prepare-for-the-installation) to download istio. Once you are the root of the downloaded directory. Run following command to install Istio.
```bash
kubectl apply -f install/kubernetes/istio-demo-auth.yaml
```
Verify the istio is running fine as mentioned [here](https://istio.io/docs/setup/kubernetes/install/kubernetes/#verifying-the-installation).

### Build Operator Image

#### If using Minikube, use the following instructions:
```bash
eval $(minikube docker-env)
operator-sdk build devimage
```
By exporting all the minikube docker environment variables locally the build happens in the virtual machine directly.

Deploy image using:
```bash
kubectl apply -R -f deploy/
cat deploy/kubernetes-manifests.yaml | sed 's|deislabs/smi-adapter-istio:latest|devimage|g'| sed 's|imagePullPolicy: Always|imagePullPolicy: Never|g' | kubectl apply -f -
```

To rebuild and redeploy after making changes to the code, run the following commands:
```bash
eval $(minikube docker-env)
operator-sdk build devimage
kubectl -n istio-system delete pod -l 'name=smi-adapter-istio'
```
## Build and Push Operator Image
If you're not using Minikube, you'll want to build and push the image to a remote container registry. Choose the container image name for the operator and build:
```bash
export OPERATOR_IMAGE=docker.io/<your username>/smi-adapter-istio:latest
make
```

Push to your container registry:
```bash
make push
```

### Developing Using Tilt
- Install [Tilt](https://docs.tilt.dev/install.html)
- Replace `deislabs/smi-adapter-istio` in [manifest](deploy/kubernetes-manifests.yaml) and [Tiltfile](Tiltfile) with your own image name i.e. `<dockeruser>/smi-adapter-istio`
- Run `$ tilt up` in project directory

This will build and deploy the operator to Kubernetes and you can iterate and watch changes get updated in the cluster!
