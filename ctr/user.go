package ctr

import (
	"net/http"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/grunmax/GinRedisApi/acs"
	"github.com/grunmax/GinRedisApi/dom"
	"github.com/grunmax/GinRedisApi/utl"
	"github.com/patrickmn/go-cache"
	"gopkg.in/validator.v2"
)

const userCacheLife = 10 * time.Minute

func AddUserRoutes(che *cache.Cache, pool *redis.Pool, routes *gin.Engine) {

	routes.PUT("/signup", func(c *gin.Context) {
		userForm := dom.UserSignupForm{}
		if err := c.Bind(&userForm); err != nil {
			c.JSON(400, utl.BodyErr(err.Error()))
			return
		}
		if err := validator.Validate(userForm); err != nil {
			c.JSON(400, utl.BodyErr(err.Error()))
			return
		}
		if userForm.Password != userForm.PasswordConfirm {
			c.JSON(400, utl.BodyErr("check password confirm"))
			return
		}

		if hash, err := acs.UserCreate(userForm, pool); err != nil {
			c.JSON(400, utl.BodyErr("Todo create error"))
		} else {
			http.SetCookie(c.Writer, &http.Cookie{Name: "k_yak", Value: hash, Path: "/", Domain: "", Expires: time.Now().AddDate(0, 0, 10)})
			//c.Writer.Header().Add("id", id)
			c.JSON(200, utl.BodyOk("Autorized"))
		}
	})

	routes.PUT("/login", func(c *gin.Context) {
		userForm := dom.UserLoginForm{}
		if err := c.Bind(&userForm); err != nil {
			c.JSON(400, utl.BodyErr(err.Error()))
			return
		}
		if err := validator.Validate(userForm); err != nil {
			c.JSON(400, utl.BodyErr(err.Error()))
			return
		}

		if isok, hash, err := acs.UserCheck(userForm, pool); err != nil || !isok {
			c.JSON(400, utl.BodyErr("check user error"))
		} else {
			http.SetCookie(c.Writer, &http.Cookie{Name: "k_yak", Value: hash, Path: "/", Domain: "", Expires: time.Now().AddDate(0, 0, 10)})
			c.JSON(200, utl.BodyOk("Autorized"))
		}
	})

	routes.GET("/logout", func(c *gin.Context) {
		cookievalue, err := utl.GetCookieValue("k_yak", c)
		if err == nil {
			che.Delete(cookievalue)
			http.SetCookie(c.Writer, &http.Cookie{Name: "k_yak", Value: "", Path: "/", Domain: "", Expires: time.Now().AddDate(0, 0, -1)})
			c.JSON(200, utl.BodyOk("Unautorize ok"))
		} else {
			c.JSON(200, utl.BodyOk("Unautorize error"))
		}
	})

}
