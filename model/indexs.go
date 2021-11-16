package model

//coll索引文件
import (
	"context"
	"log"
	"reflect"

	"gitee.com/binze/binzekeji/pkg/mongors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UpsertJsonSchemaValidator(coll *mongo.Collection, model interface{}) error {
	refModel := reflect.ValueOf(model)
	GetCollNameMethod := refModel.MethodByName("GetCollName")
	collName := GetCollNameMethod.Call(make([]reflect.Value, 0))[0].String()

	return mongors.UpsertJsonSchemaValidator(coll, model, collName)
}

//UpsertCollIndex 仅测试用，生产上使用手工维护
func UpsertCollIndex(c *mongo.Collection) error {
	indexView := c.Indexes()
	_, err := indexView.CreateOne(context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "accounName", Value: 1}},
			Options: options.Index().SetUnique(true),
		})

	if err != nil {
		log.Println(err)
	}
	log.Println("order coll index update ok")
	return err
}
