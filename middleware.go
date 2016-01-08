package main

import (
	"github.com/gin-gonic/gin"
	"github.com/grunmax/GinRedisApi/acs"
	"github.com/grunmax/GinRedisApi/utl"
)

func useMiddleware(api *gin.Engine) {
	api.Use(CORSMiddlewareMiddleware)
	if cfg.AuthWorking {
		api.Use(AuthMiddleware)
		utl.Log("Auth on", "")
	} else {
		utl.Log("Auth off", "")
	}
	if cfg.CacheWorking {
		api.Use(CacheMiddleware)
		utl.Log("Cache on", "")
	} else {
		utl.Log("Cache off", "")
	}
}

func DummyMiddleware(c *gin.Context) {
	utl.Log("dummy says: ", c.Request.Method+c.Request.RequestURI)
	c.Next()
}

func AuthMiddleware(c *gin.Context) {
	if c.Request.Method != "PUT" {
		cookievalue, err := utl.GetCookieValue("k_yak", c)
		if err != nil {
			c.JSON(401, utl.BodyErr("Access denied"))
			c.Abort()
		} else {
			if isexist, err := acs.KeyExistsCached(cookievalue, che, pool); err != nil || !isexist {
				c.JSON(401, utl.BodyErr("Access denied"))
				c.Abort()
				//c.AbortWithStatus(401)
			}
		}
	}
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
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, PUT, GET, DELETE")
	c.Next()
}
