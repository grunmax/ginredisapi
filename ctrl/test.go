package ctrl

import (
	"bytes"
	"image/jpeg"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/grunmax/GinRedisApi/accs"
	"github.com/grunmax/GinRedisApi/util"
	"github.com/patrickmn/go-cache"
)

func AddTestRoutes(che *cache.Cache, pool *redis.Pool, routes *gin.Engine) {

	routes.GET("/", func(c *gin.Context) {
		if item, err := accs.TestFunc(1, "gin: redis works", pool); err != nil {
			c.JSON(400, util.BodyErr(err.Error()))
		} else {
			c.JSON(200, item)
		}
	})

	routes.POST("/file", func(c *gin.Context) {
		file, _, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(400, util.BodyErr(err.Error()))
			return
		}
		const maxSize = 200000
		buf := make([]byte, maxSize)
		//buf := new(bytes.Buffer)

		n, err := file.Read(buf)
		if err != nil {
			c.JSON(400, util.BodyErr(err.Error()))
			return
		}
		if n == maxSize {
			c.JSON(400, util.BodyErr("big file"))
			return
		}
		if item, err := accs.RedisSet("file:", buf, pool); err != nil {
			c.JSON(400, util.BodyErr("file store error"))
		} else {
			c.JSON(200, item)
		}
	})

	routes.GET("/file/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")
		if file, err := accs.RedisGet("file:", id, pool); err != nil {
			c.JSON(400, util.BodyErr("file get error"))
		} else {
			//c.Writer.Header().Set("Content-Type", "application/octet-stream")
			c.Data(200, "application/octet-stream", file)
		}

	})

	routes.POST("/img", func(c *gin.Context) {
		file, _, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(400, util.BodyErr(err.Error()))
			return
		}
		img, err := jpeg.Decode(file)
		if err != nil {
			c.JSON(400, util.BodyErr(err.Error()))
			return
		}
		buf := new(bytes.Buffer)
		if err := jpeg.Encode(buf, img, nil); err != nil {
			c.JSON(400, util.BodyErr(err.Error()))
			return
		}
		if buf.Len() > 20000 {
			c.JSON(400, util.BodyErr("big file"))
			return
		}
		if item, err := accs.RedisSet("img:", buf.Bytes(), pool); err != nil {
			c.JSON(400, util.BodyErr(err.Error()))
		} else {
			c.JSON(200, item)
		}
	})

	routes.GET("/img/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")
		if file, err := accs.RedisGet("img:", id, pool); err != nil {
			c.JSON(400, util.BodyErr(err.Error()))
		} else {
			che.Set(c.Request.RequestURI, file, todoCacheLife)
			c.Data(200, "image/jpeg", file)
		}

	})

}
