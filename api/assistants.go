package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"openai-sdk-experiment/service"
	"openai-sdk-experiment/service/assistants"
	"strconv"
)

// CreateAssistant 创建AI助手
// @Summary 创建AI助手
// @Description 创建AI助手
// @Tags AI 助手
// @Accept  application/json
// @Product application/json
// @Param data body assistants.AssistantRequest true "创建参数"
// @Success 200 {object} service.Response{} "{"code": 200, "data": [...]}"
// @Success 500 {object} service.Response
// @Router /v1/assistant [post]
func CreateAssistant(c *gin.Context) {
	var v assistants.AssistantRequest
	if err := c.ShouldBindJSON(&v); err != nil {
		c.JSON(http.StatusBadRequest, service.Response{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
		})
		return
	}
	newAssistants := service.NewAssistants()
	response := newAssistants.CreateAssistant(c.Request.Context(), v)
	c.JSON(response.Code, response)
	return
}

// AddAndRun 创建消息并且运行
// @Summary 创建消息并且运行
// @Description 创建消息并且运行
// @Tags AI 助手
// @Accept  application/json
// @Product application/json
// @Param data body service.RunRequest true "创建参数"
// @Success 200 {object} service.Response{} "{"code": 200, "data": [...]}"
// @Success 500 {object} service.Response
// @Router /v1/thread/run [post]
func AddAndRun(c *gin.Context) {
	var v service.RunRequest
	if err := c.ShouldBindJSON(&v); err != nil {
		c.JSON(http.StatusBadRequest, service.Response{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
		})
		return
	}

	ass := service.NewAssistants()
	response := ass.AddAndRun(c.Request.Context(), v)
	c.JSON(response.Code, response)
	return
}

// ListMessages 消息列表
// @Summary 消息列表
// @Description 消息列表
// @Tags AI 助手
// @Accept  application/json
// @Product application/json
// @Param threadId path string true "线程ID"
// @Param limit query int false "分页单位，不能超过一页100条"
// @Param order query string false "排序，aes,desc,默认desc"
// @Param after query string false "传message id，某个id之后"
// @Param before query string false "传message id，某个id之前"
// @Success 200 {object} service.Response{} "{"code": 200, "data": [...]}"
// @Success 500 {object} service.Response
// @Router /v1/threads/{threadId}/messages [get]
func ListMessages(c *gin.Context) {
	threadId := c.Param("threadId")
	after := c.Query("after")
	before := c.Query("before")
	pageTools := assistants.PageTools{}

	defaultLimit := c.DefaultQuery("limit", "10")
	defaultOrder := c.DefaultQuery("order", "desc")

	atoi, err := strconv.Atoi(defaultLimit)
	if err != nil {
		atoi = 10
	}
	pageTools.Limit = atoi
	pageTools.Order = defaultOrder

	if len(after) != 0 {
		pageTools.After = after
	}

	if len(before) != 0 {
		pageTools.Before = before
	}
	ass := service.NewAssistants()
	response := ass.ListMessages(c.Request.Context(), threadId, pageTools)
	c.JSON(response.Code, response)
	return
}

// Clear 删除线程
// @Summary 删除线程
// @Description 删除线程
// @Tags AI 助手
// @Accept  application/json
// @Product application/json
// @Param threadId path string true "线程ID"
// @Success 200 {object} service.Response{} "{"code": 200, "data": [...]}"
// @Success 500 {object} service.Response
// @Router /v1/threads/{threadId} [delete]
func Clear(c *gin.Context) {
	threadId := c.Param("threadId")
	ass := service.NewAssistants()
	response := ass.Clear(c.Request.Context(), threadId)
	c.JSON(response.Code, response)
	return
}
