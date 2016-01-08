package ctrl

import (
	"net/http"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/grunmax/GinRedisApi/accs"
	"github.com/grunmax/GinRedisApi/domn"
	"github.com/grunmax/GinRedisApi/util"
	"github.com/patrickmn/go-cache"
	"gopkg.in/validator.v2"
)

func AddUserRoutes(che *cache.Cache, pool *redis.Pool, routes *gin.Engine) {

	routes.PUT("/signup", func(c *gin.Context) {
		userForm := domn.UserSignupForm{}
		if err := c.Bind(&userForm); err != nil {
			c.JSON(400, util.BodyErr(err.Error()))
			return
		}
		if err := validator.Validate(userForm); err != nil {
			c.JSON(400, util.BodyErr(err.Error()))
			return
		}
		if userForm.Password != userForm.PasswordConfirm {
			c.JSON(400, util.BodyErr("check password confirm"))
			return
		}

		if hash, err := accs.UserCreate(userForm, pool); err != nil {
			c.JSON(400, util.BodyErr("Todo create error"))
		} else {
			http.SetCookie(c.Writer, &http.Cookie{Name: "k_yak", Value: hash, Path: "/", Domain: "", Expires: time.Now().AddDate(0, 0, 10)})
			//c.Writer.Header().Add("id", id)
			c.JSON(200, util.BodyOk("Autorized"))
		}
	})

	routes.PUT("/login", func(c *gin.Context) {
		userForm := domn.UserLoginForm{}
		if err := c.Bind(&userForm); err != nil {
			c.JSON(400, util.BodyErr(err.Error()))
			return
		}
		if err := validator.Validate(userForm); err != nil {
			c.JSON(400, util.BodyErr(err.Error()))
			return
		}

		if isok, hash, err := accs.UserCheck(userForm, pool); err != nil || !isok {
			c.JSON(400, util.BodyErr("check user error"))
		} else {
			http.SetCookie(c.Writer, &http.Cookie{Name: "k_yak", Value: hash, Path: "/", Domain: "", Expires: time.Now().AddDate(0, 0, 10)})
			c.JSON(200, util.BodyOk("Autorized"))
		}
	})

	routes.GET("/logout", func(c *gin.Context) {
		cookievalue, err := util.GetCookieValue("k_yak", c)
		if err == nil {
			che.Delete(cookievalue)
			http.SetCookie(c.Writer, &http.Cookie{Name: "k_yak", Value: "", Path: "/", Domain: "", Expires: time.Now().AddDate(0, 0, -1)})
			c.JSON(200, util.BodyOk("Unautorize ok"))
		} else {
			c.JSON(200, util.BodyOk("Unautorize error"))
		}
	})

}
