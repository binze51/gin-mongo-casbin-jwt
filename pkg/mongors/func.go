package mongors

import (
	"context"
	"fmt"
	"log"
	"time"

	"gitee.com/binze/binzekeji/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// CheckSlowLogSave 慢op记录
func CheckSlowLogSave(startTime time.Duration, zlog *zap.Logger, collName, opMethod string, err1 error, docs ...interface{}) {
	duration := utils.TimeSince(startTime)
	if duration > SlowThreshold {
		content, _ := utils.JsonMarshal(docs)
		zlog.Error("mongo",
			zap.String("Mongo", fmt.Sprintf("[SLOW_LOG] mongo(%s) - slowcall - %s - fail(%s) - %s", collName, opMethod, err1.Error(), string(content))),
			zap.Duration("latency", duration),
		)
	}
}

// //正常返回
// return startSessCtx, err
func TransCtxPipline(startSessCtx mongo.SessionContext, NextColls *mongo.Collection, NextData interface{}) (mongo.SessionContext, error) {
	return nil, nil
}

//UpsertJsonSchemaValidator 更新表doc结构为严格模式
func UpsertJsonSchemaValidator(coll *mongo.Collection, collStruct interface{}, collname string) error {
	return createCollWithValidator(coll, collname, genJsonSchemaByPo(collStruct))
}

func createCollWithValidator(coll *mongo.Collection, collName string, validator map[string]interface{}) error {
	cursor, _ := coll.Database().ListCollections(context.Background(), bson.M{"name": collName})
	if !cursor.Next(context.Background()) {
		coll.Database().RunCommand(context.Background(), bson.D{{Key: "create", Value: collName}})
	}
	//bson严格模式
	result := coll.Database().RunCommand(context.Background(), bson.D{{Key: "collMod", Value: collName}, {Key: "validator", Value: validator}})
	if result.Err() == nil {
		log.Println(collName + ": ok---->")
		return nil
	}
	log.Println(collName + ": fail ---->")
	return result.Err()
}

//创建集合表验证schema结构对象
func genJsonSchemaByPo(po interface{}) (validator map[string]interface{}) {
	flect := &Reflector{ExpandedStruct: true, RequiredFromJSONSchemaTags: true, AllowAdditionalProperties: true}
	ob := flect.Reflect(po)
	bts, _ := utils.JsonMarshal(&ob)
	var o map[string]interface{}
	_ = utils.JsonUnmarshal(bts, &o)
	return bson.M{"$jsonSchema": o}
}
