package id

import (
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/qml-123/AppService/pkg/redis"
)

const (
	generateIDLock = "generate_id_lock"
	generateIDKey  = "generate_id_key"
)

var (
	_node *snowflake.Node
)

func InitGen() error {
	var id int64
	for {
		success, err := redis.SetNX(generateIDLock, 1, 30*time.Second)
		if err == nil && success {
			id, err = redis.Incr(generateIDKey)
			if err == nil {
				redis.Del(generateIDLock)
				break
			}
			redis.Del(generateIDLock)
		}

		time.Sleep(2 * time.Second)
	}

	var err error
	_node, err = snowflake.NewNode(id)
	if err != nil {
		return err
	}
	return nil
}

func Generate() snowflake.ID {
	return _node.Generate()
}
