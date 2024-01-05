package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"io"
	"openai-sdk-experiment/api"
	"openai-sdk-experiment/docs"
	"os"
)

func route(engine *gin.Engine) {
	group := engine.Group("")
	sysSwaggerRouter(group)
	group.GET("/version", api.Version)    //版本
	group.GET("/health", api.HealthCheck) //健康检查

	v1 := group.Group("v1")
	v1.POST("/assistant", api.CreateAssistant)
	v1.POST("/thread/run", api.AddAndRun)
	v1.DELETE("/threads/:threadId", api.Clear)
	v1.GET("/threads/:threadId/messages", api.ListMessages)

}

// swagger路径
func sysSwaggerRouter(r *gin.RouterGroup) {
	version := ""
	f, err := os.Open("./VERSION")
	if err != nil {
		version = "latest"
	} else {
		apiVersion, _ := io.ReadAll(f)
		version = string(apiVersion)
	}
	docs.SwaggerInfo.Version = version
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func Server() {
	// 创建 Gin 路由器
	router := gin.Default()
	// 定义路由规则和处理函数
	route(router)
	// 启动 HTTP 服务器
	router.Run(":8820")
}
