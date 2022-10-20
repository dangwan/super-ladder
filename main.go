package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"

	"github.com/dangwan/super-ladder/internal/log"
	"github.com/pkg/errors"
)

func handleClientRequest(conn net.Conn) error {
	if conn == nil {
		return errors.Wrap(nil, "Client conn is nil")
	}
	defer conn.Close()
	log.Debug(conn.LocalAddr().String(), " is reqeusting")
	var b [1024]byte
	n, err := conn.Read(b[:])
	if err != nil {
		fmt.Println(err)
		return err
	}
	var method, host string
	fmt.Sscanf(string(b[:bytes.IndexByte(b[:], '\n')]), "%s%s", &method, &host)
	log.Debug("method:", method, "host:", host)
	reqUrl, err := url.Parse(host)
	if err != nil {
		fmt.Println(err)
		return err
	}
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
	if method == "connect" {
		fmt.Fprint(conn, "http:/1.1 200 connection established\r\n")
	} else {
		svr.Write(b[:n])
	}
	go io.Copy(svr, conn)
	io.Copy(conn, svr)
	return nil
}
func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	fmt.Println("Start listen 8080")
	for {
		client, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go handleClientRequest(client)

	}
}
