package sqlserver

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cyralinc/sidecar-failopen/internal/config"
	"github.com/cyralinc/sidecar-failopen/internal/keys"
	"github.com/cyralinc/sidecar-failopen/internal/repository"
	"github.com/cyralinc/sidecar-failopen/internal/repository/genericsql"

	// Sqlserver DB driver
	_ "github.com/denisenkom/go-mssqldb"
)

// SqlServer is the name registered by the DB driver.
const SqlServer = "sqlserver"

type sqlServerRepository struct {
	// The majority of the repository.Repository functionality is delegated to
	// a generic SQL repository instance (genericSqlRepo).
	*genericsql.GenericSqlRepository
}

// *sqlServerRepository implements repository.Repository
var _ repository.Repository = (*sqlServerRepository)(nil)

func NewSqlServerRepository(_ context.Context, cfg config.RepoConfig) (repository.Repository, error) {
	connStr := fmt.Sprintf(
		"sqlserver://%s:%s@%s:%s?database=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		strconv.Itoa(cfg.Port),
		cfg.Database,
	)

	genericSqlRepo, err := genericsql.NewGenericSqlRepository(
		cfg.RepoName,
		SqlServer,
		cfg.Database,
		connStr,
	)
	if err != nil {
		return nil, fmt.Errorf("could not instantiate generic sql repository: %w", err)
	}

	return &sqlServerRepository{GenericSqlRepository: genericSqlRepo}, nil
}

func init() {
	repository.Register(keys.SQLServerRepoKey, NewSqlServerRepository)
}
