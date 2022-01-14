package postgresql

import (
	"context"
	"fmt"

	"github.com/cyralinc/cloudformation-sidecar-failopen/internal/config"
	"github.com/cyralinc/cloudformation-sidecar-failopen/internal/keys"
	"github.com/cyralinc/cloudformation-sidecar-failopen/internal/logging"
	"github.com/cyralinc/cloudformation-sidecar-failopen/internal/repository"
	"github.com/cyralinc/cloudformation-sidecar-failopen/internal/repository/genericsql"

	// Postgresql DB driver
	_ "github.com/lib/pq"
)

// PostgreSQL is the name registered by the DB driver.
const PostgreSQL = "postgres"

type postgresqlRepository struct {
	// The majority of the repository.Repository functionality is delegated to
	// a generic SQL repository instance (genericSqlRepo).
	genericSqlRepo *genericsql.GenericSqlRepository
}

// *postgresqlRepository implements repository.Repository
var _ repository.Repository = (*postgresqlRepository)(nil)

func NewPostgresqlRepository(_ context.Context, cfg config.RepoConfig) (repository.Repository, error) {

	connStr := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	logging.Info("instantiating postgres repository at %s:%d", cfg.Host, cfg.Port)

	sqlRepo, err := genericsql.NewGenericSqlRepository(cfg.RepoName, PostgreSQL, cfg.Database, connStr)
	if err != nil {
		return nil, fmt.Errorf("could not instantiate generic sql repository: %w", err)
	}

	return &postgresqlRepository{genericSqlRepo: sqlRepo}, nil
}

func (repo *postgresqlRepository) Ping(ctx context.Context) error {
	return repo.genericSqlRepo.Ping(ctx)
}

func (repo *postgresqlRepository) Close() error {
	return repo.genericSqlRepo.Close()
}

func (repo *postgresqlRepository) Type() string {
	return keys.PGRepoKey
}

func init() {
	repository.Register(keys.PGRepoKey, NewPostgresqlRepository)
}
