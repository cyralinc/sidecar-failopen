# CloudFormation template for sidecar DNS Fail-Open

This is a template that creates a DNS Fail-Open system, utilizing a Lambda function, CloudWatch alarms and Route53 health-checks to
switch a DNS recordset between a sidecar and a repo, considering the current status of the sidecar.

The overrall design is based on an event that triggers on a specific interval, which performs the healthcheck and activates
an alarm in case the sidecar is failing. The architecture is described in the chart below:

![Architectural Chart](./img/chart.png)

The lambda that is used is in [its own repo](https://github.com/cyralinc/health-check-aws).

The architecture is based on AWS’ own way of liveness probing resources in private subnets, which can be found in this link.
Their architecture emulates the possibilities for Route53 health checks with the lambda acting as a bridge between the health-check that can only monitor publicly named resources and the sidecar that is contained in a private subnet.


## Deployment

### Pre-Requisites
- The hosted zone that will be used to create the RecordSets

- A subnet for the lambda function with access to CloudWatch and SecretsManager.

- An ECR to store the lambda function that needs to be either public or in the same region as the lambda.

### Configuration

The template asks for sidecar configuration, such as its FQDN and the port that will be checked.

The DB configuration is based on environment variables and asks for the username, password and database for the native credentials on the repo. These values can be changed to be retrieved from secrets, which would be an improvement if necessary.

The lambda configuration needs its VPC and subnets. Keep in mind that the subnet/VPC needs to be the in the same VPC as the sidecar, and the VPC needs to attend to the second pre-requisite of having access to outbound internet.

| Variable                      | Description                                                                                                                                                                                  |
| ---                           | ---                                                                                                                                                                                          |
| SidecarAddress                | The hostname of the sidecar                                                                                                                                                                  |
| SidecarNamePrefix             | The name prefix of the sidecar. This is not explicitly necessary, but it’s good for identifying the metric that will be set by the lambda function.                                          |
| SidecarPort                   | Port that has been allocated for the sidecar on the control plane.                                                                                                                           |
| DBSecretLocation              | Location of the secret on AWS that contains the configuration for the repository. The secret must be in the same region as the lambda function.                                              |
| DBAddress                     | The hostname of the database. Can either be an IP address or a FQDN.                                                                                                                         |
| HostedZoneID                  | The ID of the hosted zone that will contain the recordset that will serve as a fail over.                                                                                                    |
| RecordSetName                 | The name of the recordset that will serve as a failover. Must be a name that is included in the Hosted Zone that was selected.                                                               |
| DBRecordSetType               | The type of the recordset for the Database. If it’s an IP, it should be AAAA, if it’s a FQDN, should be CNAME.                                                                               |
| Subnets                       | The subnets for the lambda function. Should contain the sidecar subnet and must have access to the Secret Manager and CloudWatch APIs, either via internet or via an Interface VPC Endpoint. |
| VPC                           | The VPC for the lambda function.                                                                                                                                                             |
| ImageUri                      | The URI for the image that will actually perform the healthcheck on the sidecar. Must be stored either in a ECR registry on the same region as the lambda or on a public ECR registry.       |
| NumberOfRetries               | Number of times the lambda function will retry to perform the health check if unable to connect before stopping.                                                                             |
| ConsecutiveFailuresForTrigger | Number of times the healthcheck must fail in a row to trigger the alarm and the failover. This will increase the total time the sidecar needs to be down before the fail over triggers.      |
