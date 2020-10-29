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
	log.SetFlags(log.Llongfile | log.LstdFlags)

	addr := "0.0.0.0:443"
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	defer listen.Close()

	srv := grpc.NewServer()
	proto.RegisterSimpleServerServer(srv, &server{})
	reflection.Register(srv)
}

type server struct{}

func (s *server) Simple(ctx context.Context, req *proto.Request) (resp *proto.Response, err error) {
	log.Println(req.Msg)

	return &proto.Response{
		Msg: req.Msg + "Spc",
	}, nil
}
