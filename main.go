package main

import (
	"fmt"
	"net"
	"time"

	"github.com/cloudwego/kitex/server"
	"github.com/qml-123/AppService/middleware"
	"github.com/qml-123/AppService/pkg/consumer/binlog"
	"github.com/qml-123/AppService/pkg/consumer/delay_task"
	"github.com/qml-123/AppService/pkg/db"
	"github.com/qml-123/AppService/pkg/id"
	"github.com/qml-123/AppService/pkg/log"
	"github.com/qml-123/AppService/pkg/redis"
	"github.com/qml-123/app_log/common"
	"github.com/qml-123/app_log/kitex_gen/app/appservice"
	"github.com/qml-123/app_log/logger"
)

const (
	configPath = "config/services.json"
)

func main() {
	ctx := id.NewContext()
	//ffmpeg.Test()
	conf, err := common.GetJsonFromFile(configPath)
	if err != nil {
		panic(err)
	}

	if err = log.InitLogger(conf.EsUrl); err != nil {
		panic(err)
	}

	if err = db.InitDB(); err != nil {
		panic(err)
	}

	if err = binlog.InitBinlog(); err != nil {
		panic(err)
	}

	if err = redis.InitRedis(); err != nil {
		panic(err)
	}

	if err = id.InitGen(); err != nil {
		panic(err)
	}

	if err = delay_task.Init(); err != nil {
		panic(err)
	}

	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:"+fmt.Sprintf("%d", conf.ListenPort))
	if err != nil {
		panic(err)
	}
	svr := appservice.NewServer(new(AppServiceImpl), server.WithServiceAddr(addr), server.WithMiddleware(middleware.ErrResponseMW), server.WithReadWriteTimeout(5 * time.Minute))

	addr, _ = net.ResolveTCPAddr("tcp", conf.ListenIp+":"+fmt.Sprintf("%d", conf.ListenPort))
	if err = common.InitConsul(addr, conf); err != nil {
		panic(err)
	}

	defer common.CloseConsul(addr, conf)

	err = svr.Run()
	if err != nil {
		logger.Warn(ctx, "Failed to run server: %v", err)
	}
}