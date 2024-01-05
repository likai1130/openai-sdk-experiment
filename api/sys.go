package api

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

// HealthCheck 健康检查
// @Summary 健康检查
// @Description 系统健康检查
// @Tags 系统服务
// @Accept  application/json
// @Product application/json
// @Success 200 {object} Response "{"code": 200, "data":{"status": "ok"} }"
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, &Response{
		Code: 0,
		Data: gin.H{"status": "ok"},
	})
}

// Version 系统版本
// @Summary 系统版本
// @Description 系统版本查询
// @Tags 系统服务
// @Accept  application/json
// @Product application/json
// @Success 200 {object} Response "{"code": 200, "data":{"version": "latest"} }"
// @Router /version [get]
func Version(c *gin.Context) {
	wd, _ := os.Getwd()
	fileName := "VERSION"
	version := "latest"
	fp := filepath.Join(wd, fileName)

	f, err := os.Open(fp)
	if os.IsNotExist(err) {
		create, _ := os.Create(fp)
		_, _ = create.Write([]byte(version))
	}

	file, err := io.ReadAll(f)
	if err == nil {
		version = string(file)
	}
	c.JSON(http.StatusOK, &Response{
		Code: 0,
		Data: gin.H{"version": version},
	})
}
