package main

import (
	"github.com/garyburd/redigo/redis"
	"github.com/grunmax/GinRedisApi/utl"
	"net/http"
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
	routes := makeRoutes()
	http.ListenAndServe(cfg.HttpUrl, routes)
}
