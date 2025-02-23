package router

import (
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "gdragon/internal/handler"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()

    r.Use(cors.Default())

    r.POST("/start", handler.StartTest)
    r.GET("/check", handler.TestStatus)

    return r
}
