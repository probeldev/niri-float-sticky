package niriworkspaces

import (
	"encoding/json"
	"fmt"
	"github.com/probeldev/niri-float-sticky/bash"
)

func GetWorkspaces() ([]Workspace, error) {
	// Выполняем команду и получаем вывод
	output, err := bash.RunCommand("niri msg --json workspaces")
	if err != nil {
		return nil, fmt.Errorf("failed to get workspaces: %v", err)
	}

	// Парсим вывод
	return parseCurrentWorkspace(output)
}

// GetCurrentWorkspace возвращает ID текущего workspace (отмеченного *)
func GetCurrentWorkspace() (Workspace, error) {
	workspaces, err := GetWorkspaces()
	if err != nil {
		return Workspace{}, err
	}
	for _, w := range workspaces {
		if w.IsFocused {
			return w, nil
		}
	}
	return Workspace{}, fmt.Errorf("current workspace not found")
}

// parseCurrentWorkspace парсит вывод команды и находит текущий workspace
func parseCurrentWorkspace(output []byte) ([]Workspace, error) {
	var workspaces []Workspace
	if err := json.Unmarshal(output, &workspaces); err != nil {
		return nil, fmt.Errorf("error unmarshalling workspaces: %w", err)
	}
	return workspaces, nil
}
