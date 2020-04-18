aws-auth-watcher
=========

Watches the `aws-auth` ConfigMap in EKS for changes. Fires and SNS notification on change.

Deployment
----------

1. Configure environment variables in `deploy/deployment.yaml` to match region and topic.
2. `kubectl apply -f deploy/`
