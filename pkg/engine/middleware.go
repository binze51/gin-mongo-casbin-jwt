package engine

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"gitee.com/binze/binzekeji/pkg/response"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/ratelimit"
	"go.uber.org/zap"
)

// func (e *Engine) Swagger() gin.HandlerFunc {
// 	return ginSwagger.WrapHandler(swaggerFiles.Handler)
// }

//panic 恢复
func (e *Engine) RecoveryZap(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			// Check for a broken connection, as it is not really a
			// condition that warrants a panic stack trace.
			var brokenPipe bool
			if ne, ok := err.(*net.OpError); ok {
				if se, ok := ne.Err.(*os.SyscallError); ok {
					if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
						brokenPipe = true
					}
				}
			}

			httpRequest, _ := httputil.DumpRequest(c.Request, false)
			if brokenPipe {
				e.Log.Error(c.Request.URL.Path,
					zap.Any("error", err),
					zap.String("request", string(httpRequest)),
				)
				// If the connection is dead, we can't write a status to it.
				c.Error(err.(error)) // nolint: errcheck
				c.Abort()
				return
			}

			e.Log.Error("[Recovery from panic]",
				zap.Time("time", time.Now()),
				zap.Any("error", err),
				zap.String("request", string(httpRequest)),
				zap.String("stack", string(debug.Stack())),
			)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}()
	c.Next()
}

//日志器
func (e *Engine) LogZap(c *gin.Context) {
	start := time.Now()
	// some evil middlewares modify this values
	path := c.Request.URL.Path
	query := c.Request.URL.RawQuery
	c.Next()

	end := time.Now()
	latency := end.Sub(start)

	if len(c.Errors) > 0 {
		// Append error field if this is an erroneous request.
		for _, err := range c.Errors.Errors() {
			e.Log.Error(err)
		}
	} else {
		e.Log.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("time", end.Format(time.RFC3339)),
			zap.Duration("latency", latency),
		)
	}
}

//检查注入攻击
func (e *Engine) DenyInjection(c *gin.Context) {
	if c.GetHeader("Content-Type") == "application/json" {
		KindWord := []byte("$")
		safe := &io.LimitedReader{R: c.Request.Body, N: maxMemory}
		requestBody, _ := ioutil.ReadAll(safe)
		if bytes.Contains(requestBody, KindWord) {
			c.Abort()
			response.FailWithError(response.ParamBindError, "body 不能包含$字符", c)
			return
		}
		//避免 req.Body的sawEOF:true情况
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
		c.Next()
	}
	c.Next()
}

//权限检查,黑名单检查
func (e *Engine) AuthPromit(c *gin.Context) {
	var token string
	auth := c.GetHeader("Authorization")
	if auth == "" {
		e.Log.Error("authorization 参数缺失,请提交后重试")
		c.Abort()
		response.FailWithError(response.AuthorizationError, "authorization 参数缺失,请提交后重试", c)
		return
	}
	prefix := "Bearer "
	if strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	}
	denyColl := e.DBConn.Database("binze").Collection("gd_denytokenlist")
	err := denyColl.Find(context.Background(), bson.M{"useToken": token}).One(denyColl)
	if err == nil {
		e.Log.Error("该请求token已注销，请退出重新登录")
		c.Abort()
		response.FailWithError(response.AuthorizationError, "该请求token已注销，请重新登录", c)
		return
	}
	userid, role, err := e.Jwt.ParseUser(c, token)
	if err != nil {
		e.Log.Error(err.Error())
		c.Abort()
		if strings.Contains(err.Error(), "expired") {
			response.FailWithError(response.TokenExpiredError, "token已过期，请重新登录", c)
			return
		}
		response.FailWithError(response.AuthorizationError, "token解析失败，请重试", c)
		return
	}
	// Enforce执行(sub, obj, act)准入决策
	obj := c.Request.URL.Path
	act := c.Request.Method
	if ok, err := e.Ce.Enforce(role, obj, act); err != nil {
		e.Log.Error(err.Error())
		c.Abort()
		response.FailWithError(response.AuthorizationError, "鉴权失败，请重试", c)
		return
	} else if !ok {
		msg := fmt.Sprintf("抱歉！您对 %s 资源 暂无 %s 权限", obj, act)
		e.Log.Error(msg)
		c.Abort()
		response.FailWithError(response.AuthorizationError, msg, c)
		return
	}
	c.Set("userid", userid)
	c.Set("role", []string{role})
	c.Next()
}

//添加api记录
func (e *Engine) AddApiRow(path, method, group, description string) gin.HandlerFunc {
	apiColl := e.DBConn.Database(e.DataName).Collection("sys_api")
	ret, _ := apiColl.InsertOne(context.Background(), bson.M{
		"path":        path,
		"group":       group,
		"method":      method,
		"description": description,
	}) //path+method 要加组合唯一索引约束
	if ret.InsertedID == nil {
		//打印日志
		print("已存在")
	}
	return func(c *gin.Context) {
		c.Next()
	}
}

func (e *Engine) LeakBucket(limiter ratelimit.Limiter) gin.HandlerFunc {
	prev := time.Now()
	return func(ctx *gin.Context) {
		now := limiter.Take()
		log.Printf("%v", now.Sub(prev))
		prev = now
	}
}
