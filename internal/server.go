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

		// 读取流程
		// 读取流程和写入流程相反
		// 以bucketName和s3文件名为key，从name表中查找到对应的objectid
		// 以objectid为key，从object表中查找到对应的block列表
		// 根据block列表中的长度信息，根据用户读请求的区间，计算出需要读取哪些blockid
		// 以blockid为key，从data表中读取对应的数据，并进行截断和拼装，返回结果。(假设上传一个10MB的文件，1MB一个block，指定读取区间 [1MB+1B, 2M+100KB], 则需要读取第1和第2个block，并进行截断)

		c.String(http.StatusOK, "Hello %s %s", bucketName, objectName)
	})

	// 以 bucketName + objectName 作为key(主键)，上传数据(比如一段视频)
	r.PUT("/:bucket_name/:object_name", func(c *gin.Context) {
		bucketName := c.Param("bucket_name")
		objectName := c.Param("object_name")

		// 写入流程
		// 接入层接受到数据之后，将数据切割为若干个block，每个block分配一个id(先不管哪里来的)。以blockid为主键，将数据写入到data表中
		// 分配一个object_id(先忽略从哪里分配)，将步骤1中的blockid和数据表的表名压缩到一段pb中，然后以object_id为主键，将这段pb写入到object表中
		// 以bucketName和s3文件名组成主键，将object_id写入到name表中

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
