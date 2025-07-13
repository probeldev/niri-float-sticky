package niriwindows

import (
	"fmt"

	"github.com/probeldev/niri-float-sticky/bash"
)

// MoveWindowToWorkspace перемещает окно на указанный workspace
func MoveWindowToWorkspace(windowID uint64, workspaceID uint8) error {
	// Если windowID = 0, перемещается текущее фокусированное окно
	cmd := fmt.Sprintf("niri msg action move-window-to-workspace --window-id %d %d ", windowID, workspaceID)

	_, err := bash.RunCommand(cmd)
	return err
}
