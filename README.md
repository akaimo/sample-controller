# sample-controller

A custom controller that deletes Jobs that have passed the specified time after completion.

## Install

```
kubectl apply -f aartifacts/00-crds.yaml
kubectl apply -f artifacts/01-namespace.yaml
kubectl apply -f artifacts/02-controller.yaml
```

## Configuration
```
apiVersion: samplecontroller.k8s.io/v1
kind: SampleResource
metadata:
  name: test
  namespace: default
spec:
  # It will be deleted after this time
  # Specify time as "s", "m", "h"
  time: "1m"
```
