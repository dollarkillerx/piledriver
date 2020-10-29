package main

import (
	"context"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
	"plumber/plumber2/rpc"
	"strconv"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	addr, err := net.ResolveTCPAddr("tcp", ":8081")
	if err != nil {
		log.Fatalln(err)
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Panic(err)
	}

	addrS := "0.0.0.0:8086"
	s := New(addrS)
	for {
		client, err := l.Accept()
		if err != nil {
			log.Panic(err)
		}
		go s.handleClientRequest(client)
	}
}

type server struct {
	client rpc.PlumberClient
}

func New(addr string) *server {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	client := rpc.NewPlumberClient(conn)
	return &server{client: client}
}

func (s *server) handleClientRequest(client net.Conn) {
	if client == nil {
		return
	}
	defer client.Close()
	var b [1024]byte
	n, err := client.Read(b[:])
	if err != nil {
		log.Println(err)
		return
	}
	if b[0] == 0x05 { //只处理Socket5协议
		//客户端回应：Socket服务端不需要验证方式
		client.Write([]byte{0x05, 0x00})
		n, err = client.Read(b[:])
		var host, port string
		switch b[3] {
		case 0x01: //IP V4
			host = net.IPv4(b[4], b[5], b[6], b[7]).String()
		case 0x03: //域名
			host = string(b[5 : n-2]) //b[4]表示域名的长度
		case 0x04: //IP V6
			host = net.IP{b[4], b[5], b[6], b[7], b[8], b[9], b[10], b[11], b[12], b[13], b[14], b[15], b[16], b[17], b[18], b[19]}.String()
		}
		port = strconv.Itoa(int(b[n-2])<<8 | int(b[n-1]))

		addr := net.JoinHostPort(host, port)

		client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}) //响应客户端连接成功
		//进行转发
		stream, err := s.client.Plumber(context.TODO())
		if err != nil {
			log.Println(err)
			return
		}

		if err := stream.Send(&rpc.PlumberRequest{
			Addr: addr,
		}); err != nil {
			log.Println(err)
			return
		}
		go copy1(stream, client)
		copy2(client, stream)
	}
}

func copy1(server rpc.Plumber_PlumberClient, client io.Reader) {
	for {
		var b [1024]byte
		read, err := client.Read(b[:])
		if err != nil {
			if err == io.EOF {
				server.Send(&rpc.PlumberRequest{Over: true})
				break
			}
			break
		}

		if err := server.Send(&rpc.PlumberRequest{Data: b[:read]}); err != nil {
			log.Println(err)
			break
		}
	}
}

func copy2(client io.Writer, server rpc.Plumber_PlumberClient) {
	for {
		recv, err := server.Recv()
		if err != nil {
			log.Println(err)
			break
		}
		if _, err := client.Write(recv.Data); err != nil {
			log.Println(err)
		}
	}
}
