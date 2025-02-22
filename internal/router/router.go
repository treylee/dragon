package router

import (
	"github.com/gin-gonic/gin"
	"gdragon/internal/handler"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/start", handler.StartTest)
	r.GET("/check", handler.TestStatus)

	return r
}
