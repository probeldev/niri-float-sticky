package nirisocket

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"sync"
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
		defer func() { _ = s.conn.Close() }()
		defer close(linesCh)

		scanner := bufio.NewScanner(s.conn)
		for scanner.Scan() {
			linesCh <- scanner.Bytes()
		}
		if err := scanner.Err(); err != nil {
			log.Errorf("error scan response from NIRI_SOCKET: %v", err)
		}
	}()

	return linesCh
}

func (s *Socket) Close() {
	_ = s.conn.Close()
}
