package main

import (
	"fmt"
	"net"

	"github.com/sssseraphim/httpfromtcp/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Printf("Cant listen to tcp: %v", err)
		return
	}
	defer listener.Close()
	for {
		conection, err := listener.Accept()
		if err != nil {
			fmt.Print(err)
			return
		}
		fmt.Println("Connection success!")
		req, err := request.RequestFromReader(conection)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n", req.RequestLine.Method, req.RequestLine.RequestTarget, req.RequestLine.HttpVersion)
		fmt.Println("Connection closed!")
	}
}
