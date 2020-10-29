package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"plumber/test_stream/proto"
)

func main() {
	addr := "0.0.0.0:8086"
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	defer listen.Close()

	srv := grpc.NewServer()
	proto.RegisterStreamServerServer(srv, &server{})
	reflection.Register(srv)

	if err := srv.Serve(listen); err != nil {
		log.Fatalln(err)
	}
}

type server struct{}

func (s *server) Stream(ser proto.StreamServer_StreamServer) error {
	recv, err := ser.Recv()
	if err != nil {
		return err
	}
	log.Println(recv.Msg)
	ser.Send(&proto.Response{Msg: recv.Msg})
	return nil
}
