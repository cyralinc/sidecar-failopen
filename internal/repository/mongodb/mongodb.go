package mongodb

import (
	"context"
	"fmt"

	"github.com/cyralinc/sidecar-failopen/internal/config"
	"github.com/cyralinc/sidecar-failopen/internal/keys"
	"github.com/cyralinc/sidecar-failopen/internal/logging"
	"github.com/cyralinc/sidecar-failopen/internal/repository"
	"github.com/cyralinc/sidecar-failopen/internal/repository/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoDBRepository struct {
	repoName string
	repoType string
	client   *mongo.Client
	database string
}

// *mongoDBRepository implements repository.Repository
var _ repository.Repository = (*mongoDBRepository)(nil)

// connection string following: https://docs.mongodb.com/v5.0/reference/connection-string/
const connectionStringFmt string = "mongodb://%s:%s@%s:%d/%s"

func NewMongoDBRepo(ctx context.Context, config config.RepoConfig) (repository.Repository, error) {
	logging.Debug("connecting to mongodb repo on %s:%d", config.Host, config.Port)
	connStringOpts := util.ParseOptString(config)
	var connStr string
	connStr = fmt.Sprintf(connectionStringFmt,
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
	return &mongoDBRepository{
		repoName: config.RepoName,
		client:   client,
		repoType: keys.MongoDBRepoKey,
	}, nil
}

func (repo *mongoDBRepository) Ping(ctx context.Context) error {
	return repo.client.Ping(ctx, readpref.Primary())
}

func (repo *mongoDBRepository) Close() error {
	return repo.client.Disconnect(context.TODO())
}
func (repo *mongoDBRepository) Type() string {
	return repo.repoType
}

func init() {
	repository.Register(keys.MongoDBRepoKey, NewMongoDBRepo)
}
