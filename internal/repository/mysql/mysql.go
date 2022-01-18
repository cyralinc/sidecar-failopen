package mysql

import (
	"context"
	"fmt"

	"github.com/cyralinc/sidecar-failopen/internal/config"
	"github.com/cyralinc/sidecar-failopen/internal/keys"
	"github.com/cyralinc/sidecar-failopen/internal/logging"
	"github.com/cyralinc/sidecar-failopen/internal/repository"
	"github.com/cyralinc/sidecar-failopen/internal/repository/genericsql"
	_ "github.com/go-sql-driver/mysql"
)

// MySQL is the name registered by the driver.
const MySQL = "mysql"

type mySqlRepository struct {
	// The majority of the repository.Repository functionality is delegated to
	// a generic SQL repository instance (genericSqlRepo).
	genericSqlRepo *genericsql.GenericSqlRepository
}

// *mySqlRepository implements repository.Repository
var _ repository.Repository = (*mySqlRepository)(nil)

func NewMySQLRepository(repoType string) func(_ context.Context, cfg config.RepoConfig) (repository.Repository, error) {
	return func(_ context.Context, cfg config.RepoConfig) (repository.Repository, error) {
		connStr := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.Database,
		)
		logging.Info("instantiating %s repository at %s:%d", repoType, cfg.Host, cfg.Port)

		sqlRepo, err := genericsql.NewGenericSqlRepository(cfg.RepoName, MySQL, cfg.Database, connStr)
		if err != nil {
			return nil, fmt.Errorf("could not instantiate generic sql repository: %w", err)
		}

		return &mySqlRepository{genericSqlRepo: sqlRepo}, nil
	}
}

func (repo *mySqlRepository) Ping(ctx context.Context) error {
	return repo.genericSqlRepo.Ping(ctx)
}

func (repo *mySqlRepository) Close() error {
	return repo.genericSqlRepo.Close()
}

func (repo *mySqlRepository) Type() string {
	return keys.MySQLRepoKey
}

func init() {
	// registering the constructors to factory method. This will be run on import
	repository.Register(keys.MySQLRepoKey, NewMySQLRepository(keys.MySQLRepoKey))
	repository.Register(keys.MariaDBRepoKey, NewMySQLRepository(keys.MariaDBRepoKey))
}
