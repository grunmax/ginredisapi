package accs

import (
	"github.com/garyburd/redigo/redis"
	"github.com/grunmax/GinRedisApi/domn"
	"github.com/grunmax/GinRedisApi/util"
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
			util.Log("SCAN read error", err)
			return nil, err
		}

		values, err = redis.Scan(values, &cursor, &items)
		if err != nil {
			util.Log("SCAN scan error", err)
			return nil, err
		}

		results = append(results, items...)
		if cursor == 0 {
			break
		}
	}
	return results, nil
}

func TodoCreate(item domn.TodoForm, pool *redis.Pool) (*domn.TodoForm, error) {
	c := pool.Get()
	defer c.Close()
	item.Id = util.NewId()
	key := "todo:" + item.Id
	if _, err := c.Do("HMSET", key,
		"title", item.Title,
		"completed", item.Completed,
		"order", item.Order); err != nil {
		util.Log("Todo create error", err)
		return nil, err
	} else {
		return &item, nil
	}
}

func TodoEdit(id string, item domn.TodoForm, pool *redis.Pool) (*domn.TodoForm, error) {
	c := pool.Get()
	defer c.Close()
	item.Id = id
	key := "todo:" + item.Id
	if exist, err := redis.Bool(c.Do("EXISTS", key)); err != nil {
		util.Log("check todo key error", err)
		return nil, err
	} else if !exist {
		return nil, util.MyErr{"no key"}
	}
	if _, err := c.Do("HMSET", key,
		"title", item.Title,
		"completed", item.Completed,
		"order", item.Order); err != nil {
		util.Log("Todo edit error", err)
		return nil, err
	} else {
		return &item, nil
	}
}

func TodoGetId(id string, pool *redis.Pool) (*domn.TodoItem, error) {
	c := pool.Get()
	defer c.Close()
	key := "todo:" + id

	values, err := redis.Values(c.Do("HGETALL", key))
	if err != nil {
		util.Log("HGET key error", err)
		return nil, err
	}
	if len(values) == 0 {
		return nil, util.MyErr{"no data"}
	}

	var todo domn.TodoItem
	if err := redis.ScanStruct(values, &todo); err != nil {
		util.Log("HGET parse key error", err)
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
		util.Log("Todo delete error", err)
		return err
	} else {
		return nil
	}

}
