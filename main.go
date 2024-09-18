package main

import (
	"github.com/gin-gonic/gin"
	"github.com/prakyy/Go-Drive/controllers"
	"github.com/prakyy/Go-Drive/initializers"
	"github.com/prakyy/Go-Drive/middleware"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)
	r.POST("/upload", middleware.RequireAuth, controllers.UploadFile)
	r.GET("/files/:file_id", middleware.RequireAuth, controllers.GetFile)
	r.Run() // localhost:3000
}