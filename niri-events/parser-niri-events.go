package nirievents

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/probeldev/niri-float-sticky/bash"
	log "github.com/sirupsen/logrus"
)

func GetEventStream() (<-chan any, error) {
	linesCh, err := bash.RunAndListenCommand("niri msg --json event-stream")
	if err != nil {
		return nil, fmt.Errorf("error getting event-stream: %w", err)
	}

	eventStreamCh := make(chan any)

	go func() {
		defer close(eventStreamCh)

		for line := range linesCh {
			if len(line) < 2 {
				continue
			}
			events := line[bytes.IndexByte(line, '{')+1 : bytes.IndexByte(line, ':')]
			var event any
			switch {
			case bytes.Equal(events, []byte("\"WorkspaceActivated\"")):
				event = &WorkspaceActivatedEvent{}
			case bytes.Equal(events, []byte("\"WindowsChanged\"")):
				event = &WindowsChangedEvent{}
			case bytes.Equal(events, []byte("\"WindowClosed\"")):
				event = &WindowClosedEvent{}
			case bytes.Equal(events, []byte("\"WindowOpenedOrChanged\"")):
				event = &WindowOpenedOrChangedEvent{}
			default:
				continue
			}
			if err = json.Unmarshal(line, &event); err != nil {
				log.Errorf("error unmarshalling event: %v", err)
				continue
			}
			eventStreamCh <- event
		}
	}()

	return eventStreamCh, nil
}
