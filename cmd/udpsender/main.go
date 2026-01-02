package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":42069")
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		s, err := reader.ReadString(byte('\n'))
		if err != nil {
			fmt.Println(err)
		}
		_, err = conn.Write([]byte(s))
		if err != nil {
			fmt.Println(err)
		}
	}

}
