package main

import (
	"fmt"
	"log"
	"net"

	"github.com/cloudwego/kitex/server"
	"github.com/qml-123/AppService/cgo/rav1e"
	"github.com/qml-123/AppService/cgo/test/file"
	app "github.com/qml-123/AppService/kitex_gen/app/appservice"
	"github.com/qml-123/GateWay/common"
)

const (
	configPath = "config/services.json"
)

func main() {
	_test()
	conf, err := common.GetJsonFromFile(configPath)
	if err != nil {
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
		log.Println(err.Error())
	}
}

func _test() {
	input := file.GetFileData("/root/go/src/github.com/qml-123/AppService/output/bin/1.txt")
	got, err := rav1e.Encode(input)
	if err != nil {
		log.Panicf("err: %v", err)
	}
	file.SaveFileData(got, "/root/go/src/github.com/qml-123/AppService/output/bin/en_1.txt")

	input = file.GetFileData("/root/go/src/github.com/qml-123/AppService/output/bin/en_1.txt")
	got, err = rav1e.Decode(input)
	if err != nil {
		log.Panicf("err: %v", err)
	}
	file.SaveFileData(got, "/root/go/src/github.com/qml-123/AppService/output/bin/de_1.txt")
}
