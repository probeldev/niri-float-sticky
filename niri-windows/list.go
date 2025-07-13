// Package niriwindows
package niriwindows

import (
	"github.com/probeldev/niri-float-sticky/bash"
)

func GetWindowsList() ([]Window, error) {
	output, err := bash.RunCommand("niri msg --json windows")
	if err != nil {
		return nil, err
	}

	windows, err := ParseWindows(output)

	return windows, err
}

func GetFloatWindows() ([]Window, error) {
	windows, err := GetWindowsList()
	if err != nil {
		return nil, err
	}

	response := make([]Window, 0, len(windows))

	for _, w := range windows {
		if w.IsFloating {
			response = append(response, w)
		}
	}

	return response, nil
}
