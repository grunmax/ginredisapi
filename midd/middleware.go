package midd

import (
	"strings"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/grunmax/GinRedisApi/accs"
	"github.com/grunmax/GinRedisApi/util"
	"github.com/patrickmn/go-cache"
)

var pool *redis.Pool
var che *cache.Cache
var cfg *util.Config

func InitMiddleware(che_ *cache.Cache, pool_ *redis.Pool, cfg_ *util.Config) {
	pool = pool_
	che = che_
	cfg = cfg_
}

func UseMiddleware(api *gin.Engine) {
	//api.Use(dummyMiddleware)
	api.Use(corsMiddlewareMiddleware)
	if cfg.AuthWorking {
		api.Use(authMiddleware)
		util.Log("Auth on", "")
	} else {
		util.Log("Auth off", "")
	}
	if cfg.CacheWorking {
		api.Use(cacheMiddleware)
		util.Log("Cache on", "")
	} else {
		util.Log("Cache off", "")
	}
}

func dummyMiddleware(c *gin.Context) {
	util.Log("dummy says: ", c.Request.Method+c.Request.RequestURI)
	c.Next()
}

func authMiddleware(c *gin.Context) {
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

func cacheMiddleware(c *gin.Context) {
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

func corsMiddlewareMiddleware(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, PUT, GET, DELETE")
	c.Next()
}
