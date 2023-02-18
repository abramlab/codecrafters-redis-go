package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
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

		echo, err := parseSimpleEchoCommand(buf)
		if err != nil {
			errResp := []byte(fmt.Sprintf("-ERR %s: %q\r\n", err, string(buf)))
			if _, err = conn.Write(errResp); err != nil {
				exitWithError(fmt.Errorf("failed to write to connection: %s", err))
			}
			continue
		}

		if _, err := conn.Write(echo.response()); err != nil {
			exitWithError(fmt.Errorf("failed to write to connection: %s", err))
		}
	}
}

func exitWithError(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}

type echoCommand struct {
	arg string
}

func (c *echoCommand) response() []byte {
	return []byte(fmt.Sprintf("+%s\r\n", c.arg))
}

func parseSimpleEchoCommand(raw []byte) (*echoCommand, error) {
	parts := bytes.Split(bytes.TrimSuffix(raw, []byte("\r\n")), []byte("\r\n"))

	if len(parts) != 5 {
		return nil, fmt.Errorf("unsupported command")
	}

	if parts[0][0] != '*' {
		return nil, fmt.Errorf("unsupported command")
	}
	if !(bytes.Equal(parts[1], []byte("$4")) && bytes.Equal(parts[2], []byte("ECHO"))) {
		return nil, fmt.Errorf("unsupported command")
	}
	return &echoCommand{arg: string(parts[4])}, nil
}
