package ctrl

import (
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/grunmax/GinRedisApi/accs"
	"github.com/grunmax/GinRedisApi/domn"
	"github.com/grunmax/GinRedisApi/util"
	"github.com/patrickmn/go-cache"
	"gopkg.in/validator.v2"
)

const todoCacheLife = 10 * time.Minute

func AddTodoRoutes(che *cache.Cache, pool *redis.Pool, routes *gin.Engine) {

	routes.GET("/todo", func(c *gin.Context) {
		if keys, err := accs.GetKeys("todo:*", pool); err != nil {
			c.JSON(400, util.BodyErr("Todo get keys error"))
		} else {
			c.JSON(200, keys)
		}
	})

	routes.GET("/todo/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")
		if item, err := accs.TodoGetId(id, pool); err != nil {
			c.JSON(400, util.BodyErr(err.Error()))
		} else {
			che.Set(c.Request.RequestURI, item, todoCacheLife)
			c.JSON(200, item)
		}
	})

	routes.POST("/todo", func(c *gin.Context) {
		todoForm := domn.TodoForm{}
		if err := c.Bind(&todoForm); err != nil {
			c.JSON(400, util.BodyErr(err.Error()))
			return
		}
		if err := validator.Validate(todoForm); err != nil {
			c.JSON(400, util.BodyErr(err.Error()))
			return
		}
		if item, err := accs.TodoCreate(todoForm, pool); err != nil {
			c.JSON(400, util.BodyErr("Todo create error"))
		} else {
			c.Writer.Header().Add("id", item.Id)
			c.JSON(200, item)
		}
	})

	routes.POST("/todo/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")
		if id == "" {
			c.JSON(400, util.BodyErr("Empty id"))
			return
		}
		todoForm := domn.TodoForm{}
		if err := c.Bind(&todoForm); err != nil {
			c.JSON(400, util.BodyErr(err.Error()))
			return
		}
		if err := validator.Validate(todoForm); err != nil {
			c.JSON(400, util.BodyErr(err.Error()))
			return
		}
		if item, err := accs.TodoEdit(id, todoForm, pool); err != nil {
			c.JSON(400, util.BodyErr(err.Error()))
		} else {
			c.Writer.Header().Add("id", item.Id)
			che.Delete(c.Request.RequestURI)
			c.JSON(200, item)
		}
	})

	//routes.DELETE("/todos", func(c *gin.Context) {
	//	c.JSON(200, todo.DeleteAll())
	//})

	routes.DELETE("/todo/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")
		if err := accs.TodoDeleteId(id, pool); err != nil {
			c.JSON(400, util.BodyErr("Todo delete id error"))
		} else {
			che.Delete(c.Request.RequestURI)
			c.JSON(200, util.BodyOk("deleted"))
		}
	})
}
