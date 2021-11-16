package engine

import (
	"context"
	"net/http"
	"time"

	"gitee.com/binze/binzekeji/pkg/jwtauth"
	casbin "github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/qmgo"
	"go.uber.org/zap"
)

const (
	SlowThreshold       = time.Millisecond * 500 //慢日志阀值500毫秒
	maxMemory     int64 = 64 << 20
)

//构建全局引擎
func New(dataName string) *Engine {
	return &Engine{DataName: dataName}
}

type EngineConfig struct {
	DataName string `yaml:"dataname"`
	Rps      int    `yaml:"rps"`
}

//service 全局引擎
type Engine struct {
	DataName string //该服务持久化数据库
	DBConn   *qmgo.Client
	Log      *zap.Logger
	Ce       *casbin.Enforcer
	Jwt      *jwtauth.JWTAuth
}

func (e *Engine) GetDatabase() *qmgo.Database {
	return e.DBConn.Database(e.DataName)
}

//Close 回收资源
func (e *Engine) Close() error {
	//如果输出为文件，flush日志缓存落盘
	_ = e.Log.Sync()
	return e.DBConn.Close(context.Background())
}

//Healthz prometheus需要使用原生状态码
func (e *Engine) Healthz(c *gin.Context) {
	type res struct {
		Status string `json:"status"`
		Err    string `json:"err"`
	}
	err := e.DBConn.Ping(int64(time.Second) * 4)
	if err != nil {
		e.Log.Error("服务db资源并没有就绪", zap.Error(err))
		c.JSON(http.StatusBadRequest, res{Status: "no", Err: err.Error()})
	}
	c.JSON(http.StatusOK, res{Status: "ok"})
}
