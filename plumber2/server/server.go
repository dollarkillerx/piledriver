package main

import (
	"io"
	"log"
	"net"
	"plumber/plumber2/rpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log.SetFlags(log.Llongfile | log.LstdFlags)
	addr := "0.0.0.0:8086"
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	defer listen.Close()

	srv := grpc.NewServer()
	rpc.RegisterPlumberServer(srv, &server{})
	reflection.Register(srv)

	if err := srv.Serve(listen); err != nil {
		log.Fatalln(err)
	}
}

type server struct{}

func (s *server) Plumber(ser rpc.Plumber_PlumberServer) error {
	// 第一次建立链接进行拨号
	recv, err := ser.Recv()
	if err != nil {
		return err
	}
	addr, err := net.ResolveTCPAddr("tcp", recv.Addr)
	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	if err := conn.SetLinger(0); err != nil {
		log.Println(err)
		return err
	}

	go copy1(conn, ser)
	copy2(ser, conn)
	return nil
}

func copy2(server rpc.Plumber_PlumberServer, client io.Reader) {
	for {
		var b [1024]byte
		read, err := client.Read(b[:])
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println(err)
			break
		}

		if err := server.Send(&rpc.PlumberResponse{Data: b[:read]}); err != nil {
			log.Println(err)
			break
		}
	}
}

func copy1(server io.Writer, ser rpc.Plumber_PlumberServer) {
	for {
		recv, err := ser.Recv()
		if err != nil {
			log.Println(err)
			break
		}
		if recv.Over {
			break
		}
		if _, err := server.Write(recv.Data); err != nil {
			log.Println(err)
			break
		}
	}
}
