package main

import (
	"context"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"plumber/test_simple/proto"
)

func main() {
	// 创建一个tcp拨号
	conn, e := grpc.Dial(os.Args[1], grpc.WithInsecure()) // grpc.WithInsecure() 不安全的传输
	if e != nil {
		panic(e.Error())
	}
	client := proto.NewSimpleServerClient(conn)
	now := time.Now()
	simple, err := client.Simple(context.TODO(), &proto.Request{
		Msg: "hello world",
	})
	since := time.Since(now)
	log.Println("use ", since.Seconds())
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(simple)
}
