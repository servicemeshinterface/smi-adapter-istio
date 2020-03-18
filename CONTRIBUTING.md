# Contributing

## Support Channels

Whether you are a user or contributor, official support channels include:

- [Issues](https://github.com/servicemeshinterface/smi-spec/issues)
- #smi Slack channel in the [CNCF Slack](https://cloud-native.slack.com)

## Sign Your Work

The sign-off is a simple line at the end of the explanation for a commit. All
commits needs to be signed. Your signature certifies that you wrote the patch or
otherwise have the right to contribute the material. The rules are pretty simple,
if you can certify the below (from [developercertificate.org](https://developercertificate.org/)):

```text
Developer Certificate of Origin
Version 1.1

Copyright (C) 2004, 2006 The Linux Foundation and its contributors.
1 Letterman Drive
Suite D4700
San Francisco, CA, 94129

Everyone is permitted to copy and distribute verbatim copies of this
license document, but changing it is not allowed.

Developer's Certificate of Origin 1.1

By making a contribution to this project, I certify that:

(a) The contribution was created in whole or in part by me and I
    have the right to submit it under the open source license
    indicated in the file; or

(b) The contribution is based upon previous work that, to the best
    of my knowledge, is covered under an appropriate open source
    license and I have the right under that license to submit that
    work with modifications, whether created in whole or in part
    by me, under the same open source license (unless I am
    permitted to submit under a different license), as indicated
    in the file; or

(c) The contribution was provided directly to me by some other
    person who certified (a), (b) or (c) and I have not modified
    it.

(d) I understand and agree that this project and the contribution
    are public and that a record of the contribution (including all
    personal information I submit with it, including my sign-off) is
    maintained indefinitely and may be redistributed consistent with
    this project or the open source license(s) involved.
```

Then you just add a line to every git commit message:

```text
    Signed-off-by: Joe Smith <joe.smith@example.com>
```

Use your real name (sorry, no pseudonyms or anonymous contributions.)

If you set your `user.name` and `user.email` git configs, you can sign your
commit automatically with `git commit -s`.

Note: If your git config information is set properly then viewing the
 `git log` information for your commit will look something like this:

```text
Author: Joe Smith <joe.smith@example.com>
Date:   Thu Feb 2 11:41:15 2018 -0800

    Update README

    Signed-off-by: Joe Smith <joe.smith@example.com>
```

Notice the `Author` and `Signed-off-by` lines match. If they don't
your PR will be rejected by the automated DCO check.

# Community Code of Conduct

Service Mesh Interface follows the [CNCF Code of Conduct](https://github.com/cncf/foundation/blob/master/code-of-conduct.md).

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
cat deploy/kubernetes-manifests.yaml | sed 's|servicemeshinterface/smi-adapter-istio:latest|devimage|g'| sed 's|imagePullPolicy: Always|imagePullPolicy: Never|g' | kubectl apply -f -
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
- Replace `servicemeshinterface/smi-adapter-istio` in [manifest](deploy/kubernetes-manifests.yaml) and [Tiltfile](Tiltfile) with your own image name i.e. `<dockeruser>/smi-adapter-istio`
- Run `$ tilt up` in project directory

This will build and deploy the operator to Kubernetes and you can iterate and watch changes get updated in the cluster!
