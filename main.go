package main

import (
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/grunmax/GinRedisApi/utl"
	"github.com/patrickmn/go-cache"
)

var cfg *utl.Config
var pool *redis.Pool
var che *cache.Cache

func init() {
	utl.InitLog()
	cfg = utl.ReadConfig()
	pool = utl.InitRedisPool(cfg.RedisUrl, cfg.RedisPassword, cfg.MaxConnections)
	che = cache.New(time.Duration(cfg.CacheExpired)*time.Minute, time.Duration(cfg.CacheCheck)*time.Minute)
}

func main() {
	defer pool.Close()
	api := gin.Default()
	api.NoRoute(is404)
	//gin.SetMode(gin.ReleaseMode) // debug default
	useMiddleware(api)
	makeRoutes(api)
	api.Run(cfg.HttpUrl)
}
