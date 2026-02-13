package ipc

import (
	"encoding/json"
	"fmt"
	"net"
)

func SendRequest(c Command) error {
	socketPath, err := SocketPath()
	if err != nil {
		return err
	}
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return err
	}
	defer conn.Close()

	data, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if _, err := fmt.Fprintf(conn, "%s\n", data); err != nil {
		return fmt.Errorf("failed to write JSON to socket: %w", err)
	}
	return nil
}
