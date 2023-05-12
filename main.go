package main

import (
	"context"
	"fmt"
	"net"

	"github.com/cloudwego/kitex/server"
	"github.com/qml-123/AppService/cgo/ffmpeg"
	app "github.com/qml-123/AppService/kitex_gen/app/appservice"
	"github.com/qml-123/AppService/pkg/log"
	"github.com/qml-123/GateWay/common"
	"github.com/qml-123/app_log/logger"
)

const (
	configPath = "config/services.json"
)

func main() {
	ctx := context.Background()
	ffmpeg.Test()
	conf, err := common.GetJsonFromFile(configPath)
	if err != nil {
		panic(err)
	}

	if err = log.InitLogger(conf.EsUrl); err != nil {
		panic(err)
	}

	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:"+fmt.Sprintf("%d", conf.ListenPort))
	if err != nil {
		panic(err)
	}
	svr := app.NewServer(new(AppServiceImpl), server.WithServiceAddr(addr))

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
