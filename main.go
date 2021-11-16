package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"gitee.com/binze/binzekeji/api"
	"gitee.com/binze/binzekeji/config"
	"gitee.com/binze/binzekeji/pkg/casbinx"
	"gitee.com/binze/binzekeji/pkg/engine"
	"gitee.com/binze/binzekeji/pkg/httpserver"
	"gitee.com/binze/binzekeji/pkg/jwtauth"
	"gitee.com/binze/binzekeji/pkg/logger"
	"gitee.com/binze/binzekeji/pkg/mongors"
	"gitee.com/binze/binzekeji/pkg/router"
	"gitee.com/binze/binzekeji/pkg/shutdown"
	"gitee.com/binze/binzekeji/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/ratelimit"

	"go.uber.org/zap"
)

func main() {
	//环境变量注入后执行：export envflag=dev && go run -race main.go
	envflag, ok := os.LookupEnv("envflag")
	if !ok {
		panic("envflag环境变量缺失，请设置envflag后重试")
	}
	//构造依赖对象
	conf := config.FromYamlFile(envflag)
	//log记录器
	zlog := logger.NewJSONLogger(&conf.ZapConfig)
	//mongo 副本集
	mongors := mongors.New(&conf.MongoRSConfig)
	//jwt
	jwt := jwtauth.New(&conf.JwtConfig)
	//casbin
	ca := casbinx.NewAdapter(&conf.MongoAdapterConf)

	//构造全家engine，注入依赖对象
	_engine := engine.New(conf.EngineConfig.DataName)
	_engine.Log = zlog
	_engine.Ce = ca
	_engine.Jwt = jwt
	_engine.DBConn = mongors

	//初始化服务模块资源对象
	service.InitService(_engine, conf)

	//全局route
	_router := router.New()
	//流量限速
	limitrate := ratelimit.New(conf.EngineConfig.Rps)
	_router.Use(_engine.LeakBucket(limitrate))
	// //go运行时指标
	// pprof.RouteRegister(_router)
	//日志，panic拦截并记录
	_router.Use(_engine.LogZap, _engine.RecoveryZap)
	//检查mongo js注入攻击
	_router.Use(_engine.DenyInjection)

	//注册系统所有模块路由
	api.RegisterRouter(_router)

	//http服务
	httpServer := httpserver.New(&conf.HttpServConfig, _router)
	fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] RATE"+fmt.Sprintf("   Current Rate Limit: %d requests/s", conf.Rps))
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zlog.Fatal("http server startup err", zap.Error(err))
		}
	}()

	//释放资源
	shutdown.New().Close(
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()
			//优雅释放h2Server连接资源
			if err := httpServer.Shutdown(ctx); err != nil {
				zlog.Error("server shutdown fail", zap.Error(err))
			} else {
				zlog.Info("server shutdown success")
			}
		},
		func() {
			//关闭服务的数据库资源,日志缓存flush
			if err := _engine.Close(); err != nil {
				zlog.Error("svcResources release fail", zap.Error(err))
			} else {
				zlog.Info("svcResources release success")
			}
		},
	)
}
