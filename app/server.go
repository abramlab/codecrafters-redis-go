package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

var store = newStorage()

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
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			exitWithError(fmt.Errorf("failed to read from connection: %s", err))
		}
		buf = buf[:n]

		fmt.Printf("Raw request: %q\n", string(buf))

		msg, err := ParseFromReader(bufio.NewReader(bytes.NewReader(buf)))
		if err != nil {
			fmt.Println("Parse request failed:", err.Error())
			write(conn, fmt.Sprintf("-ERR parse request failed: %q\r\n", string(buf)))
			continue
		}
		if msg.Type != MessageMutli {
			write(conn, fmt.Sprintf("-ERR unsupported request type: %q\r\n", string(buf)))
			continue
		}

		switch strings.ToUpper(string(msg.Multi[0].Bulk)) {
		case "PING":
			write(conn, "+PONG\r\n")
		case "ECHO":
			if len(msg.Multi) != 2 {
				write(conn, "-ERR wrong number of arguments for 'ECHO' command\r\n")
				continue
			}

			write(conn, fmt.Sprintf("+%s\r\n", msg.Multi[1].Bulk))
		case "SET":
			if len(msg.Multi) < 3 {
				write(conn, "-ERR wrong number of arguments for 'SET' command\r\n")
				continue
			}

			key, value := string(msg.Multi[1].Bulk), string(msg.Multi[2].Bulk)
			store.set(key, value)
			write(conn, "+OK\r\n")
		case "GET":
			if len(msg.Multi) != 2 {
				write(conn, "-ERR wrong number of arguments for 'GET' command\r\n")
				continue
			}

			key := string(msg.Multi[1].Bulk)
			value, ok := store.get(key)
			if !ok {
				write(conn, "$-1\r\n")
				continue
			}
			write(conn, fmt.Sprintf("$%d\r\n%s\r\n", len(value), value))
		default:
			write(conn, "-ERR unknown command\r\n")
		}
	}
}

func write(conn net.Conn, resp string) {
	_, err := conn.Write([]byte(resp))
	if err != nil {
		exitWithError(fmt.Errorf("failed to write to connection: %s", err))
	}
}

func exitWithError(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}
