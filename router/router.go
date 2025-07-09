package router

import (
	"filesys/endpoint"
	"filesys/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/login", endpoint.Login)

	auth := r.Group("/api", middleware.AuthMiddleware())
	auth.POST("/user", endpoint.CreateUser)

	fileGroup := r.Group("/api/file", middleware.AuthMiddleware(), middleware.FilePermissionMiddleware())
	fileGroup.POST("/:file_id/new", endpoint.CreateFolder)

	return r
}
