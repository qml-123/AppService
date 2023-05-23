package binlog

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/qml-123/AppService/pkg/db"
	"github.com/qml-123/AppService/pkg/id"
	"github.com/qml-123/AppService/pkg/redis"
	"github.com/qml-123/AppService/pkg/utils"
	"github.com/qml-123/app_log/logger"
)

type CompressEventHandler struct {
	canal.DummyEventHandler
}

func InitBinlog() error {
	cfg := canal.NewDefaultConfig()
	cfg.Addr = "bj-cynosdbmysql-grp-4u4v5eag.sql.tencentcdb.com:21216"
	cfg.User = "root"
	cfg.Password = "QmlGls08280709"
	cfg.Dump.ExecutionPath = ""
	cfg.Dump.TableDB = "app"
	cfg.Dump.Tables = []string{"files"}

	c, err := canal.NewCanal(cfg)
	if err != nil {
		return err
	}

	// Register a handler to handle RowsEvent
	c.SetEventHandler(&CompressEventHandler{})

	go func() {
		// Start canal
		c.Run()
	}()
	return nil
}

func interfaceToBool(i interface{}) (bool, error) {
	switch v := i.(type) {
	case bool:
		return v, nil
	case int, int8, int16, int32, int64:
		return v == 1, nil
	case uint, uint8, uint16, uint32, uint64:
		return v == 1, nil
	case string:
		return strconv.ParseBool(v)
	default:
		return false, fmt.Errorf("unsupported type: %v", reflect.TypeOf(i))
	}
}

func (h *CompressEventHandler) OnRow(e *canal.RowsEvent) error {
	if e.Action != canal.InsertAction || e.Table.Name != "files" {
		return nil
	}
	ctx := id.NewContext()

	testKeys := make([]string, 0)
	file_keys := make([]string, 0)
	//logger.Info(ctx, "rows: %v", e.Rows)
	for _, row := range e.Rows {
		var has_more, is_compressed, ok bool
		var err error
		var file_key string
		var is_video, is_txt_ bool
		for field, value := range row {
			if e.Table.Columns[field].Name == "file_type" {
				fileType, ok := value.(string)
				if ok && fileType == db.VideoColumnValue {
					is_video = true
				}
				if ok && fileType == "_txt_" {
					is_txt_ = true
				}
			}
			if e.Table.Columns[field].Name == "has_more" {
				has_more, err = interfaceToBool(value)
				if err != nil {
					logger.Warn(ctx, "has_more is not bool, value: %v, err: %v", value, err)
				}
			}
			if e.Table.Columns[field].Name == "is_compressed" {
				is_compressed, err = interfaceToBool(value)
				if err != nil {
					logger.Warn(ctx, "is_compressed is not bool, value: %v, err: %v", value, err)
				}
			}
			if e.Table.Columns[field].Name == "file_key" {
				file_key, ok = value.(string)
				if !ok {
					logger.Warn(ctx, "file_key is not string, value: %v, type: %v", value, reflect.ValueOf(value).Kind())
				}
			}
		}

		logger.Info(ctx, "file_key: %s, has_more: %v, is_compressed: %v", file_key, has_more, is_compressed)

		if is_video && !has_more && !is_compressed && file_key != "" {
			if file_key != "" {
				logger.Info(ctx, "file(%s) is upload end", file_key)
			}
			file_keys = append(file_keys, file_key)
		}

		if is_txt_ && !has_more && !is_compressed && file_key != "" {
			testKeys = append(testKeys, file_key)
		}
	}
	//go file.ContactFileChunks(ctx, testKeys)
	file_keys = utils.StrMapToSlice(utils.StrSliceToMap(file_keys))
	if len(file_keys) == 0 {
		return nil
	}
	logger.Info(ctx, "get upload end, file_keys: %v", file_keys)
	_, err := redis.SAdd("file_keys", file_keys...)
	if err != nil {
		logger.Warn(ctx, "SAdd failed, file_keys: %v, err: %v", file_keys, err)
	}
	return nil
}

func (h *CompressEventHandler) String() string {
	return "CompressEventHandler"
}
