// This package implements the health checking procedure for repositories.
// The procedure checks for liveness on the sidecar and on failure, checks
// for liveness in the repo. The test is repeated "FAIL_OPEN_N_RETRIES" times on
// sidecar failure.
package healthcheck

import (
	"context"
	"fmt"

	"github.com/cyralinc/sidecar-failopen/internal/config"
	"github.com/cyralinc/sidecar-failopen/internal/logging"
	"github.com/cyralinc/sidecar-failopen/internal/repository"
)

func singleHealthCheck(ctx context.Context, sidecar repository.Repository, repo repository.Repository) (sErr, rErr error) {
	if sidecar.Type() != repo.Type() {
		return fmt.Errorf("sidecar repo type '%s' does not match '%s'", sidecar.Type(), repo.Type()), nil
	}

	sErr = sidecar.Ping(ctx)
	if sErr != nil { // in case the sidecar fails, we test the repo
		logging.Debug("sidecar could not respond. error: %s", sErr.Error())
		rErr = repo.Ping(ctx)
		if rErr != nil {
			logging.Debug("repository could not respond. error: %s", rErr.Error())

		}
	}

	return sErr, rErr // returing both errors so that we can chech which failed
}

// HealthCheck performs the full health check procedure, including the retries.
func HealthCheck(ctx context.Context, cfg *config.LambdaConfig) error {

	sidecar, err := repository.Recover(cfg.Sidecar.RepoType)(ctx, cfg.Sidecar)
	if err != nil {
		return err
	}
	repo, err := repository.Recover(cfg.Repo.RepoType)(ctx, cfg.Repo)
	if err != nil {
		return err
	}

	var sErr, rErr error
	for i := 0; i < cfg.NumberOfRetries; i++ {
		sErr, rErr = singleHealthCheck(ctx, sidecar, repo)
		if sErr == nil { // if the sidecar responded without an error, no retries are needed
			logging.Info("sidecar healthy, returning health check as healthy")
			return nil
		}
	}
	if rErr != nil { // if both the sidecar and the repository are failing, we don't trigger the fail open
		logging.Info("both sidecar and repo unhealthy. returning health check as healthy")
		return nil
	}

	logging.Info("sidecar unhealthy but repo healthy. return health check as unhealthy")
	return sErr
}
