package postgresql

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

type postgresqlRepository struct {
	// The majority of the repository.Repository functionality is delegated to
	// a generic SQL repository instance (genericSqlRepo).
	*genericsql.GenericSqlRepository
}

// *postgresqlRepository implements repository.Repository
var _ repository.Repository = (*postgresqlRepository)(nil)

// NewPostgresqlRepository creates a creator for a repo with postgresql driver with the given
// repo type
func NewPostgresqlRepository(repoType string) func(context.Context, config.RepoConfig) (repository.Repository, error) {
	return func(_ context.Context, cfg config.RepoConfig) (repository.Repository, error) {
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

		logging.Info("instantiating %s repository at %s:%d", repoType, cfg.Host, cfg.Port)

		sqlRepo, err := genericsql.NewGenericSqlRepository(cfg.RepoName, PostgreSQL, cfg.Database, connStr)
		if err != nil {
			return nil, fmt.Errorf("could not instantiate generic sql repository: %w", err)
		}

		return &postgresqlRepository{GenericSqlRepository: sqlRepo}, nil
	}
}

func init() {
	// registering the constructors to factory method. This will be run on import
	repository.Register(keys.PGRepoKey, NewPostgresqlRepository(keys.PGRepoKey))
	repository.Register(keys.DenodoRepoKey, NewPostgresqlRepository(keys.DenodoRepoKey))
	repository.Register(keys.RedshiftRepoKey, NewPostgresqlRepository(keys.RedshiftRepoKey))
}
