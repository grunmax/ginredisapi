package acs

import (
	"github.com/garyburd/redigo/redis"
	"github.com/grunmax/GinRedisApi/dom"
	"github.com/grunmax/GinRedisApi/utl"
)

func TestFunc(id int, text string, pool *redis.Pool) (*dom.TestItem, error) {
	c := pool.Get()
	defer c.Close()

	redisKey := "key:test"
	if _, err := c.Do("SET", redisKey, text); err != nil {
		utl.Log("redis SET error", err)
		return nil, err
	}
	if redisValue, err := redis.String(c.Do("GET", redisKey)); err != nil {
		utl.Log("redis GET error", err)
		return nil, err
	} else {
		item := new(dom.TestItem)
		item.Id = id
		item.Text = redisValue
		return item, nil
	}
}
