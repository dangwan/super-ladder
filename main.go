package main

import (
	"flag"
	"fmt"
	"net"
)

func main() {
	host := flag.String("host", "", "Please input host address")
	port := flag.String("port", "8080", "Please input port")
	flag.Parse()
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", *host, *port))
	if err != nil {
		panic(err)
	}
	fmt.Println("Start listen " + fmt.Sprintf("%s:%s", *host, *port))
	_ = listener
	//for {
	//	client, err := listener.Accept()
	//	if err != nil {
	//		panic(err)
	//	}
	//	go request.handleClientRequest(client)
	//
	//}
}
