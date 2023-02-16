package main

import (
	"fmt"
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
	for {
		conn, err := l.Accept()
		if err != nil {
			exitWithError(fmt.Errorf("failed to accept connection: %s", err))
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
