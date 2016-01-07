package acs

import (
	"github.com/garyburd/redigo/redis"
	"github.com/grunmax/GinRedisApi/dom"
	"github.com/grunmax/GinRedisApi/utl"
)

func GetKeys(match string, pool *redis.Pool) ([]string, error) {
	c := pool.Get()
	defer c.Close()
	var items []string
	const step = 1000

	results := make([]string, 0)
	cursor := 0
	for {
		values, err := redis.Values(c.Do("SCAN", cursor, "MATCH", match, "COUNT", step))
		if err != nil {
			utl.Log("SCAN read error", err)
			return nil, err
		}

		values, err = redis.Scan(values, &cursor, &items)
		if err != nil {
			utl.Log("SCAN scan error", err)
			return nil, err
		}

		results = append(results, items...)
		if cursor == 0 {
			break
		}
	}
	return results, nil
}

func TodoCreate(item dom.TodoForm, pool *redis.Pool) (*dom.TodoForm, error) {
	c := pool.Get()
	defer c.Close()
	item.Id = utl.NewId()
	key := "todo:" + item.Id
	if _, err := c.Do("HMSET", key,
		"title", item.Title,
		"completed", item.Completed,
		"order", item.Order); err != nil {
		utl.Log("Todo create error", err)
		return nil, err
	} else {
		return &item, nil
	}
}

func TodoEdit(id string, item dom.TodoForm, pool *redis.Pool) (*dom.TodoForm, error) {
	c := pool.Get()
	defer c.Close()
	item.Id = id
	key := "todo:" + item.Id
	if exist, err := redis.Bool(c.Do("EXISTS", key)); err != nil {
		utl.Log("check todo key error", err)
		return nil, err
	} else if !exist {
		return nil, utl.MyErr{"no key"}
	}
	if _, err := c.Do("HMSET", key,
		"title", item.Title,
		"completed", item.Completed,
		"order", item.Order); err != nil {
		utl.Log("Todo edit error", err)
		return nil, err
	} else {
		return &item, nil
	}
}

func TodoGetId(id string, pool *redis.Pool) (*dom.TodoItem, error) {
	c := pool.Get()
	defer c.Close()
	key := "todo:" + id

	values, err := redis.Values(c.Do("HGETALL", key))
	if err != nil {
		utl.Log("HGET key error", err)
		return nil, err
	}
	if len(values) == 0 {
		return nil, utl.MyErr{"no data"}
	}

	var todo dom.TodoItem
	if err := redis.ScanStruct(values, &todo); err != nil {
		utl.Log("HGET parse key error", err)
		return nil, err
	}
	todo.Id = id
	return &todo, nil
}

func TodoDeleteId(id string, pool *redis.Pool) error {
	c := pool.Get()
	defer c.Close()
	key := "todo:" + id
	if _, err := c.Do("DEL", key); err != nil {
		utl.Log("Todo delete error", err)
		return err
	} else {
		return nil
	}

}
