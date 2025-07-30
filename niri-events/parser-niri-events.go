package nirievents

import (
	"bytes"
	"encoding/json"
	"fmt"
	nirisocket "github.com/probeldev/niri-float-sticky/niri-socket"
	log "github.com/sirupsen/logrus"
)

func GetEventStream() (<-chan any, error) {
	eventStreamCh := make(chan any)
	socket := nirisocket.GetSocket()

	go func() {
		defer nirisocket.ReleaseSocket(socket)
		defer socket.Close()
		defer close(eventStreamCh)

		for line := range socket.RecvStream() {
			if len(line) < 2 {
				continue
			}
			events := line[1:bytes.IndexByte(line, ':')]
			var event any
			switch {
			case bytes.Equal(events, []byte("\"WorkspaceActivated\"")):
				event = &WorkspaceActivatedEvent{}
			case bytes.Equal(events, []byte("\"WorkspacesChanged\"")):
				event = &WorkspacesChangedEvent{}
			case bytes.Equal(events, []byte("\"WindowsChanged\"")):
				event = &WindowsChangedEvent{}
			case bytes.Equal(events, []byte("\"WindowClosed\"")):
				event = &WindowClosedEvent{}
			case bytes.Equal(events, []byte("\"WindowOpenedOrChanged\"")):
				event = &WindowOpenedOrChangedEvent{}
			default:
				continue
			}
			if err := json.Unmarshal(line, &event); err != nil {
				log.Errorf("error unmarshalling event: %v", err)
				continue
			}
			eventStreamCh <- event
		}
	}()

	if err := socket.SendRequest("\"EventStream\""); err != nil {
		return nil, fmt.Errorf("error during request event-stream: %w", err)
	}

	return eventStreamCh, nil
}
