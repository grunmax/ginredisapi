package main

import (
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/grunmax/GinRedisApi/utl"
)

var cfg *utl.Config
var pool *redis.Pool

func init() {
	utl.InitLog()
	cfg = utl.ReadConfig()
	pool = utl.InitRedisPool(cfg.RedisUrl, cfg.RedisPassword, cfg.MaxConnections)
}

func main() {
	defer pool.Close()
	api := gin.Default()
	//gin.SetMode(gin.ReleaseMode) // debug default
	useMiddleware(api)
	makeRoutes(api)
	api.Run(cfg.HttpUrl)
}
