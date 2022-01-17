package cloudwatch

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/cyralinc/sidecar-failopen/internal/config"
)

var client *cloudwatch.CloudWatch

const dimensionFormat = "%s %s %s Health Check"
const metricNameFormat = "%s-%s-%s: %s (Health Check for resource %s)"

func logValue(value float64) error {
	namespace := "CyralSidecarHealthChecks"
	cfg := config.Config()
	metricName := fmt.Sprintf(metricNameFormat,
		cfg.Sidecar.Host,
		cfg.Repo.RepoType,
		cfg.Repo.Host,
		cfg.StackName,
		cfg.Sidecar.Host)
	dimentionName := fmt.Sprintf(dimensionFormat,
		cfg.Sidecar.Host,
		cfg.Repo.RepoType,
		cfg.Repo.Host,
	)

	_, err := client.PutMetricData(
		&cloudwatch.PutMetricDataInput{
			Namespace: &namespace,
			MetricData: []*cloudwatch.MetricDatum{
				{
					MetricName: &metricName,
					Dimensions: []*cloudwatch.Dimension{
						{
							Name:  &dimentionName,
							Value: &dimentionName,
						},
					},
					Value: &value,
				},
			},
		},
	)

	return err
}

func LogHealthy(ctx context.Context) error {
	return logValue(1)
}

func LogUnhealthy(ctx context.Context) error {
	return logValue(0)
}

func init() {
	mySession := session.Must(session.NewSession())
	client = cloudwatch.New(mySession)
}
