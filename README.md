# sample-controller

A custom controller that deletes Jobs that have passed the specified time after completion.

## Install

```
kubectl apply -f aartifacts/00-crds.yaml
kubectl apply -f artifacts/01-namespace.yaml
kubectl apply -f artifacts/02-controller.yaml
```
