package main

import (
	"blogapp.com/controllers"
	"blogapp.com/database"
	"blogapp.com/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// conntecting to db
	database.Connect()

	// Router's
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{
		"http://localhost:4200",
	}
	corsConfig.AllowMethods = []string{
		"GET", "POST", "PUT", "DELETE", "OPTIONS",
	}
	corsConfig.AllowHeaders = []string{
		"Origin", "Content-Type", "Authorization",
	}
	corsConfig.AllowCredentials = true

	r.Use(cors.New(corsConfig))

	// Blog CRUD route's
	r.GET("/blogs", controllers.GetAllBlogs)
	//r.GET("/blogs/:id", controllers.GetBlog)
	r.POST("/blogs", middleware.RequireAuth, controllers.CreateBlog)
	r.PUT("/blogs/:id", controllers.UpdateBlog)
	r.DELETE("/blogs/:id", controllers.DeleteBlog)
	r.GET("/blogs/:slug", controllers.GetBlogBySlugHandler)

	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)

	// starting server
	r.Run(":8080")
}
