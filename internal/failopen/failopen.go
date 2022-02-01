package failopen

import (
	"context"

	"github.com/cyralinc/sidecar-failopen/internal/cloudwatch"
	"github.com/cyralinc/sidecar-failopen/internal/config"
	"github.com/cyralinc/sidecar-failopen/internal/healthcheck"
	"github.com/cyralinc/sidecar-failopen/internal/logging"

	// registering drivers for the healthcheck
	// mysql covers mysql and mariadb
	_ "github.com/cyralinc/sidecar-failopen/internal/repository/mysql"
	_ "github.com/cyralinc/sidecar-failopen/internal/repository/oracle"

	// snowflake
	_ "github.com/cyralinc/sidecar-failopen/internal/repository/snowflake"

	// pg covers postgresql, denodo and redshift
	_ "github.com/cyralinc/sidecar-failopen/internal/repository/postgresql"
)

func Run(ctx context.Context) error {
	cfg := config.Config()
	logging.Init(cfg.LogLevel)

	logging.Info("performing the health check")
	err := healthcheck.HealthCheck(ctx, cfg)
	if err != nil {
		logging.Info("health check performed, sidecar unhealthy. Setting metric on cloudwatch. Sidecar error: %s", err)
		err := cloudwatch.LogUnhealthy(ctx, cfg)

		if err != nil {
			logging.Error("error when connecting to cloudwatch: %s", err)
			return err
		}
	} else {
		logging.Info("health check performed, sidecar healthy. Setting metric on cloudwatch")
		err := cloudwatch.LogHealthy(ctx, cfg)
		if err != nil {
			logging.Error("error when connecting to cloudwatch: %s", err)
			return err
		}
	}
	return nil
}
