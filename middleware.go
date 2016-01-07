package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	//	"github.com/patrickmn/go-cache"
)

func useMiddleware(api *gin.Engine) {
	//api.Use(DummyMiddleware)
	api.Use(CORSMiddlewareMiddleware)
	api.Use(CacheMiddleware)
}

func DummyMiddleware(c *gin.Context) {
	fmt.Println("dummy says url: " + c.Request.Method + c.Request.RequestURI)
	c.Next()
}

func CacheMiddleware(c *gin.Context) {
	if c.Request.Method == "GET" {
		item, isexist := che.Get(c.Request.RequestURI)
		if isexist {
			c.JSON(200, item)
			c.Abort()
		}
	}
	c.Next()
}

func CORSMiddlewareMiddleware(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, DELETE, OPTIONS")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	c.Next()
}
