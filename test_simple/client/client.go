package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"plumber/test_simple/proto"
)

func main() {
	// 创建一个tcp拨号
	conn, e := grpc.Dial(":9001", grpc.WithInsecure()) // grpc.WithInsecure() 不安全的传输
	if e != nil {
		panic(e.Error())
	}
	client := proto.NewSimpleServerClient(conn)
	simple, err := client.Simple(context.TODO(), &proto.Request{
		Msg: "hello world",
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(simple)
}
