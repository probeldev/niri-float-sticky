package ipc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

type Command struct {
	Action   string `json:"action"`
	WindowID uint64 `json:"window_id"`
}

func SocketPath() (string, error) {
	dir := os.Getenv("XDG_RUNTIME_DIR")
	if dir == "" {
		return "", fmt.Errorf("XDG_RUNTIME_DIR is not set")
	}
	return filepath.Join(dir, "niri-float-sticky.sock"), nil
}

func StartIPC(ctx context.Context, cmdChan chan<- Command) error {
	socketPath, err := SocketPath()
	if err != nil {
		return err
	}
	if err := os.Remove(socketPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove old socket: %v", err)
	}

	l, err := net.Listen("unix", socketPath)
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		l.Close()
		os.Remove(socketPath)
	}()

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return
				}

				log.Errorf("accept error: %v", err)
				continue
			}

			go handleConn(conn, cmdChan)
		}
	}()
	return nil
}

func handleConn(c net.Conn, cmdChan chan<- Command) {
	defer c.Close()
	var req Command
	decoder := json.NewDecoder(io.LimitReader(c, 512))
	if err := decoder.Decode(&req); err != nil {
		log.Errorf("failed to handle connection: %v", err)
		return
	}

	cmdChan <- req
}
