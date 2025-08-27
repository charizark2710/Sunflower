package socket

import (
	"encoding/json"
	"io"
	"net"
	"time"
)

func GuessHandler(sock net.Listener) {
	for {
		conn, err := sock.Accept()
		if err != nil {
			panic(err)
		}
		var result []byte
		buf := make([]byte, 4096)
		go func(c net.Conn) {
			for {
				c.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
				n, err := c.Read(buf)
				if n > 0 {
					result = append(result, buf[:n]...)
				}
				if err != nil {
					if err == io.EOF || len(result) > 0 {
						break
					}
					errorJson, _ := json.Marshal(map[string]string{"error": "unexpected read error"})
					c.Write(errorJson)
					return
				}
			}
			defer c.Close()
		}(conn)

	}
}
