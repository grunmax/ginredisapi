package main

import (
	"github.com/gin-gonic/gin"
	"github.com/grunmax/GinRedisApi/acs"
	"github.com/grunmax/GinRedisApi/utl"
	"net/http"
)

func makeRoutes() http.Handler {

	cors := func(c *gin.Context) {
		c.Writer.Header().Add("access-control-allow-origin", "*")
		c.Writer.Header().Add("access-control-allow-headers", "accept, content-type")
		c.Writer.Header().Add("access-control-allow-methods", "GET,HEAD,POST,DELETE,OPTIONS,PUT,PATCH")
	}

	//gin.SetMode(gin.ReleaseMode) // or debug

	routes := gin.Default()
	routes.Use(cors)
	routes.GET("/", func(c *gin.Context) {
		if item, err := acs.TestFunc(1, "gin&redis works", pool); err != nil {
			c.JSON(400, utl.BodyErr("redis r/w error"))
		} else {
			c.JSON(200, item)
		}
	})

	AddTodoRoutes(routes)

	return routes
}
