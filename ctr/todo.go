package ctr

import (
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/grunmax/GinRedisApi/acs"
	"github.com/grunmax/GinRedisApi/dom"
	"github.com/grunmax/GinRedisApi/utl"
	"gopkg.in/validator.v2"
)

func AddTodoRoutes(pool *redis.Pool, routes *gin.Engine) {

	routes.GET("/todo", func(c *gin.Context) {
		if keys, err := acs.TodoGetKeys("todo:*", pool); err != nil {
			c.JSON(400, utl.BodyErr("Todo get keys error"))
		} else {
			c.JSON(200, keys)
		}
	})

	routes.GET("/todo/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")
		if item, err := acs.TodoGetId(id, pool); err != nil {
			c.JSON(400, utl.BodyErr("Todo getid error"))
		} else {
			c.Writer.Header().Add("id", item.Id)
			c.JSON(200, item)
		}
	})

	routes.POST("/todo", func(c *gin.Context) {
		todoForm := dom.TodoForm{}
		if err := c.Bind(&todoForm); err != nil {
			c.JSON(400, utl.BodyErr(err.Error()))
			return
		}
		if err := validator.Validate(todoForm); err != nil {
			c.JSON(400, utl.BodyErr(err.Error()))
			return
		}
		if item, err := acs.TodoCreate(todoForm, pool); err != nil {
			c.JSON(400, utl.BodyErr("Todo create error"))
		} else {
			c.Writer.Header().Add("id", item.Id)
			c.JSON(200, item)
		}
	})

	routes.POST("/todo/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")
		if id == "" {
			c.JSON(400, utl.BodyErr("Empty id"))
			return
		}
		todoForm := dom.TodoForm{}
		if err := c.Bind(&todoForm); err != nil {
			c.JSON(400, utl.BodyErr(err.Error()))
			return
		}
		if err := validator.Validate(todoForm); err != nil {
			c.JSON(400, utl.BodyErr(err.Error()))
			return
		}
		if item, err := acs.TodoEdit(id, todoForm, pool); err != nil {
			c.JSON(400, utl.BodyErr(err.Error()))
		} else {
			c.Writer.Header().Add("id", item.Id)
			c.JSON(200, item)
		}
	})

	//routes.DELETE("/todos", func(c *gin.Context) {
	//	c.JSON(200, todo.DeleteAll())
	//})

	routes.DELETE("/todo/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")
		if err := acs.TodoDeleteId(id, pool); err != nil {
			c.JSON(400, utl.BodyErr("Todo delete id error"))
		} else {
			c.Writer.Header().Add("id", id)
			c.JSON(200, utl.BodyOk("deleted"))
		}
	})
}
