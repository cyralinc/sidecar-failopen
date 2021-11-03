# CloudFormation Template for Cyral Sidecar DNS Fail Open for AWS

## Introduction

This repository contains the CloudFormation template that deploys the Cyral Sidecar Fail Open feature.
This feature provides automatic fail open/close to a Cyral Sidecar and its respectives target repositories,
allowing customers to keep existing databases reachable even when Cyral Sidecar experiences transient
failures.

![Cyral Sidecar Fail Open - Overview](./img/fail_open_overview.png)

The Cyral Sidecar Fail Open for AWS is built on top of CloudWatch, Lambda and Route53. The CloudFormation 
template contained here deploys the whole infrastructure and will provide the fail open feature out of
the box.

The overrall design is based on an event that triggers on a specific interval and perform health checks.
In case health checks fail, an alarm will be raised which will then trigger the change in the DNS record.
The architecture is based on AWS' own way of liveness probing resources in private subnets, which can be 
found [here](https://aws.amazon.com/blogs/networking-and-content-delivery/performing-route-53-health-checks-on-private-resources-in-a-vpc-with-aws-lambda-and-amazon-cloudwatch/).
Our architecture is described in the image below:

![Cyral Sidecar Fail Open for AWS - Architecture](./img/fail_open_aws.png)


The lambda function lives in [its own repo](https://github.com/cyralinc/health-check-aws).

# Limitations

Some limitations apply to the operation of the fail open as follows:

## DNS CNAME

One CNAME must be used per Cyral Sidecar and per repository. It means that if two different 
repositories are bound to the same Cyral Sidecar, then one CNAME must be created to represent
each repository. This will allow separate health checks as Cyral Sidecar is designed as a
modular architecture and support for different databases are completely independent at
Sidecar level.

![One CNAME per Cyral Sidecar per repository](./img/fail_open_cname_conf.png)

## Credentials

Only repositories that are bind to Cyral Sidecars and that accept native credentials are supported.
This means that repositories that use SSO credentials exclusively are **not supported**. This 
limitation is due to the fact that the same credentials used to check if the Cyral Sidecar is
healthy are also used to check if the database is healthy, so native database credentials are 
required.

## Port Definition

To allow client applications to restore database connectivity after a fail open event, the port
allocated in the sidecar for the repository must be the same as the database port. For example, 
If your MySQL instance is listening on port `3306`, then you need to make this repository 
available in the sidecar on port `3306` as well. This parity will guarantee that client 
applications that refers to `finance.db.acme.com` on port `3306` will still be able to connect
if the CNAME moves from `sidecar.db.acme.com` to `mysql-01.db.acme.com` and vice-versa.

## Deployment

One health check must be deployed [per CNAME](#dns-cname), thus
one lambda will be deployed for each pair Sidecar CNAME + Repository CNAME.


# Stack Deployment

The stack can be deployed to CloudFormation via AWS Console or AWS CLI. All the necessary parameters
and detailed information on each of them are in the `Parameters` section of the 
[CloudFormation Template](./cft_sidecar_failopen.yaml). Before you deploy, read the entire 
[Deployment Pre-Requisites](#deployment-pre-prequisites) section and make sure to have all of them
ready before you give it a go.

## Deployment Pre-Requisites

- The sidecar load balancer DNS name (output variable `SidecarLoadBalancerDNS` from CFT sidecar).

- The hosted zone that will be used to create the record sets.

- A list of subnets with access to the sidecar, the repository, AWS CloudWatch and AWS SecretsManager.
These subnets are the ones where the lambda will be deployed to. One way to configure access to 
CloudWatch and SecretsManager from Labmda is to [follow these steps](https://aws.amazon.com/premiumsupport/knowledge-center/internet-access-lambda-function/) 
up to item `Verify that your network ACL allows outbound requests from your Lambda function, and inbound traffic as needed`
and attach the lambda to the created subnets.

- Database secrets stored in AWS Secrets Manager in the format:

```json
{
  "username": "",
  "password": ""
}
```

- The database name to which the health check commands will be executed against.

