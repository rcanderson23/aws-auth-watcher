aws-auth-watcher
=========

Watches the `aws-auth` ConfigMap in EKS for changes. Fires an SNS notification on change.

Deployment
----------

1. Create IAM Policy and Role for sns:Publish permissions. This can be done by attaching the role to the instance or using [IRSA](https://docs.aws.amazon.com/eks/latest/userguide/iam-roles-for-service-accounts.html) 
2. Configure environment variables in `deploy/deployment.yaml` to match region and topic.
3. If using IRSA, modify `deploy/rbac.yaml` annotation in the ServiceAccount to use desired role.
4. `kubectl apply -f deploy/`
5. Make test modification to `aws-auth` to confirm desired behavior.

IAM
----------
Example Policy:
```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "VisualEditor0",
            "Effect": "Allow",
            "Action": "sns:Publish",
            "Resource": "arn:aws:sns:<region>:<account>:aws-auth"
        }
    ]
}
```

Example Email:
----------
```
Old ConfigMap Data:
mapRoles="- groups:
  - system:bootstrappers
  - system:nodes
  - system:masters
  rolearn: arn:aws:iam::123456789012:role/eksctl-cluster-nodegroup-sta-NodeInstanceRole-1IGP8NHFBU3G6
  username: system:node:{{EC2PrivateDNSName}}
"


New ConfigMap Data:
mapRoles="- groups:
  - system:bootstrappers
  - system:nodes
  rolearn: arn:aws:iam::123456789012:role/windows-ng-NodeInstanceRole-3NDOIJ74SOZ7
  username: system:node:{{EC2PrivateDNSName}}
- groups:
  - system:bootstrappers
  - system:nodes
  - system:masters
  rolearn: arn:aws:iam::123456789012:role/eksctl-cluster-nodegroup-sta-NodeInstanceRole-1IGP8NHFBU3G6
  username: system:node:{{EC2PrivateDNSName}}
"
```

TODO
----------

* Add Slack functionality
