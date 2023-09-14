package rpc

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
	"mime/multipart"
	"net/http"
	"tiny_oss/pb"
)

var db *gorm.DB

type App struct {
	app *gin.Engine

	// 解耦
	// data server client
	// metadata server client
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
		// 以bucketName和s3文件名为key，从name表中查找到对应的objectid
		nameModel := new(NameModel)
		db.Where("bucket_name = ?", bucketName).Where("object_name = ?", objectName).First(nameModel)

		// 以objectid为key，从object表中查找到对应的block列表
		objectModel := new(ObjectModel)
		db.Where("object_id = ?", nameModel.objectId).First(objectModel)

		// 根据block列表中的长度信息，根据用户读请求的区间，计算出需要读取哪些blockid
		blocks := new(pb.DataBlocks)
		proto.Unmarshal(objectModel.blocks, blocks)

		// 以blockid为key，从data表中读取对应的数据，并进行截断和拼装，返回结果。(假设上传一个10MB的文件，1MB一个block，指定读取区间 [1MB+1B, 2M+100KB], 则需要读取第1和第2个block，并进行截断)
		dataModels := make([]*DataModel, 0)
		for _, block := range blocks.Blocks {
			model := new(DataModel)
			db.Model(block.TableName).Where("block_id = ?", block.BlockId).First(model)
			dataModels = append(dataModels, model)
		}

		c.String(http.StatusOK, "Hello %s %s", bucketName, objectName)
	})

	// 以 bucketName + objectName 作为key(主键)，上传数据(比如一段视频)
	r.PUT("/:bucket_name/:object_name", func(c *gin.Context) {
		bucketName := c.Param("bucket_name")
		objectName := c.Param("object_name")
		file, _ := c.FormFile("file")

		// 写入流程
		// 将数据切割为若干个block，每个block分配一个id(先不管哪里来的)。以blockid为主键，将数据写入到data表中
		dataBlocks := new(pb.DataBlocks)
		split := func(file *multipart.FileHeader) {
			id := genId()
			dataBlocks.Blocks = append(dataBlocks.Blocks, &pb.Block{
				TableName: "",
				Size:      0,
				BlockId:   id,
			})
			block := DataModel{
				blockId: id,
				data:    nil,
				shardId: 0,
			}
			db.Save(block) // 可以改为一次存储多个block
		}
		split(file)

		// 分配一个object_id(先忽略从哪里分配)，将步骤1中的blockid和数据表的表名压缩到一段pb中，然后以object_id为主键，将这段pb写入到object表中
		objectId := genId()
		blocks, _ := proto.Marshal(dataBlocks)
		object := ObjectModel{
			objectId: objectId,
			blocks:   blocks,
		}
		db.Save(object)

		// 以bucketName和s3文件名组成主键，将object_id写入到name表中
		nameModel := &NameModel{
			bucketName: bucketName,
			key:        objectName,
			objectId:   objectId,
		}
		db.Save(nameModel) // 多次保存数据，可以放在一个事务里

		// 路由层收到key之后， hash(key) % 2233 得到一个[0, 2232]之间的数字，比如X， 即shardX。
		// 查询路由MySQL，得到shardX所对应后端的数据存储层MySQL服务器的地址。
		// 与数据存储层MySQL进行通信，获取数据。

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

// 唯一标识一段数据，不能有重复
// 使用mysql自增 or 其他外部的ID分配器
// 如果使用mysql自增，替换了以后就不能用了
func genId() int64 {
	// 一般不使用value的crc
	// 使用CRC的问题在于，同一个文件两次上传(使用不同的文件名)，以block的CRC和整个的CRC分别作为blockid和objectid，会导致在data表和object表中，只有一条记录(CRC相同)，删除的时候会导致实际的数据被删除
	// 当然我们可以使用CRC做去重，就是增加一些工作量
	return 0
}
