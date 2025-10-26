#  FlowCD

A lightweight GitOps controller for Kubernetes clusters inspired by Argo CD.

## why i am building this?

Managing and deploying gitops applications with argocd or Fluxcd is quite complex as we want to understand each of the compenent how it working like as in argo we have two servers and controller for syncing and ui ,cli. Now i want make things simple for gitops deployment in k8s so i am building this Flowcd, just write your custom resouse like this 

```yaml
apiVersion: flowcd.io/v1alpha1
kind: FlowCD
metadata:
  name: nginx-app
  namespace: default
spec:
  source:
    repoURL: "https://github.com/sagnik3788/Gitops-controller.git"
    branch: "main"
    path: "manifests"
  destination:
    namespace: "default"
  deploymentStrategy:
    type: "QuickSync"
```

and write `kubectl apply -f flowcd.yaml` and your flowcd deployed on the default ns


no complexity simple architecture which you can scale easily and separately
we divided into 3 components

1. flowcd server 
2. flowcd controller 
3. flowctl cli

I liked the sync strategies of pipecd so i follow those  strategies 

- quick-sync
- pipeline-sync
- custom-sync

I know this project is not perfect for production env but ig we can improve it right :)