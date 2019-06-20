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

## Developer setup

### Prerequisites

- Minikube machine with atleast `4GB` memory.
- `operator-sdk` installed, download from [here](https://github.com/operator-framework/operator-sdk/releases).

### Download Istio

Follow [instructions in istio documentation](https://istio.io/docs/setup/kubernetes/download/#download-and-prepare-for-the-installation) to download istio.

### Install Istio

Once you are the root of the downloaded directory. Run following command to install Istio.

```bash
kubectl apply -f install/kubernetes/istio-demo-auth.yaml
```

Verify the istio is running fine as mentioned [here](https://istio.io/docs/setup/kubernetes/install/kubernetes/#verifying-the-installation).

### Now build the operator image

```bash
eval $(minikube docker-env)
operator-sdk build devimage
```

By exporting all the minikube docker environment variables locally the build happens in the virtual machine directly.

### Deploy the operator and related configs

```bash
kubectl apply -R -f deploy/
cat deploy/operator.yaml | sed 's|OPERATOR_IMAGE|devimage|g'| sed 's|imagePullPolicy: Always|imagePullPolicy: Never|g' | kubectl apply -f -
```

### To rebuild and redeploy

After making changes to the code run following commands

```bash
eval $(minikube docker-env)
operator-sdk build devimage
kubectl -n istio-system delete pod -l 'name=smi-adapter-istio'
```
