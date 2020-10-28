package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"

	"plumber/test_simple/proto"
)

func main() {
	// 创建一个tcp服务
	addr := "0.0.0.0:9001"
	listener, e := net.Listen("tcp", addr)
	if e != nil {
		panic(e.Error())
	}

	// 创建一个没有注册的grpc
	srv := grpc.NewServer()
	proto.RegisterSimpleServerServer(srv, &server{}) // 注册上去
	reflection.Register(srv)                         // 注册在给定的gRPC服务器上注册服务器反射服务

	// 监听
	log.Println("Listen :", addr)
	if e := srv.Serve(listener); e != nil {
		panic(e.Error())
	}
}

type server struct{}

func (s *server) Simple(ctx context.Context, req *proto.Request) (resp *proto.Response, err error) {
	log.Println(req.Msg)

	return &proto.Response{
		Msg: req.Msg + "Spc",
	}, nil
}
