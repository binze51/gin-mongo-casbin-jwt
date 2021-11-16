package router

import "github.com/gin-gonic/gin"

type Router struct {
	*gin.Engine
}

//New全局路由
func New() *Router {
	Router := &Router{}
	Router.Engine = gin.New()
	return Router
}
