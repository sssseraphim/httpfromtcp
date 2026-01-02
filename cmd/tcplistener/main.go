package main

import (
	"fmt"
	"io"
	"net"
	"strings"
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
		lines := getLinesChannel(conection)
		for line := range lines {
			fmt.Printf("%s\n", line)
		}
		fmt.Println("Connection closed!")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	res := make(chan string)
	go func() {
		defer f.Close()
		defer close(res)
		b := make([]byte, 8)
		var line string
		for {
			_, err := f.Read(b)
			if err == io.EOF {
				break
			}
			s := string(b)
			parts := strings.Split(s, "\n")
			line += parts[0]
			for i := 1; i < len(parts); i++ {
				res <- line
				line = parts[i]
			}
		}
		if len(line) > 0 {
			res <- line
		}
	}()
	return res
}
