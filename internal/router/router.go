package router

import (
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "time"

    "gdragon/internal/handler"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()

    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"}, 
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

    r.POST("/start", handler.StartTest)
    r.GET("/status", handler.TestStatus)
    r.GET("", handler.GetTests)
    r.GET("/tests", handler.GetPaginatedTests)

    return r
}
