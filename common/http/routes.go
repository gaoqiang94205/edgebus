package http

import (
	"edgebus/server/pkg/handler"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRoutes() *gin.Engine {
	router := gin.New()
	//router.LoadHTMLGlob(config.GetEnv().TemplatePath + "/*")
	//if config.GetEnv().Debug {
	//	pprof.Register(router) // 性能分析工具
	//}

	router.Use(gin.Logger())
	//异常处理
	router.Use(handler.HandleErrors())

	//认证相关暂不做处理
	//router.Use(filters.RegisterSession()) // 全局session
	//router.Use(filters.RegisterCache())   // 全局cache
	//router.Use(auth.RegisterGlobalAuthDriver("cookie", "web_auth")) // 全局auth cookie
	//cookierouter.Use(auth.RegisterGlobalAuthDriver("jwt", "jwt_auth"))    // 全局auth jwt

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "找不到该路由",
		})
	})

	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "找不到该方法",
		})
	})
	// ReverseProxy
	// router.Use(proxy.ReverseProxy(map[string] string {
	// 	"localhost:4000" : "localhost:9090",
	// }))

	return router
}
