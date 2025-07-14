package niriwindows

import (
	"fmt"
	nirisocket "github.com/probeldev/niri-float-sticky/niri-socket"
)

// MoveWindowToWorkspace перемещает окно на указанный workspace
func MoveWindowToWorkspace(windowID, workspaceID uint64) error {
	// Если windowID = 0, перемещается текущее фокусированное окно
	action := fmt.Sprintf(`{"Action":{"MoveWindowToWorkspace":{"window_id":%d,"focus":false,"reference":{"Id":%d}}}}`,
		windowID, workspaceID,
	)
	socket := nirisocket.GetSocket()
	defer nirisocket.ReleaseSocket(socket)
	if err := socket.SendRequest(action); err != nil {
		return fmt.Errorf("failed to move window to workspace: %w", err)
	}
	return nil
}
