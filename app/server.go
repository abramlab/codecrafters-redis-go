package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	// Uncomment this block to pass the first stage
	// "net"
	// "os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		exitWithError(fmt.Errorf("failed to bind to port 6379: %s", err))
	}

	conn, err := l.Accept()
	if err != nil {
		exitWithError(fmt.Errorf("failed to accept connection: %s", err))
	}

	for {
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			exitWithError(fmt.Errorf("failed to read from connection: %s", err))
		}
		pong := []byte("+PONG\r\n")
		_, err = conn.Write(pong)
		if err != nil {
			exitWithError(fmt.Errorf("failed to write to connection: %s", err))
		}
	}
}

func exitWithError(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}
