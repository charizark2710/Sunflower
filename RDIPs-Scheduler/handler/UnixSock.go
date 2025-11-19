package handler

import (
	LogConstant "RDIPs-Scheduler/constant/LogConst"
	connection "RDIPs-Scheduler/handler/Connection"
	"RDIPs-Scheduler/utils"
	"context"
	"fmt"
	"net"
	"path/filepath"
	"sync"
	"time"
)

var sockPool *connection.Pool = &connection.Pool{}
var mu sync.Mutex

func InitializeUnixSock(path string) error {
	factoryFn := func() (interface{}, error) {
		mu.Lock()
		defer mu.Unlock()
		conn := connect(path)
		time.Sleep(2 * time.Second)
		return conn, nil
	}

	closeFn := func(conn interface{}) error {
		unixConn, ok := conn.(*net.Conn)
		if !ok {
			return fmt.Errorf("%v", "wrong unix connection format")
		}
		(*unixConn).Close()
		return nil
	}

	pingFn := func(conn interface{}) error {
		unixConn, ok := conn.(*net.Conn)
		if !ok {
			return fmt.Errorf("%v", "wrong unix connection format")
		}
		// Test if connection is still alive
		(*unixConn).SetWriteDeadline(time.Now().Add(2 * time.Second))
		_, err := (*unixConn).Write([]byte{})
		if err != nil {
			return fmt.Errorf("connection dead: %v", err)
		}
		(*unixConn).SetWriteDeadline(time.Time{}) // Clear deadline
		return nil
	}

	poolData := connection.PoolData{
		FactoryFn: factoryFn,
		CloseFn:   closeFn,
		PingFn:    pingFn,
	}

	err := sockPool.FillPool(poolData)

	if err != nil {
		return err
	}

	utils.Log(LogConstant.Info, "Finish Setup Unix Sock Pool, err:", err)
	return err
}

func GetUnixSock() *net.Conn {
	conn, _, err := sockPool.Get()
	if err != nil {
		utils.Log(LogConstant.Error, err)
		return nil
	}
	unixConn, ok := conn.(*net.Conn)
	if !ok {
		utils.Log(LogConstant.Error, "wrong unix connection format")
		return nil
	}
	return unixConn
}

func PutUnixSock(conn *net.Conn) {
	sockPool.Release(conn, true, context.Background())
}

func connect(path string) *net.Conn {
	for {
		conn, err := net.Dial("unix", path)
		if err == nil {
			return &conn
		}
		time.Sleep(100 * time.Millisecond)
	}
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
