package nirisocket

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
)

type Socket struct {
	conn net.Conn
}

var pool = sync.Pool{
	New: func() any {
		conn, err := net.Dial("unix", os.Getenv("NIRI_SOCKET"))
		if err != nil {
			log.Panicf("failed to connect to NIRI_SOCKET: %v", err)
		}
		return &Socket{conn: conn}
	},
}

func GetSocket() *Socket {
	return pool.Get().(*Socket)
}

func ReleaseSocket(socket *Socket) {
	pool.Put(socket)
}

func (s *Socket) SendRequest(req string) error {
	_, err := fmt.Fprintf(s.conn, "%s\n", req)
	return err
}

func (s *Socket) RecvStream() <-chan []byte {
	linesCh := make(chan []byte)

	go func() {
		defer s.conn.Close()
		defer close(linesCh)

		decoder := json.NewDecoder(s.conn)

		for {
			var raw json.RawMessage
			if err := decoder.Decode(&raw); err != nil {
				if err == io.EOF {
					return
				}
				log.Errorf("error decode response from NIRI_SOCKET: %v", err)
				return
			}

			linesCh <- raw
		}
	}()

	return linesCh
}

func (s *Socket) Close() {
	_ = s.conn.Close()
}
