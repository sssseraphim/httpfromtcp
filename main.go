package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		fmt.Println(err)
	}
	lines := getLinesChannel(f)
	for line := range lines {
		fmt.Printf("read: %s\n", line)
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
