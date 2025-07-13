package niriworkspaces

import (
	"encoding/json"
	"fmt"
	"github.com/probeldev/niri-float-sticky/bash"
)

// GetCurrentWorkspace возвращает ID текущего workspace (отмеченного *)
func GetCurrentWorkspace() (uint8, error) {
	// Выполняем команду и получаем вывод
	output, err := bash.RunCommand("niri msg --json workspaces")
	if err != nil {
		return 0, fmt.Errorf("failed to get workspaces: %v", err)
	}

	// Парсим вывод
	return parseCurrentWorkspace(output)
}

// parseCurrentWorkspace парсит вывод команды и находит текущий workspace
func parseCurrentWorkspace(output []byte) (uint8, error) {
	var workspaces []Workspace
	if err := json.Unmarshal(output, &workspaces); err != nil {
		return 0, fmt.Errorf("error unmarshalling workspaces: %w", err)
	}
	for _, w := range workspaces {
		if w.IsFocused {
			return w.WorkspaceOnMonitorID, nil
		}
	}
	return 0, fmt.Errorf("current workspace not found")
}
