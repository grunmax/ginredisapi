package accs

import (
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/grunmax/GinRedisApi/domn"
	"github.com/grunmax/GinRedisApi/util"
	"github.com/patrickmn/go-cache"
)

func UserCreate(item domn.UserSignupForm, pool *redis.Pool) (string, error) {
	c := pool.Get()
	defer c.Close()
	item.Id = util.NewId()
	hash := util.GetMD5Hash(item.Password)
	key := "user:" + hash
	if _, err := c.Do("HMSET", key, "id", item.Id, "email", item.Email); err != nil {
		util.Log("User create error", err)
		return "", err
	} else {
		util.Log("hash", hash)
		return hash, nil
	}
}

func UserCheck(item domn.UserLoginForm, pool *redis.Pool) (bool, string, error) {
	c := pool.Get()
	defer c.Close()
	hash := util.GetMD5Hash(item.Password)
	key := "user:" + hash
	email, err := redis.String(c.Do("HGET", key, "email"))
	if err != nil {
		util.Log("user check error", err)
		return false, "", err
	}
	return email == item.Email, hash, nil
}

func KeyExists(hash string, pool *redis.Pool) (bool, error) {
	c := pool.Get()
	defer c.Close()
	key := "user:" + hash
	exist, err := redis.Bool(c.Do("EXISTS", key))
	if err != nil {
		util.Log("check todo key error", err)
		return false, err
	}
	return exist, nil
}

func KeyExistsCached(hash string, che *cache.Cache, pool *redis.Pool) (bool, error) {
	_, isexist := che.Get(hash)
	if isexist {
		return true, nil
	}
	c := pool.Get()
	defer c.Close()
	key := "user:" + hash
	exist, err := redis.Bool(c.Do("EXISTS", key))
	if err != nil {
		util.Log("check todo key error", err)
		return false, err
	}
	if exist {
		che.Set(hash, "foo", 10*time.Minute)
	}
	return exist, nil
}
