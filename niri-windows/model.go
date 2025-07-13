package niriwindows

type Window struct {
	Title       string  `json:"title,omitempty"`
	AppID       string  `json:"app_id,omitempty"`
	WindowID    uint64  `json:"id"`
	PID         *int32  `json:"pid"`
	WorkspaceID *uint64 `json:"workspace_id"`
	IsFocused   bool    `json:"is_focused"`
	IsFloating  bool    `json:"is_floating"`
	IsUrgent    bool    `json:"is_urgent"`
}
