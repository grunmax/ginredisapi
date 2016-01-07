package main

import (
	"github.com/gin-gonic/gin"
	"github.com/grunmax/GinRedisApi/acs"
	"github.com/grunmax/GinRedisApi/ctr"
	"github.com/grunmax/GinRedisApi/utl"
)

func makeRoutes(api *gin.Engine) {

	api.GET("/", func(c *gin.Context) {
		if item, err := acs.TestFunc(1, "gin: redis works", pool); err != nil {
			c.JSON(400, utl.BodyErr(err.Error()))
		} else {
			c.JSON(200, item)
		}
	})

	ctr.AddTodoRoutes(che, pool, api)
}

func is404(c *gin.Context) {
	c.JSON(404, utl.BodyErr("wrong page"))
}
