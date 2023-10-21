package main

import (
	"todo-dodo/api"
	"todo-dodo/db"
	"todo-dodo/pages"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func main() {
	// Init db
	db.Open()
	defer db.Close()

	// Setup router
	engine := gin.Default()

	// Middleware
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	engine.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedPaths([]string{"/api/**"})))

	{ // JSON-based APIs
		api_g := engine.Group("/api/v1")
		api_g.POST("/task/create", api.TaskCreateBatchEnpoint)
		api_g.POST("/task/fetch", api.TaskFetchAllEnpoint)
		api_g.POST("/task/delete", api.TaskDeleteBatchEnpoint)
		api_g.POST("/task/update", api.TaskUpdateEndpoint)

		api_g.POST("/tag/delete", api.TagDeleteBatchEnpoint)
	}

	// Templates
	engine.LoadHTMLGlob("./templates/*")

	{ // Js files
		js_g := engine.Group("/js")
		js_g.StaticFile("/index.js", "./pages/index.js")
	}

	{
		css_g := engine.Group("/css")
		css_g.StaticFile("/index.css", "./pages/index.css")
	}

	{ // Website
		site_g := engine.Group("")
		site_g.GET("/", pages.Index)
	}

	// Run server
	engine.Run()
}
