package snowflake

import (
	"context"
	"fmt"

	"github.com/cyralinc/sidecar-failopen/internal/config"
	"github.com/cyralinc/sidecar-failopen/internal/keys"
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
	connStr := fmt.Sprintf(
		"%s:%s@%s/%s?role=%s&warehouse=%s",
		cfg.User,
		cfg.Password,
		cfg.SnowflakeConfig.Account,
		cfg.Database,
		cfg.SnowflakeConfig.Role,
		cfg.SnowflakeConfig.Warehouse,
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

func init() {
	repository.Register(keys.SnowflakeRepoKey, NewSnowflakeRepository)
}
