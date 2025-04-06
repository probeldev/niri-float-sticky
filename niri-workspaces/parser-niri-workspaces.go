package niriworkspaces

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/probeldev/niri-float-sticky/bash"
)

// GetCurrentWorkspace возвращает ID текущего workspace (отмеченного *)
func GetCurrentWorkspace() (int, error) {
	// Выполняем команду и получаем вывод
	output, err := bash.RunCommand("niri msg workspaces")
	if err != nil {
		return 0, fmt.Errorf("failed to get workspaces: %v", err)
	}

	// Парсим вывод
	return parseCurrentWorkspace(output)
}

// parseCurrentWorkspace парсит вывод команды и находит текущий workspace
func parseCurrentWorkspace(output string) (int, error) {
	lines := strings.Split(output, "\n")
	re := regexp.MustCompile(`\*\s*(\d+)`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if matches := re.FindStringSubmatch(line); matches != nil {
			workspaceID, err := strconv.Atoi(matches[1])
			if err != nil {
				return 0, fmt.Errorf("invalid workspace ID: %v", err)
			}
			return workspaceID, nil
		}
	}

	return 0, fmt.Errorf("current workspace not found")
}
