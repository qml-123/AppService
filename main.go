package main

import (
	"fmt"
	"log"
	"net"

	"github.com/cloudwego/kitex/server"
	"github.com/qml-123/AppService/cgo/av1"
	"github.com/qml-123/AppService/cgo/test/file"
	app "github.com/qml-123/AppService/kitex_gen/app/appservice"
	"github.com/qml-123/GateWay/common"
)

const (
	configPath = "config/services.json"
)

func main() {
	test()
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

func test() {
	input := file.GetFileData("1.txt")
	got, err := av1.AV1Encode(input)
	if err != nil {
		log.Panicf("err: %v", err)
	}
	file.SaveFileData(got, "en_1.txt")

	input = file.GetFileData("en_1.txt")
	got, err = av1.AV1Decode(input)
	if err != nil {
		log.Panicf("err: %v", err)
	}
	file.SaveFileData(got, "de_1.txt")
}