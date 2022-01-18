package failopen

import (
	"context"

	"github.com/cyralinc/sidecar-failopen/internal/cloudwatch"
	"github.com/cyralinc/sidecar-failopen/internal/config"
	"github.com/cyralinc/sidecar-failopen/internal/healthcheck"
	"github.com/cyralinc/sidecar-failopen/internal/logging"
	_ "github.com/cyralinc/sidecar-failopen/internal/repository/denodo"
	_ "github.com/cyralinc/sidecar-failopen/internal/repository/mariadb"
	_ "github.com/cyralinc/sidecar-failopen/internal/repository/mysql"
	_ "github.com/cyralinc/sidecar-failopen/internal/repository/oracle"
	_ "github.com/cyralinc/sidecar-failopen/internal/repository/postgresql"
	_ "github.com/cyralinc/sidecar-failopen/internal/repository/redshift"
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
