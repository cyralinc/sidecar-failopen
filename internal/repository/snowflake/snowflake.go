package snowflake

import (
	"context"
	"fmt"

	"github.com/cyralinc/sidecar-failopen/internal/config"
	"github.com/cyralinc/sidecar-failopen/internal/keys"
	"github.com/cyralinc/sidecar-failopen/internal/logging"
	"github.com/cyralinc/sidecar-failopen/internal/repository"
	"github.com/cyralinc/sidecar-failopen/internal/repository/genericsql"

	// Snowflake DB driver
	_ "github.com/snowflakedb/gosnowflake"
)

// Snowflake is the name registered by the DB driver.
const Snowflake = "snowflake"

type snowflakeRepository struct {
	// The majority of the repository.Repository functionality is delegated to
	// a generic SQL repository instance (genericSqlRepo).
	*genericsql.GenericSqlRepository
}

// *snowflakeRepository implements repository.Repository
var _ repository.Repository = (*snowflakeRepository)(nil)

func NewSnowflakeRepository(_ context.Context, cfg config.RepoConfig) (repository.Repository, error) {
	logging.Debug("instantiating snowflake repo at %s:%d/%s?role=%s&warehouse=%s&account=%s",
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.SnowflakeConfig.Role,
		cfg.SnowflakeConfig.Warehouse,
		cfg.SnowflakeConfig.Account,
	)

	connStr := fmt.Sprintf(
		"%s:%s@%s:%d/%s?role=%s&warehouse=%s&account=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.SnowflakeConfig.Role,
		cfg.SnowflakeConfig.Warehouse,
		cfg.SnowflakeConfig.Account,
	)

	sqlRepo, err := genericsql.NewGenericSqlRepository(
		cfg.RepoName,
		Snowflake,
		cfg.Database,
		connStr,
	)
	if err != nil {
		return nil, fmt.Errorf("could not instantiate generic sql repository: %w", err)
	}

	return &snowflakeRepository{GenericSqlRepository: sqlRepo}, nil
}

func (repo *snowflakeRepository) Ping(ctx context.Context) error {
	errChan := make(chan error)
	go func() {
		errChan <- repo.GenericSqlRepository.Ping(ctx)
	}()

	select {
	case err := <-errChan:
		logging.Debug("request complete")
		return err
	case <-ctx.Done():
		logging.Debug("deadline exceeded, returning error")
		return ctx.Err()
	}
}

func init() {
	repository.Register(keys.SnowflakeRepoKey, NewSnowflakeRepository)
}
