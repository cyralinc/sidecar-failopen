package failopen

import (
	"context"

	"github.com/cyralinc/cloudformation-sidecar-failopen/internal/cloudwatch"
	"github.com/cyralinc/cloudformation-sidecar-failopen/internal/config"
	"github.com/cyralinc/cloudformation-sidecar-failopen/internal/healthcheck"
	"github.com/cyralinc/cloudformation-sidecar-failopen/internal/logging"
	_ "github.com/cyralinc/cloudformation-sidecar-failopen/internal/repository/oracle"
	_ "github.com/cyralinc/cloudformation-sidecar-failopen/internal/repository/mysql"
	_ "github.com/cyralinc/cloudformation-sidecar-failopen/internal/repository/postgresql"
)

func Run(ctx context.Context) error {

	logging.Info("performing the health check")
	err := healthcheck.HealthCheck(ctx, config.Config())
	if err != nil {
		logging.Info("health check performed, sidecar unhealthy. Setting metric on cloudwatch")
		err := cloudwatch.LogUnhealthy(ctx)

		if err != nil {
			logging.Error("error when connecting to cloudwatch: %s", err)
			return err
		}
	} else {
		logging.Info("health check performed, sidecar healthy. Setting metric on cloudwatch")
		err := cloudwatch.LogHealthy(ctx)
		if err != nil {
			logging.Error("error when connecting to cloudwatch: %s", err)
			return err
		}
	}
	return nil
}
