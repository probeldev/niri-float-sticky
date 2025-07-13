package nirievents

import (
	niriwindows "github.com/probeldev/niri-float-sticky/niri-windows"
)

type WorkspaceActivatedEvent struct {
	Event struct {
		WorkspaceID uint64 `json:"id"`
		Focused     bool   `json:"focused"`
	} `json:"WorkspaceActivated"`
}

type WindowsChangedEvent struct {
	Event struct {
		Windows []niriwindows.Window `json:"windows"`
	} `json:"WindowsChanged"`
}

type WindowClosedEvent struct {
	Event struct {
		WindowID uint64 `json:"id"`
	} `json:"WindowClosed"`
}

type WindowOpenedOrChangedEvent struct {
	Event struct {
		Window niriwindows.Window `json:"window"`
	} `json:"WindowOpenedOrChanged"`
}
