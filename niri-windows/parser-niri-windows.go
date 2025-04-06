package niriwindows

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Window struct {
	ID          int
	Title       string
	AppID       string
	IsFloating  bool
	PID         int
	WorkspaceID int
	IsFocused   bool
}

func ParseWindows(input string) ([]Window, error) {
	var windows []Window
	lines := strings.Split(input, "\n")
	var currentWindow *Window

	// Регулярные выражения для парсинга
	idRe := regexp.MustCompile(`Window ID (\d+)(: \(focused\))?`)
	titleRe := regexp.MustCompile(`\s*Title: "(.*)"`)
	appIDRe := regexp.MustCompile(`\s*App ID: "(.*)"`)
	floatingRe := regexp.MustCompile(`\s*Is floating: (yes|no)`)
	pidRe := regexp.MustCompile(`\s*PID: (\d+)`)
	workspaceRe := regexp.MustCompile(`\s*Workspace ID: (\d+)`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if matches := idRe.FindStringSubmatch(line); matches != nil {
			// Сохраняем предыдущее окно, если оно есть
			if currentWindow != nil {
				windows = append(windows, *currentWindow)
			}

			id, err := strconv.Atoi(matches[1])
			if err != nil {
				return nil, fmt.Errorf("invalid window ID: %v", err)
			}

			currentWindow = &Window{
				ID:        id,
				IsFocused: matches[2] != "",
			}
		} else if currentWindow != nil {
			switch {
			case titleRe.MatchString(line):
				matches := titleRe.FindStringSubmatch(line)
				currentWindow.Title = matches[1]
			case appIDRe.MatchString(line):
				matches := appIDRe.FindStringSubmatch(line)
				currentWindow.AppID = matches[1]
			case floatingRe.MatchString(line):
				matches := floatingRe.FindStringSubmatch(line)
				currentWindow.IsFloating = matches[1] == "yes"
			case pidRe.MatchString(line):
				matches := pidRe.FindStringSubmatch(line)
				pid, err := strconv.Atoi(matches[1])
				if err != nil {
					return nil, fmt.Errorf("invalid PID: %v", err)
				}
				currentWindow.PID = pid
			case workspaceRe.MatchString(line):
				matches := workspaceRe.FindStringSubmatch(line)
				workspace, err := strconv.Atoi(matches[1])
				if err != nil {
					return nil, fmt.Errorf("invalid workspace ID: %v", err)
				}
				currentWindow.WorkspaceID = workspace
			}
		}
	}

	// Добавляем последнее окно, если оно есть
	if currentWindow != nil {
		windows = append(windows, *currentWindow)
	}

	return windows, nil
}
