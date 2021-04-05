package main

import (
	"log"

	"smicro/server"

	"smicro/output/router"

	"smicro/output/generate/com/google/hello"
)

var routerServer = &router.RouterServer{}

func main() {

	err := server.Init("com.google.hello")
	if err != nil {
		log.Fatal("init service failed, err:%v", err)
		return
	}

	hello.RegisterHelloServiceServer(server.GRPCServer(), routerServer)
	server.Run()
}
