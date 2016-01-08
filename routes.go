package main

import (
	"github.com/gin-gonic/gin"
	"github.com/grunmax/GinRedisApi/ctrl"
	"github.com/grunmax/GinRedisApi/util"
)

func makeRoutes(api *gin.Engine) {
	util.Log("Server start", gin.Mode())

	ctrl.AddTodoRoutes(che, pool, api)
	ctrl.AddUserRoutes(che, pool, api)
	ctrl.AddTestRoutes(che, pool, api)
}

func is404(c *gin.Context) {
	c.JSON(404, util.BodyErr("wrong URL"))
}
