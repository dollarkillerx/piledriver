package main

import (
	"bufio"
	"bytes"
	"context"
	"google.golang.org/grpc/reflection"
	"io"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
	"plumber/rpc/proto"
)

func main() {
	log.SetFlags(log.Llongfile | log.LstdFlags)

	addr := "0.0.0.0:8091"
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	defer listen.Close()

	srv := grpc.NewServer()
	proto.RegisterPlumberServer(srv, &plumber{})
	reflection.Register(srv)

	log.Println("Success in", addr)
	if err := srv.Serve(listen); err != nil {
		log.Fatalln(err)
	}
}

type plumber struct {
	mu  sync.Mutex
	scp map[string]*net.TCPConn
}

func NewPlumber() {

}

func (p *plumber) send(data []byte, id string) {
	p.mu.Lock()

}

func (p *plumber) PlumberLinks(ctx context.Context, request *proto.PlumberRequest) (response *proto.PlumberResponse, err error) {
	log.Println("In")
	response = &proto.PlumberResponse{
		Data: []byte{},
	}
	addr, err := net.ResolveTCPAddr("tcp", request.Addr)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	if err := conn.SetLinger(0); err != nil {
		log.Println(err)
		return nil, err
	}

	if _, err := conn.Write(request.Data); err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(response.Data)
	writer := bufio.NewWriter(buffer)
	if _, err := io.Copy(writer, conn); err != nil {
		return nil, err
	}

	return response, nil
}

func (p *plumber) PlumberExchange(ctx context.Context, request *proto.PlumberRequest) (response *proto.PlumberResponse, err error) {
	log.Println("In")
	response = &proto.PlumberResponse{
		Data: []byte{},
	}
	addr, err := net.ResolveTCPAddr("tcp", request.Addr)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	if err := conn.SetLinger(0); err != nil {
		log.Println(err)
		return nil, err
	}

	if _, err := conn.Write(request.Data); err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(response.Data)
	writer := bufio.NewWriter(buffer)
	if _, err := io.Copy(writer, conn); err != nil {
		return nil, err
	}

	return response, nil
}

func (p *plumber) PlumberDisconnect(ctx context.Context, request *proto.PlumberRequest) (response *proto.PlumberResponse, err error) {
	log.Println("In")
	response = &proto.PlumberResponse{
		Data: []byte{},
	}
	addr, err := net.ResolveTCPAddr("tcp", request.Addr)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	if err := conn.SetLinger(0); err != nil {
		log.Println(err)
		return nil, err
	}

	if _, err := conn.Write(request.Data); err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(response.Data)
	writer := bufio.NewWriter(buffer)
	if _, err := io.Copy(writer, conn); err != nil {
		return nil, err
	}

	return response, nil
}
