package main

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/grunmax/GinRedisApi/accs"
	"github.com/grunmax/GinRedisApi/util"
)

func useMiddleware(api *gin.Engine) {
	//api.Use(DummyMiddleware)
	api.Use(CORSMiddlewareMiddleware)
	if cfg.AuthWorking {
		api.Use(AuthMiddleware)
		util.Log("Auth on", "")
	} else {
		util.Log("Auth off", "")
	}
	if cfg.CacheWorking {
		api.Use(CacheMiddleware)
		util.Log("Cache on", "")
	} else {
		util.Log("Cache off", "")
	}
}

func DummyMiddleware(c *gin.Context) {
	util.Log("dummy says: ", c.Request.Method+c.Request.RequestURI)
	c.Next()
}

func AuthMiddleware(c *gin.Context) {
	if c.Request.Method != "PUT" {
		cookievalue, err := util.GetCookieValue("k_yak", c)
		if err != nil {
			c.JSON(401, util.BodyErr("Access denied"))
			c.Abort()
		} else {
			if isexist, err := accs.KeyExistsCached(cookievalue, che, pool); err != nil || !isexist {
				c.JSON(401, util.BodyErr("Access denied"))
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
			//util.Log("cache!", c.Request.RequestURI)
			if strings.HasPrefix(c.Request.RequestURI, "/img/") {
				c.Data(200, "image/jpeg", item.([]byte))
				c.Abort()
			}
			if strings.HasPrefix(c.Request.RequestURI, "/file/") {
				c.Data(200, "application/octet-stream", item.([]byte))
				c.Abort()
			}
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
