package main

import (
	"context"
	"io"
	"log"

	"google.golang.org/grpc"
	"plumber/test_stream/proto"
)

func main() {
	addr := "0.0.0.0:8086"
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	client := proto.NewStreamServerClient(conn)
	req := &proto.Request{
		Msg: "ppp",
	}
	stream, err := client.Stream(context.TODO())
	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		for {
			if err := stream.Send(req); err != nil {
				if err == io.EOF {
					log.Println("O2")
					break
				}
				log.Fatalln(err)
			}
		}
	}()

	for {
		recv, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				log.Println("O1")
				break
			}
			log.Fatalln(err)
		}
		log.Println(recv)
	}

}
