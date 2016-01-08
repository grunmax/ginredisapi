package accs

import (
	"github.com/garyburd/redigo/redis"
	"github.com/grunmax/GinRedisApi/domn"
	"github.com/grunmax/GinRedisApi/util"
)

func TestFunc(id int, text string, pool *redis.Pool) (*domn.TestItem, error) {
	c := pool.Get()
	defer c.Close()

	redisKey := "key:test"
	if _, err := c.Do("SET", redisKey, text); err != nil {
		util.Log("redis SET error", err)
		return nil, err
	}
	if redisValue, err := redis.String(c.Do("GET", redisKey)); err != nil {
		util.Log("redis GET error", err)
		return nil, err
	} else {
		item := new(domn.TestItem)
		item.Id = id
		item.Text = redisValue
		return item, nil
	}
}

func RedisSet(prefix string, file []byte, pool *redis.Pool) (*domn.FormFile, error) {
	c := pool.Get()
	defer c.Close()
	id := util.NewId()
	key := prefix + id
	var ff domn.FormFile
	if _, err := c.Do("SET", key, file); err != nil {
		util.Log("key create error", err)
		ff.Id = ""
		return &ff, err
	} else {
		ff.Id = id
		return &ff, nil
	}
}

func RedisGet(prefix string, id string, pool *redis.Pool) ([]byte, error) {
	c := pool.Get()
	defer c.Close()
	key := prefix + id
	if file, err := redis.Bytes(c.Do("GET", key)); err != nil {
		util.Log("get by key error", err)
		return nil, err
	} else {
		return file, nil
	}
}
