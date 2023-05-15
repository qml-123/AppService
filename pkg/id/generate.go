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
		if err != nil {
			return err
		}
		if success {
			id, err = redis.Incr(generateIDKey)
			redis.Del(generateIDLock)
			if err == nil {
				break
			}
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
