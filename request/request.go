package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"
)

// HandleRequest 处理客户端请求
func HandleRequest(conn net.Conn) error {
	if conn == nil {
		return errors.New("client conn is nil")
	}
	defer conn.Close()
	fmt.Println(conn.LocalAddr().String(), " is requesting")
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
