package utl

import (
	"github.com/garyburd/redigo/redis"
	"github.com/nu7hatch/gouuid"
)

func InitRedisPool(url string, password string, maxConnections int) *redis.Pool {
	redisPool := redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", url)
		if err != nil {
			return nil, err
		}
		c.Do("AUTH", password)
		return c, err
	}, maxConnections)
	return redisPool
}

func NewId() string {
	id, _ := uuid.NewV4()
	return id.String()
}
