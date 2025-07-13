// Package niriworkspaces
package niriworkspaces

type Workspace struct {
	Name                 string  `json:"name,omitempty"`
	Output               string  `json:"output,omitempty"`
	WorkspaceID          uint64  `json:"id"`
	WorkspaceOnMonitorID uint8   `json:"idx"`
	ActiveWindowID       *uint64 `json:"active_window_id"`
	IsUrgent             bool    `json:"is_urgent"`
	IsActive             bool    `json:"is_active"`
	IsFocused            bool    `json:"is_focused"`
}
