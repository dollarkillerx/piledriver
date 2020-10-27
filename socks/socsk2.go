package main

import (
	"bufio"
	"bytes"
	"context"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"plumber/rpc/proto"
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
	
	for {
		client, err := l.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleClientRequest2(client)
	}
}
func handleClientRequest2(client net.Conn) {
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

		log.Println(addr)
		client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}) //响应客户端连接成功
		//进行转发

		dial, err := grpc.Dial("0.0.0.0：8091", grpc.WithInsecure())
		if err != nil {
			log.Fatalln(err)
		}

		plumber := proto.NewPlumberClient(dial)
		data, err := ral(client)
		if err != nil {
			log.Println(err)
			return
		}


		response, err := plumber.Plumber(context.TODO(), &proto.PlumberRequest{Data: data, Addr: addr})
		if err != nil {
			log.Println(err)
			return
		}

		client.Write(response.Data)
	}
}

func ral(client io.Reader) (data []byte,err error) {
	buffer := bytes.NewBuffer(data)
	writer := bufio.NewWriter(buffer)
	for {
		var b [1024]byte
		read, err := client.Read(b[:])
		if err != nil {
			if err == io.EOF {
				log.Println("bt")
				break
			}
			log.Println(err)
			break
		}
		log.Println("In")
		writer.Write(b[:read])
	}
	return buffer.Bytes(), nil
}

func copy2(server io.Writer, client io.Reader) {
	for {
		var b [1024]byte
		read, err := client.Read(b[:])
		if err != nil {
			if err == io.EOF {
				break
			}
			break
		}


		if _, err := server.Write(b[:read]); err != nil {
			log.Println(err)
			break
		}
	}
}
