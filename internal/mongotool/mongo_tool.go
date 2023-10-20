package mongotool

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/danenmao/pterergate-dtf/internal/config"
)

// 镜像安全 MongoDB 客户端
var s_DefaultMongoDB *mongo.Database

func GetDefaultMongoDB() *mongo.Database {
	return s_DefaultMongoDB
}

// 初始化客户端
func InitDefaultMongoClient() {

	uri := fmt.Sprintf("mongodb://%s:%s@%s/%s?authSource=admin&replicaSet=%s",
		config.DefaultMongoDB.Username,
		config.DefaultMongoDB.Password,
		config.DefaultMongoDB.Address,
		config.DefaultMongoDB.Database,
		config.DefaultMongoDB.ReplicaSet)

	client, err := NewMongoClient(uri)
	if err != nil {
		glog.Fatal("failed to new a mongodb client: ", err)
		return
	}

	if err = client.Ping(context.Background(), readpref.Primary()); err != nil {
		glog.Fatal("failed to connect to mongodb", err)
		return
	}

	s_DefaultMongoDB = client.Database(config.DefaultMongoDB.Database)
}

// 创建新的MongoDB客户端
func NewMongoClient(uri string) (*mongo.Client, error) {

	opts := options.Client()
	opts.ApplyURI(uri)

	ctx := context.Background()
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	return client, nil
}
