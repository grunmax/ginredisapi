package acs

import (
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/grunmax/GinRedisApi/dom"
	"github.com/grunmax/GinRedisApi/utl"
	"github.com/patrickmn/go-cache"
)

func UserCreate(item dom.UserSignupForm, pool *redis.Pool) (string, error) {
	c := pool.Get()
	defer c.Close()
	item.Id = utl.NewId()
	hash := utl.GetMD5Hash(item.Password)
	key := "user:" + hash
	if _, err := c.Do("HMSET", key, "id", item.Id, "email", item.Email); err != nil {
		utl.Log("User create error", err)
		return "", err
	} else {
		utl.Log("hash", hash)
		return hash, nil
	}
}

func UserCheck(item dom.UserLoginForm, pool *redis.Pool) (bool, string, error) {
	c := pool.Get()
	defer c.Close()
	hash := utl.GetMD5Hash(item.Password)
	key := "user:" + hash
	email, err := redis.String(c.Do("HGET", key, "email"))
	if err != nil {
		utl.Log("user check error", err)
		return false, "", err
	}
	return email == item.Email, hash, nil
}

//func UserCheckCached(item dom.UserLoginForm, che *cache.Cache, pool *redis.Pool) (bool, string, error) {
//	hash := utl.GetMD5Hash(item.Password)
//	_, isexist := che.Get(hash)
//		if isexist {
//			utl.Log("user check auth cache!", hash)
//			true, hash, nil
//		}
//	}
//	c := pool.Get()
//	defer c.Close()
//	key := "user:" + hash
//	email, err := redis.String(c.Do("HGET", key, "email"))
//	if err != nil {
//		utl.Log("user check error", err)
//		return false, "", err
//	}
//	isok := email == item.Email
//	if isok {
//		che.Set(hash, email, 10*time.Minute)
//	}
//	return isok, hash, nil
//}

func KeyExists(hash string, pool *redis.Pool) (bool, error) {
	c := pool.Get()
	defer c.Close()
	key := "user:" + hash
	exist, err := redis.Bool(c.Do("EXISTS", key))
	if err != nil {
		utl.Log("check todo key error", err)
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
		utl.Log("check todo key error", err)
		return false, err
	}
	if exist {
		che.Set(hash, "foo", 10*time.Minute)
	}
	return exist, nil
}
