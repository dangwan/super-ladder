package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

func handleClientRequest(conn net.Conn) error {
	if conn == nil {
		return errors.Wrap(nil, "Client conn is nil")
	}
	defer conn.Close()
	fmt.Println(conn.LocalAddr().String(), " is reqeusting")
	var b [1024]byte
	n, err := conn.Read(b[:])
	if err != nil {
		fmt.Println(err)
		return err
	}
	var method, host string
	fmt.Sscanf(string(b[:bytes.IndexByte(b[:], '\n')]), "%s%s", &method, &host)
	fmt.Println("method:", method, "host:", host)
	reqUrl, err := url.Parse(host)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("%+v", reqUrl)
	address := reqUrl.Host
	if reqUrl.Opaque == "443" {
		address = reqUrl.Scheme + ":443"
	} else {
		if !strings.Contains(reqUrl.Host, ":") {
			address = reqUrl.Host + ":80"
		}
	}
	svr, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if method == "CONNECT" {
		fmt.Fprint(conn, "http:/1.1 200 connection established\r\n")
	} else {
		svr.Write(b[:n])
	}
	go io.Copy(svr, conn)
	io.Copy(conn, svr)
	return nil
}
func main() {
	host := flag.String("host", "", "Please input host address")
	port := flag.String("port", "8080", "Please input port")
	flag.Parse()
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", *host, *port))
	if err != nil {
		panic(err)
	}
	fmt.Println("Start listen " + fmt.Sprintf("%s:%s", *host, *port))
	for {
		client, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go handleClientRequest(client)

	}
}
