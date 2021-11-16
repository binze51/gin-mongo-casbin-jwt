package mongors

import (
	"testing"
	"time"

	"gitee.com/binze/binzekeji/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//mongo 5.0时 会返回具体字段错误
type User struct {
	Id       int                  `json:"id" bson:"_id,omitempty" `
	Age      int                  `json:"age"`
	Uid      string               `json:"uid,omitempty" bson:"uid,omitempty"`
	Pwd      string               `json:"pwd,omitempty" bson:"pwd,omitempty"`
	Name     string               `json:"name,omitempty" bson:"name,omitempty" jsonschema:"required,description=账号名称"`
	Phone    float32              `json:"phone,omitempty" bson:"phone,omitempty"`
	Remark   string               `json:"remark,omitempty" bson:"remark,omitempty"`
	RoleId   []int32              `json:"role,omitempty" bson:"role,omitempty"`
	Account  float64              `json:"account,omitempty" bson:"account,omitempty"`
	CreateAt time.Time            `json:"createat" bson:"createat,omitempty"`
	Method   string               `json:"method" bson:"method,omitempty"`
	Oid      primitive.ObjectID   `json:"oid,omitempty" bson:"oid,omitempty"`
	Oids     []primitive.ObjectID `json:"oids,omitempty" bson:"oids,omitempty"`
}

//TestReflectFromType -
func TestReflectFromType(t *testing.T) {
	u := &User{}
	flect := &Reflector{ExpandedStruct: true, RequiredFromJSONSchemaTags: true, AllowAdditionalProperties: true}
	sc := flect.Reflect(u)
	bts, _ := utils.JsonMarshal(&sc)
	var o map[string]interface{}
	_ = utils.JsonUnmarshal(bts, &o)
}
