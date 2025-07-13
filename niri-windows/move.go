package niriwindows

import (
	"fmt"

	"github.com/probeldev/niri-float-sticky/bash"
)

// MoveWindowToWorkspace перемещает окно на указанный workspace
func MoveWindowToWorkspace(windowID, workspaceID uint64) error {
	// Если windowID = 0, перемещается текущее фокусированное окно
	action := fmt.Sprintf(`{"Action":{"MoveWindowToWorkspace":{"window_id":%d,"focus":false,"reference":{"Id":%d}}}}`,
		windowID, workspaceID,
	)
	cmd := fmt.Sprintf("echo '%s' | nc -w 0 -U $NIRI_SOCKET", action)

	_, err := bash.RunCommand(cmd)
	return err
}
