package rpc

import (
	"context"
	"gorm.io/gorm"
)
import "github.com/redis/go-redis/v9"

// 虚拟shard数量
var virtualShard = 2233

// 处理具体文件数据的server

type DataServer struct {
	routerDb redis.Client       // 只需要存储 virtualShard 条数据
	dataDbs  map[string]gorm.DB // 存储 data 的db
}

func (srv *DataServer) Shard(key string) gorm.DB {
	hashKey := hash(key) % virtualShard
	dataKey, _ := srv.routerDb.Get(context.Background(), string(hashKey)).Result()
	return srv.dataDbs[dataKey]
}

func hash(key string) int {
	return 0
}
