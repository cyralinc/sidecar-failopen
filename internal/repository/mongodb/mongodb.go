package mongodb

import (
	"context"
	"fmt"

	"github.com/cyralinc/sidecar-failopen/internal/config"
	"github.com/cyralinc/sidecar-failopen/internal/logging"
	"github.com/cyralinc/sidecar-failopen/internal/repository"
	"github.com/cyralinc/sidecar-failopen/internal/repository/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	PingQuery = "SELECT 1"
)

/*
MongoDBRepository is an implementation of repository.Repository that
works for a subset of ANSI SQL compatible databases. In addition to the
standard repository.Repository methods, it also exposes some SQL-specific
functionality, which may be useful for other repository.Repository
implementations.
*/
type MongoDBRepository struct {
	repoName string
	repoType string
	client   *mongo.Client
	database string
}

// *GenericSqlRepository implements repository.Repository
var _ repository.Repository = (*MongoDBRepository)(nil)

const connectionStringFmt string = "mongodb://%s:%s@%s:%d/%s"

/*
NewGenericSqlRepository is the constructor for the GenericSqlRepository type.
Note that it returns a pointer to GenericSqlRepository rather than
repository.Repository. This is intentional, as the GenericSqlRepository type
exposes additional functionality on top of the repository.Repository
interface. If the caller is not concerned with this additional functionality,
they are free to assign to the return value to repository.Repository
*/
func NewMongoDBRepo(ctx context.Context, config config.RepoConfig) (repository.Repository, error) {
	logging.Debug("connecting to mongdb repo on mongdb://%s:%d", config.Host, config.Port)
	connStringOpts := util.ParseOptString(config)
	connStr := fmt.Sprintf(connectionStringFmt,
		config.User,
		config.Password,
		config.Host,
		config.Port,
		connStringOpts,
	)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connStr))
	if err != nil {
		return nil, fmt.Errorf("error connecting to mongodb repo @ %s:%d - error: %v", config.Host, config.Port, err)
	}
	return &MongoDBRepository{
		repoName: config.RepoName,
		client:   client,
		repoType: "mongodb",
	}, nil
}

func (repo *MongoDBRepository) Ping(ctx context.Context) error {
	return repo.client.Ping(ctx, readpref.Primary())
}

func (repo *MongoDBRepository) Close() error {
	return repo.client.Disconnect(context.TODO())
}
func (repo *MongoDBRepository) Type() string {
	return repo.repoType
}

func init() {
	repository.Register("mongodb", NewMongoDBRepo)
}
