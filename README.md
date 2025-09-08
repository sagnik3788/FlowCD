#  FlowCD

A lightweight GitOps controller for Kubernetes clusters inspired by Pipecd.

## Features

- **Quick sync with Git repo**  
  Automatically syncs manifests from your Git repository.

- **Deploy to local Kubernetes cluster**  
  Applies synced manifests to keep your cluster in sync with Git.

## Planned Features

- Drift Detection
- Sync Strategies + metrics
- cli for interaction with cp(flowctl)
- Application CRD
- rollback and multi-cluster
## Getting Started

1. Configure your Git repository URL in the controller settings.  
2. Run the controller in your Kubernetes cluster or locally.  
3. The controller will sync and apply manifests from Git automatically.

## How to Access Traefik UI

- Traefik is deployed as part of this setup with a NodePort service.  
- Access the Traefik dashboard UI using the cluster node IP and port `30090`.
