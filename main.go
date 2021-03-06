package main

import (
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/grunmax/GinRedisApi/ctrl"
	"github.com/grunmax/GinRedisApi/midd"
	"github.com/grunmax/GinRedisApi/util"
	"github.com/patrickmn/go-cache"
)

var cfg *util.Config
var pool *redis.Pool
var che *cache.Cache

func init() {
	util.InitLog()
	cfg = util.ReadConfig()
	pool = util.InitRedisPool(cfg.RedisUrl, cfg.RedisPassword, cfg.MaxConnections)
	che = cache.New(time.Duration(cfg.CacheExpired)*time.Minute, time.Duration(cfg.CacheCheck)*time.Minute)
}

func main() {
	defer pool.Close()
	api := gin.Default()
	api.NoRoute(is404)
	//gin.SetMode(gin.ReleaseMode) // debug default
	midd.InitMiddleware(che, pool, cfg)
	midd.UseMiddleware(api)
	makeRoutes(api)
	api.Run(cfg.HttpUrl)
}

func makeRoutes(api *gin.Engine) {
	util.Log("Server start", gin.Mode())

	ctrl.AddTodoRoutes(che, pool, api)
	ctrl.AddUserRoutes(che, pool, api)
	ctrl.AddTestRoutes(che, pool, api)
}

func is404(c *gin.Context) {
	c.JSON(404, util.BodyErr("wrong URL"))
}
