package casbinx

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	mongodbadapter "github.com/casbin/mongodb-adapter/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//MongoAdapterConf
type MongoAdapterConf struct {
	MongoURI          string `yaml:"mongouri"`
	DataBase          string `yaml:"database"`
	CollName          string `yaml:"collname"`
	RBACModelConfFile string `yaml:"modelfile"`
}

//NewAdapter 新建适配器
func NewAdapter(conf *MongoAdapterConf) *casbin.Enforcer {
	mongoClientOption := options.Client().ApplyURI(conf.MongoURI)
	a, err := mongodbadapter.NewAdapterWithCollectionName(mongoClientOption, conf.DataBase, conf.CollName)
	if err != nil {
		panic(err)
	}
	m, err := model.NewModelFromString(conf.RBACModelConfFile)
	if err != nil {
		panic(err)
	}

	//新建Enforcer
	e, err := casbin.NewEnforcer(m, a)
	if err != nil {
		panic(err)
	}
	err = e.LoadPolicy()
	if err != nil {
		panic(err)
	}
	return e
}
