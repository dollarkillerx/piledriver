package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"plumber/rpc"
	"plumber/storage"
	"plumber/utils"
	"strconv"
	"time"

	"github.com/dollarkillerx/easy_dns"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const pem = `
-----BEGIN CERTIFICATE-----
MIICYzCCAeqgAwIBAgIUTShO8REvwRrDoBdRq9ZzzcEMP6swCgYIKoZIzj0EAwIw
aTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoMGElu
dGVybmV0IFdpZGdpdHMgUHR5IEx0ZDEQMA4GA1UECwwHcGx1bWJlcjEQMA4GA1UE
AwwHcGx1bWJlcjAeFw0yMDEwMjcxMTU1MDZaFw0zMDEwMjUxMTU1MDZaMGkxCzAJ
BgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5l
dCBXaWRnaXRzIFB0eSBMdGQxEDAOBgNVBAsMB3BsdW1iZXIxEDAOBgNVBAMMB3Bs
dW1iZXIwdjAQBgcqhkjOPQIBBgUrgQQAIgNiAATvf5qn0hKp56iFULTyXfTNRVMf
8mBgQR8GiJlhMNs8SE/128T2lD0UFDMrAVxKxo/rHsP5ORiP/uc9NnK721dVlmcH
40XGXtW2BbThJeCFdNO1ife81fuqzxWx4oSIGWajUzBRMB0GA1UdDgQWBBQoI4sE
Nx8JzONu9VixAeN1Kr5GdzAfBgNVHSMEGDAWgBQoI4sENx8JzONu9VixAeN1Kr5G
dzAPBgNVHRMBAf8EBTADAQH/MAoGCCqGSM49BAMCA2cAMGQCMEa8IP5y+EzZOzrm
O9yFZkNqkBkFl00M0GAR5wvO1W6pyC7tJfvgxd8C3mClltakOgIwXzDvpKz9eG4h
59r69pUhx2Jc5ffDO/SbNEx51o3zOvdR77OxaJvNS/cyC2TENO8C
-----END CERTIFICATE-----
`

// ./client addr socketAddr user passwd
func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	if len(os.Args) < 5 {
		log.Fatalln("what fuck???")
	}
	addr, err := net.ResolveTCPAddr("tcp", os.Args[2])
	if err != nil {
		log.Fatalln(err)
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Socket Addr: ", os.Args[2])
	s := New(os.Args[1])
	if len(os.Args) > 5 {
		s.pac = true
		log.Println("Pac Model")
	}

	for {
		client, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go s.handleClientRequest(client)
	}
}

type server struct {
	client rpc.PlumberClient
	dns    string
	pac    bool
}

func New(addr string) *server {
	creds, err := NewClientTLSFromFile(pem, "plumber")
	if err != nil {
		log.Fatalln(err)
	}

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds), grpc.WithPerRPCCredentials(&loginCreds{
		Username: os.Args[3],
		Password: os.Args[4],
	}))

	client := rpc.NewPlumberClient(conn)
	if len(os.Args) > 5 {
		return &server{client: client, dns: os.Args[5]}
	} else {
		return &server{client: client}
	}
}

func (s *server) handleClientRequest(client net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("%s\n", err)
			return
		}
	}()

	if client == nil {
		return
	}
	defer client.Close()
	var b [1024]byte
	n, err := client.Read(b[:])
	if err != nil {
		//log.Println(err)
		return
	}
	var pac bool
	if b[0] == 0x05 { //只处理Socket5协议
		//客户端回应：Socket服务端不需要验证方式
		client.Write([]byte{0x05, 0x00})
		n, err = client.Read(b[:])
		var host, port string
		switch b[3] {
		case 0x01: //IP V4
			host = net.IPv4(b[4], b[5], b[6], b[7]).String()
			pac, err = s.isPac(host, "")
			if err != nil {
				//log.Println(err)
				return
			}
		case 0x03: //域名
			host = string(b[5 : n-2]) //b[4]表示域名的长度
			pac, err = s.isPac("", host)
			if err != nil {
				//log.Println(err)
				return
			}
		case 0x04: //IP V6
			//host = net.IP{b[4], b[5], b[6], b[7], b[8], b[9], b[10], b[11], b[12], b[13], b[14], b[15], b[16], b[17], b[18], b[19]}.String()
			return
		}
		port = strconv.Itoa(int(b[n-2])<<8 | int(b[n-1]))

		addr := net.JoinHostPort(host, port)

		if _, err := client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}); err != nil { //响应客户端连接成功
			//log.Println(err)
			return
		}

		if pac {
			simple(client, addr)
			return
		}

		//进行转发
		stream, err := s.client.Plumber(context.TODO())
		if err != nil {
			//log.Println(err)
			return
		}

		if err := stream.Send(&rpc.PlumberRequest{
			Addr: addr,
		}); err != nil {
			//log.Println(err)
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
			//log.Println(err)
			break
		}
	}
}

func copy2(client io.Writer, server rpc.Plumber_PlumberClient) {
	for {
		recv, err := server.Recv()
		if err != nil {
			//log.Println(err)
			break
		}
		if _, err := client.Write(recv.Data); err != nil {
			//log.Println(err)
			break
		}
	}
}

func (s *server) isPac(ip string, domain string) (bool, error) {
	if s.pac {
		if ip != "" {
			search, err := utils.IP2.MemorySearch(ip)
			if err != nil {
				log.Println(err)
				return false, err
			}
			if search.Country == "中国" || search.Country == "0" {
				return true, nil
			}
		}
		if domain != "" {
			//ipSimple, err := easy_dns.LookupIPSimple(domain, s.dns)
			ipSimple, err := lookIP(domain, s.dns)
			if err != nil {
				return false, nil
			}
			if ipSimple == "" {
				return false, nil
			}

			search, err := utils.IP2.MemorySearch(ipSimple)
			if err != nil {
				return false, err
			}
			if search.Country == "中国" || search.Country == "0" {
				return true, nil
			}
		}
	}
	return false, nil
}

func lookIP(domain string, dns string) (string, error) {
	var ip string
	var ttl uint32
	ipb, err := storage.Storage.Get(domain)
	if err == nil {
		return ipb.(string), nil
	}

	lookupIP, err := easy_dns.LookupIP(domain, dns)
	if err != nil {
		return "", err
	}
	for _, v := range lookupIP.Answers {
		if v.Header.Type == easy_dns.TypeA {
			ip = v.Body.GoString()
			ttl = v.Header.TTL
			storage.Storage.SetWithExpire(domain, ip, time.Second*time.Duration(ttl))
			break
		}
	}

	if ip == "" {
		return "", fmt.Errorf("not fund")
	}

	return ip, nil
}

func simple(client net.Conn, addr string) {
	addrs, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return
	}
	server, err := net.DialTCP("tcp", nil, addrs)
	if err != nil {
		return
	}

	if err := server.SetLinger(0); err != nil {
		return
	}

	defer server.Close()
	//进行转发
	go io.Copy(server, client)
	io.Copy(client, server)
}

type loginCreds struct {
	Username, Password string
}

func (c *loginCreds) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"username": c.Username,
		"password": c.Password,
	}, nil
}

func (c *loginCreds) RequireTransportSecurity() bool {
	return true
}

func NewClientTLSFromFile(certFile, serverNameOverride string) (credentials.TransportCredentials, error) {
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM([]byte(certFile)) {
		return nil, fmt.Errorf("credentials: failed to append certificates")
	}
	return credentials.NewTLS(&tls.Config{ServerName: serverNameOverride, RootCAs: cp}), nil
}
