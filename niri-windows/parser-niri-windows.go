package niriwindows

import (
	"encoding/json"
	"fmt"
)

func ParseWindows(output []byte) ([]Window, error) {
	var windows []Window
	if err := json.Unmarshal(output, &windows); err != nil {
		return nil, fmt.Errorf("error unmarshalling windows: %w", err)
	}
	return windows, nil
}
