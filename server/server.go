package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"plumber/rpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
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
const key = `
-----BEGIN EC PARAMETERS-----
BgUrgQQAIg==
-----END EC PARAMETERS-----
-----BEGIN EC PRIVATE KEY-----
MIGkAgEBBDDrNKgmBBwFBR/NgUxoVDzYeyl7wb1dqmehUSvzBPYLd3SZp/euqVli
hIOSvDX4DrqgBwYFK4EEACKhZANiAATvf5qn0hKp56iFULTyXfTNRVMf8mBgQR8G
iJlhMNs8SE/128T2lD0UFDMrAVxKxo/rHsP5ORiP/uc9NnK721dVlmcH40XGXtW2
BbThJeCFdNO1ife81fuqzxWx4oSIGWY=
-----END EC PRIVATE KEY-----
`

func main() {
	//log.SetFlags(log.Llongfile | log.LstdFlags)
	addr := "0.0.0.0:443"
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	defer listen.Close()

	creds, err := NewServerTLSFromFile(pem, key)
	if err != nil {
		log.Fatalln(err)
	}

	srv := grpc.NewServer(
		grpc.Creds(creds),
		grpc.StreamInterceptor(streamInterceptor("user", "H40XGXtW2")),
		grpc.UnaryInterceptor(unaryInterceptor("user", "H40XGXtW2")),
	)
	rpc.RegisterPlumberServer(srv, &server{})
	reflection.Register(srv)

	log.Println(addr)
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
		//log.Println(err)
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
			//log.Println(err)
			break
		}

		if err := server.Send(&rpc.PlumberResponse{Data: b[:read]}); err != nil {
			//log.Println(err)
			break
		}
	}
}

func copy1(server io.Writer, ser rpc.Plumber_PlumberServer) {
	for {
		recv, err := ser.Recv()
		if err != nil {
			//log.Println(err)
			break
		}
		if recv.Over {
			break
		}
		if _, err := server.Write(recv.Data); err != nil {
			//log.Println(err)
			break
		}
	}
}

func streamInterceptor(username string, password string) func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if err := authorize(username, password)(stream.Context()); err != nil {
			return err
		}

		return handler(srv, stream)
	}
}

func unaryInterceptor(username string, password string) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if err := authorize(username, password)(ctx); err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

func authorize(username string, password string) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if len(md["username"]) > 0 && md["username"][0] == username &&
				len(md["password"]) > 0 && md["password"][0] == password {
				return nil
			}

			return fmt.Errorf("AccessDeniedErr")
		}

		return fmt.Errorf("EmptyMetadataErr")
	}
}

func NewServerTLSFromFile(certFile, keyFile string) (credentials.TransportCredentials, error) {
	cert, err := LoadX509KeyPair([]byte(certFile), []byte(keyFile))
	if err != nil {
		return nil, err
	}
	return credentials.NewTLS(&tls.Config{Certificates: []tls.Certificate{cert}}), nil
}

func LoadX509KeyPair(cert, key []byte) (tls.Certificate, error) {
	return tls.X509KeyPair(cert, key)
}
