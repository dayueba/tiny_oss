package internal

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type App struct {
	app *gin.Engine
}

func NewApp() *App {
	r := gin.Default()

	// 以 bucketName + objectName 为key，下载数据(比如刚刚上传的视频)
	r.GET("/:bucket_name/:object_name", func(c *gin.Context) {
		bucketName := c.Param("bucket_name")
		objectName := c.Param("object_name")
		c.String(http.StatusOK, "Hello %s %s", bucketName, objectName)
	})

	// 以bucketName和指定的前缀为过滤条件，列出符合条件的key。比如以bucketName为前缀，列出已经上传的文件名 列表
	// (根据上面的例子，就是列出已经上传的视频文件名)
	r.GET("/:bucket_name", func(c *gin.Context) {
		bucketName := c.Param("bucket_name")
		objectName := c.Param("object_name")
		c.String(http.StatusOK, "Hello %s %s", bucketName, objectName)
	})

	// 以 bucketName + objectName 作为key(主键)，上传数据(比如一段视频)
	r.PUT("/:bucket_name/:object_name", func(c *gin.Context) {
		bucketName := c.Param("bucket_name")
		objectName := c.Param("object_name")
		c.String(http.StatusOK, "Hello %s %s", bucketName, objectName)
	})

	// 以 bucketName + objectName 为key，把之前上传的视频进行删除
	r.DELETE("/:bucket_name/:object_name", func(c *gin.Context) {
		bucketName := c.Param("bucket_name")
		objectName := c.Param("object_name")
		c.String(http.StatusOK, "Hello %s %s", bucketName, objectName)
	})

	return &App{
		app: r,
	}
}

func (app *App) Run() error {
	return app.app.Run()
}
