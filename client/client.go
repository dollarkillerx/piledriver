package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"plumber/utils"
	"strconv"
	"time"

	"github.com/dgraph-io/badger/v2"
	"github.com/dollarkillerx/easy_dns"
	"github.com/gorilla/websocket"
)

var localHost = flag.String("local_host", "127.0.0.1:9871", "Local Host")
var localUser = flag.String("local_user", "piledriver", "Local User")
var localPassword = flag.String("local_password", "piledriver", "Local Password")
var pacDns = flag.String("pac_dns", "8.8.8.8", "Pac Dns")
var pac = flag.Bool("pac", false, "pac")
var remoteHost = flag.String("remote_host", "127.0.0.1:8020", "Remote Host")
var remotePath = flag.String("remote_path", "/piledriver", "Remote PATH")
var token = flag.String("token", "piledriver", "token auth")

func main() {
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Llongfile)

	initStorage()

	// Local
	addr, err := net.ResolveTCPAddr("tcp", *localHost)
	if err != nil {
		log.Fatalln(err)
	}
	conn, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		accept, err := conn.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		c, err := initClient()
		if err != nil {
			//log.Println(err)
			continue
		}
		go c.accept(accept)
	}
}

type client struct {
	conn *websocket.Conn
}

func (c *client) accept(conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("%s\n", err)
			return
		}
	}()

	if conn == nil {
		return
	}

	defer conn.Close()

	var b [1024]byte
	_, err := conn.Read(b[:])
	if err != nil {
		return
	}

	if b[0] != 0x05 {
		return
	}

	conn.Write([]byte{0x05, 0x00})
	n, err := conn.Read(b[:])
	if err != nil {
		return
	}

	var domain bool
	// 解析目的地
	var host, port string
	switch b[3] {
	case 0x01: // ip
		host = net.IPv4(b[4], b[5], b[6], b[7]).String()
	case 0x03: // domain
		host = string(b[5 : n-2]) //b[4]表示域名的长度
		domain = true
	case 0x04: // ipv6
		return
	}

	isPac := usePac(host, *pac, domain)

	port = strconv.Itoa(int(b[n-2])<<8 | int(b[n-1]))
	addr := net.JoinHostPort(host, port)

	if isPac {
		simple(conn, addr)
		return
	}

	if _, err := conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}); err != nil { //响应客户端连接成功
		return
	}

	err = c.conn.WriteMessage(websocket.TextMessage, []byte(addr))
	if err != nil {
		log.Println(err)
		return
	}

	go copy1(c.conn, conn)
	copy2(conn, c.conn)
}

func initClient() (*client, error) {
	u := url.URL{Scheme: "wss", Host: *remoteHost, Path: *remotePath}

	dialer := &websocket.Dialer{TLSClientConfig: &tls.Config{RootCAs: nil, InsecureSkipVerify: true}}
	dial, _, err := dialer.Dial(u.String(), http.Header{"token": []string{*token}})
	if err != nil {
		return nil, err
	}

	return &client{conn: dial}, nil
}

func copy1(conn *websocket.Conn, client io.Reader) {
	for {
		var b [1024]byte
		read, err := client.Read(b[:])
		if err != nil {
			if err == io.EOF {
				conn.Close()
				break
			}
			break
		}

		if err := conn.WriteMessage(websocket.BinaryMessage, b[:read]); err != nil {
			break
		}
	}
}

func copy2(client io.Writer, conn *websocket.Conn) {
	for {
		_, byt, err := conn.ReadMessage()
		if err != nil {
			return
		}

		if _, err := client.Write(byt); err != nil {
			break
		}
	}
}

func usePac(host string, pac bool, website bool) bool {
	if !pac {
		return false
	}

	var ip string

	rb, err := Storage.Get([]byte(host))
	if err != nil {
		if website {
			lookupIP, err := lockDns(host, *pacDns)
			if err != nil {
				// 如果不存在则查询内网DNS
				lookupHost, err := net.LookupHost(host)
				if err != nil {
					log.Println(err)
					return false
				} else {
					if len(lookupHost) > 0 {
						ip = lookupHost[0]
					} else {
						return false
					}
				}
			}
			if len(lookupIP) > 0 {
				ip = lookupIP[0]
			} else {
				return false
			}

			Storage.Set([]byte(host), []byte(ip), 4*time.Hour)
		}
	} else {
		ip = string(rb)
	}

	// TODO: 检测IP
	search, err := utils.IP2.MemorySearch(ip)
	if err != nil {
		return false
	}
	if search.Country == "中国" {
		return true
	}

	return false
}

type storage struct {
	db *badger.DB
}

var Storage *storage

func initStorage() {
	open, err := badger.Open(badger.DefaultOptions("./dns.cache"))
	if err != nil {
		log.Fatalln(err)
	}

	Storage = &storage{db: open}
}

func (l *storage) Set(key, val []byte, tll time.Duration) error {
	return l.db.Update(func(txn *badger.Txn) error {
		entry := badger.NewEntry(key, val)
		if tll != 0 {
			entry = entry.WithTTL(tll)
		}
		return txn.SetEntry(entry)
	})
}

func (l *storage) Get(key []byte) (value []byte, err error) {
	err = l.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			value = val
			return nil
		})
	})
	if err != nil {
		return nil, err
	}

	return value, nil
}

func lockDns(domain string, dns string) ([]string, error) {
	lookupIP, err := easy_dns.LookupIP(domain, dns)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, v := range lookupIP.Answers {
		if v.Header.Type == easy_dns.TypeA {
			result = append(result, v.Body.GoString())
		}
	}

	return result, nil
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
