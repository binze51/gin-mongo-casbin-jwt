package mongors

import (
	"context"
	"runtime"
	"time"

	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	mop "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

const (
	SlowThreshold = time.Millisecond * 500 //慢日志阀值
)

type MongoRSConfig struct {
	Uri      string `yaml:"uri"`
	DataBase string `yaml:"database"`
	CollInit bool   `yaml:"collinit"`
}

//New 构建mongo 客户端
func New(conf *MongoRSConfig) *qmgo.Client {
	ops := options.ClientOptions{
		//注意这里指针 要new初始化空间
		ClientOptions: new(mop.ClientOptions),
	}
	p := uint64(runtime.NumCPU() * 2)
	minp := uint64(4)
	ops.MaxPoolSize = &p
	ops.MinPoolSize = &minp
	ops.WriteConcern = writeconcern.New(writeconcern.J(true), writeconcern.W(1))
	ops.ReadPreference = readpref.PrimaryPreferred()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := qmgo.NewClient(ctx, &qmgo.Config{
		Uri: conf.Uri,
	}, ops)
	if err != nil {
		panic(err)
	}
	err = client.Ping(int64(time.Second * 5))
	if err != nil {
		panic(err)
	}
	return client
}
