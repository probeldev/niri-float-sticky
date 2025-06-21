package niriwindows

import (
	"github.com/probeldev/niri-float-sticky/bash"
)

func GetWindowsList() ([]Window, error) {
	input, err := bash.RunCommand("niri msg windows")
	if err != nil {
		return nil, err
	}

	windows, err := ParseWindows(input)

	return windows, err
}

func GetFloatWinwods() ([]Window, error) {
	response := []Window{}

	windows, err := GetWindowsList()
	if err != nil {
		return response, err
	}

	for _, w := range windows {
		if w.IsFloating {
			response = append(response, w)
		}
	}

	return response, nil
}
