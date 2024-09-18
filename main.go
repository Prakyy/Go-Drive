package main

import (
	"github.com/gin-gonic/gin"
	"github.com/prakyy/Go-Drive/controllers"
	"github.com/prakyy/Go-Drive/initializers"
	"github.com/prakyy/Go-Drive/middleware"
)

func init() {
	initializers.LoadEnv()
	initializers.LinkDB()
	initializers.SyncDB()
}

func main() {
	r := gin.Default()
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)
	r.POST("/upload", middleware.RequireAuth, controllers.UploadFile)
	r.GET("/files/:file_id", middleware.RequireAuth, controllers.GetFile)
	r.DELETE("/delete/:file_id", middleware.RequireAuth, controllers.DeleteFile)
	r.GET("/search", middleware.RequireAuth, controllers.SearchFile)
	r.Run() // localhost:3000
}