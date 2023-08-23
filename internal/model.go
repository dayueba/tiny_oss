package internal

type Name struct {
	bucketName string
	// 用户key
	key      string
	objectId int64
}

type Object struct {
	objectId int64
	blocks   []byte //blockid列表 用PB压缩(DataBocks) binary类型
}

type Data struct {
	blockId int64
	data    []byte
}
