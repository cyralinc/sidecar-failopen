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

type MongoDBRepository struct {
	repoName string
	repoType string
	client   *mongo.Client
	database string
}

// *GenericSqlRepository implements repository.Repository
var _ repository.Repository = (*MongoDBRepository)(nil)

const connectionStringFmt string = "mongodb://%s:%s@%s:%d/%s"
const rsConnectionStringFmt string = "mongodb://%s:%s@%s/%s"

func NewMongoDBRepo(ctx context.Context, config config.RepoConfig) (repository.Repository, error) {
	logging.Debug("connecting to mongdb repo on mongdb://%s:%d", config.Host, config.Port)
	connStringOpts := util.ParseOptString(config)
	var connStr string
	if config.Port == 0 {
		connStr = fmt.Sprintf(rsConnectionStringFmt,
			config.User,
			config.Password,
			config.Host,
			connStringOpts,
		)
	} else {
		connStr = fmt.Sprintf(connectionStringFmt,
			config.User,
			config.Password,
			config.Host,
			config.Port,
			connStringOpts,
		)
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connStr))
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
