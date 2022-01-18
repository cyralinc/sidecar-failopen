package redshift

import (
	"context"
	"fmt"

	"github.com/cyralinc/sidecar-failopen/internal/config"
	"github.com/cyralinc/sidecar-failopen/internal/keys"
	"github.com/cyralinc/sidecar-failopen/internal/logging"
	"github.com/cyralinc/sidecar-failopen/internal/repository"
	"github.com/cyralinc/sidecar-failopen/internal/repository/genericsql"
	"github.com/cyralinc/sidecar-failopen/internal/repository/postgresql/util"

	// Postgresql DB driver
	_ "github.com/lib/pq"
)

// PostgreSQL is the name registered by the DB driver.
const PostgreSQL = "postgres"

type redshiftRepository struct {
	// The majority of the repository.Repository functionality is delegated to
	// a generic SQL repository instance (genericSqlRepo).
	genericSqlRepo *genericsql.GenericSqlRepository
}

// *postgresqlRepository implements repository.Repository
var _ repository.Repository = (*redshiftRepository)(nil)

func NewRedshiftRepository(_ context.Context, cfg config.RepoConfig) (repository.Repository, error) {

	parsedOpts := util.ParseOptString(cfg)
	logging.Debug("using connection string opts: %s", parsedOpts)

	connStr := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		parsedOpts,
	)

	logging.Info("instantiating redshift repository at %s:%d", cfg.Host, cfg.Port)

	sqlRepo, err := genericsql.NewGenericSqlRepository(cfg.RepoName, PostgreSQL, cfg.Database, connStr)
	if err != nil {
		return nil, fmt.Errorf("could not instantiate generic sql repository: %w", err)
	}

	return &redshiftRepository{genericSqlRepo: sqlRepo}, nil
}

func (repo *redshiftRepository) Ping(ctx context.Context) error {
	return repo.genericSqlRepo.Ping(ctx)
}

func (repo *redshiftRepository) Close() error {
	return repo.genericSqlRepo.Close()
}

func (repo *redshiftRepository) Type() string {
	return keys.RedshiftRepoKey
}

func init() {
	repository.Register(keys.RedshiftRepoKey, NewRedshiftRepository)
}
