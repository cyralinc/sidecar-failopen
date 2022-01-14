# Cyral sidecar DNS fail-open for AWS

## Introduction

This repository contains the CloudFormation template that deploys the Cyral sidecar fail-open feature.
This feature provides automatic fail-open/fail-closed operation for a Cyral sidecar and its respective target repositories,
allowing customers to keep existing databases reachable even when the Cyral sidecar experiences transient
failures.

![Cyral Sidecar Fail Open - Overview](./img/fail_open_overview.png)

Cyral sidecar fail-open for AWS is built on top of CloudWatch, Lambda, and Route53. The CloudFormation 
template contained here deploys the whole infrastructure and will provide the fail-open feature out of
the box.

The overall design is based on an event that triggers at a specific interval and performs health checks.
If a health check fails, an alarm will be raised which will then trigger the change in the DNS record.
The architecture is based on AWS' own approach to liveness probing for resources in private subnets, which is
described [here](https://aws.amazon.com/blogs/networking-and-content-delivery/performing-route-53-health-checks-on-private-resources-in-a-vpc-with-aws-lambda-and-amazon-cloudwatch/).

Our architecture is described in the image below:

![Cyral Sidecar Fail Open for AWS - Architecture](./img/fail_open_aws.png)


The lambda function lives in [its own repo](https://github.com/cyralinc/sidecar-failopen-healthcheck-aws).

# Limitations

Some limitations apply to the operation of the fail-open feature, as described below.

## Repositories supported

PostgreSQL and MySQL are supported in the current version of [the CFT](./cft_sidecar_failopen.yaml).

Click [here for Snowflake support](./snowflake) and follow the instructions.

## DNS CNAME

One CNAME must be used per Cyral sidecar-repository combination. It means that if two different 
repositories are bound to the same Cyral sidecar, then one CNAME must be created to represent
each repository. This allows an independent health check per repository, based on the sidecar's 
modular design (a sidecar handles each database independently).

![One CNAME per Cyral Sidecar per repository](./img/fail_open_cname_conf.png)

## Native credentials required

Only repositories that are bound to Cyral sidecars and that accept native credentials are supported.
This means that repositories that use SSO credentials exclusively are **not supported**. This 
limitation is due to the fact that the same credentials used to check if the Cyral sidecar is
healthy are also used to check if the database is healthy, so native database credentials are 
required.

## Port definition

To allow client applications to restore database connectivity after a fail-open event, the port
allocated in the sidecar for the repository must be the same as the database port. For example, 
if your MySQL repository is listening on port `3306`, then you need to make this repository 
available in the sidecar on port `3306` as well. This parity guarantees that client 
applications that refer to `finance.db.acme.com` on port `3306` will still be able to connect
if the CNAME moves from `sidecar.db.acme.com` to `mysql-01.db.acme.com` and vice-versa.

## Deployment

One health check must be deployed [per CNAME](#dns-cname), thus
one lambda will be deployed for each sidecar CNAME/repository CNAME pair.


# Stack deployment

The stack can be deployed to CloudFormation via AWS Console or AWS CLI. All the necessary parameters
and detailed information on each of them are shown in the `Parameters` section of the 
[CloudFormation Template](./cft_sidecar_failopen.yaml). Before you deploy, read the entire 
[Deployment Pre-Requisites](#deployment-pre-prequisites) section and make sure you have all of them
ready before you give it a go.

## Deployment prerequisites

- The sidecar load balancer DNS name (output variable `SidecarLoadBalancerDNS` from CFT sidecar).

- The hosted zone that will be used to create the record sets.

- A list of subnets with access to the sidecar, the repository, AWS CloudWatch, and AWS Secrets Manager.
  These subnets are the ones where the lambda will be deployed. One way to configure access to 
  CloudWatch and Secrets Manager from Lambda is to [follow these steps](https://aws.amazon.com/premiumsupport/knowledge-center/internet-access-lambda-function/) 
  up to the step, `Verify that your network ACL allows outbound requests from your Lambda function, and inbound traffic as needed`,
  and then attach the lambda to the created subnets.

- Database secrets stored in AWS Secrets Manager in the format:

```json
{
  "username": "",
  "password": ""
}
```

- The database name against which the health check commands will be executed.

# FAQ

> - What is the time delay between the sidecar being unavailable and the Fail Open triggering?

This depends on the configuration you've set up. The default delay is around 4 minutes from the first instant where the sidecar becomes unavailable, 2 of those 4 being due to the `ConsecutivesFailuresForTrigger` parameter, which can be reduced from 2 to 1, causing this time to diminish to around 3 minutes. This parameter can also be increased if the trigger needs to wait longer before switching the DNS record.

> - What is the time delay between the sidecar being back up and the record switching back?

This takes around 2 minutes and **cannot be configured** as this is determined by the current behaviors of Route53 and CloudWatch.

> - After the Fail Open is triggered my application is still resolving to the same address.

This may be due to the DNS TTL in your runtime. The JVM default DNS TTL is 60s, for example. Refer to the documentation on your application/operational system to change it, if necessary.

> - The health check applied by the lambda does not conform to my needs.

The repository for the lambda is open source, and can be forked and updated as needed. You can publish it to your local AWS account and refer to your lambda store location when deploying the CloudFormation template.

> - Does this support snowflake repositories

At this time, snowflake specific repos have a separate fail open configuration that can be found within the [snowflake](./snowflake) directory in this repo.
