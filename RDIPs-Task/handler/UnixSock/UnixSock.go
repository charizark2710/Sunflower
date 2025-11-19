package handler

import (
	"fmt"
	"net"
	"path/filepath"
	"time"
)

func Connect(path string) net.Conn {
	time.Sleep(2 * time.Second)
	socketPath, _ := filepath.Abs(path)
	fmt.Println(socketPath)
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		panic(err)
	}

	return conn
}

func Listen(path string) net.Listener {
	socketPath, _ := filepath.Abs(path)
	fmt.Println(socketPath)
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		panic(err)
	}

	return listener
}
