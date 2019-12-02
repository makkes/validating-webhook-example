# A minimal Kubernetes Validating Admission Webhook

This is an example repo to get started with [validating admission webhooks](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/)
on Kubernetes. It features a webhook written in Go as well as all resources needed to deploy it into a Kubernetes
cluster as a service.

## Getting Started

The easiest way to get started is to build the service, push the Docker image to Docker Hub and create the resources:

If you don't have `cert-manager` installed in your cluster, do that now:

```sh
$ make deploy-cert-manager
```

Then build, push and roll out the webhook:

```sh
$ make deploy
```

Try to spin up a pod in another terminal window:

```sh
k run debug --image=debian:latest --rm -it --restart=Never -- sh
```

Then you should be able to see all admission requests in the webhook's pod's logs:

```sh
$ kubectl logs -f deploy/minimal-webhook
```