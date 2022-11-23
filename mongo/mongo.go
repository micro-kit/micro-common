package mongo

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/micro-kit/micro-common/config"
	"github.com/micro-kit/micro-common/etcdcli"
	"github.com/micro-kit/micro-common/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/* mongo客户端 */

var (
	// Client mongo 连接对象
	Client      *mongo.Client
	mongoConfig *config.MongoConfig
)

// GetClient 获取mongo连接对象
func GetClient() *mongo.Client {
	if Client == nil {
		log.Println("mongo 客户端为nil")
		initMongo()
	}
	return Client
}

// 获取db对象
func GetDatabase() *mongo.Database {
	return GetClient().Database(mongoConfig.Db)
}

func init() {
	initMongo()
}

// 初始化mongo
func initMongo() {
	var err error
	err = config.GetMongoConfig(etcdcli.EtcdCli, func(cfg *config.MongoConfig) {
		mongoConfig = cfg
		Client, err = NewClient(cfg)
		if err != nil {
			logger.Logger.Panicw("Creating mongo connection errors", "err", err)
		}
	})
	if err != nil {
		logger.Logger.Panicw("Get mongo configuration error", "err", err)
	}
}

// NewClient 创建客户端连接
func NewClient(cfg *config.MongoConfig) (client *mongo.Client, err error) {
	if cfg == nil {
		err = errors.New("the mongo configuration file can not be empty")
		return
	}
	logger.Logger.Infow("Start connecting to mongo database")

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	client, err = mongo.Connect(context.Background(), options.Client().
		// ApplyURI(strings.Join(cfg.Address, ",")).
		SetMaxPoolSize(uint64(cfg.MaxPoolSize)).
		SetMinPoolSize(uint64(cfg.MinPoolSize)).
		SetReplicaSet(cfg.ReplicaSet).
		SetHosts(cfg.Address).
		SetAuth(options.Credential{
			AuthMechanism: "SCRAM-SHA-256",
			Username:      cfg.Username,
			Password:      cfg.Password,
			AuthSource:    cfg.AuthSource,
		}))
	if err != nil {
		return nil, err
	}

	// ping 防止断开
	go func() {
		for {
			func() {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				err := client.Ping(ctx, nil)
				if err != nil {
					logger.Logger.Errorw("mongo ping error", "err", err)
				}
			}()

			time.Sleep(time.Second * 30)
		}
	}()
	logger.Logger.Infow("Connect to mongo database successfully")
	return
}
