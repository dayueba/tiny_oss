package rpc

type NameModel struct {
	bucketName string
	// 用户key
	key      string
	objectId int64
}

type ObjectModel struct {
	objectId int64
	blocks   []byte //blockid列表 用PB压缩(DataBocks) binary类型
}

type DataModel struct {
	blockId int64
	data    []byte
	shardId int8
}

type Shard struct {
	key  int8
	name string // mysql name
}
