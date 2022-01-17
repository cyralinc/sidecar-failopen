package mariadb

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

type mariaDBRepository struct {
	// The majority of the repository.Repository functionality is delegated to
	// a generic SQL repository instance (genericSqlRepo).
	genericSqlRepo *genericsql.GenericSqlRepository
}

// *mySqlRepository implements repository.Repository
var _ repository.Repository = (*mariaDBRepository)(nil)

func NewMariaDBRepository(_ context.Context, cfg config.RepoConfig) (repository.Repository, error) {
	connStr := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)
	logging.Info("instantiating mariadb repository at %s:%d", cfg.Host, cfg.Port)

	sqlRepo, err := genericsql.NewGenericSqlRepository(cfg.RepoName, MySQL, cfg.Database, connStr)
	if err != nil {
		return nil, fmt.Errorf("could not instantiate generic sql repository: %w", err)
	}

	return &mariaDBRepository{genericSqlRepo: sqlRepo}, nil
}

func (repo *mariaDBRepository) Ping(ctx context.Context) error {
	return repo.genericSqlRepo.Ping(ctx)
}

func (repo *mariaDBRepository) Close() error {
	return repo.genericSqlRepo.Close()
}

func (repo *mariaDBRepository) Type() string {
	return keys.MariaDBRepoKey
}

func init() {
	repository.Register(keys.MariaDBRepoKey, NewMariaDBRepository)
}
