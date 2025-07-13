package main

import (
	nirievents "github.com/probeldev/niri-float-sticky/niri-events"
	niriwindows "github.com/probeldev/niri-float-sticky/niri-windows"
	niriworkspaces "github.com/probeldev/niri-float-sticky/niri-workspaces"
	log "github.com/sirupsen/logrus"
)

func main() {
	events, err := nirievents.GetEventStream()
	if err != nil {
		log.Panic(err)
	}

	floatingWindows := make(map[uint64]struct{})

	for event := range events {
		switch e := event.(type) {
		case *nirievents.WorkspaceActivatedEvent:
			log.Debugf("Workspace activated: %v", e)
			for windowID := range floatingWindows {
				var workspaceOnMonitorID uint8
				workspaceOnMonitorID, err = niriworkspaces.GetCurrentWorkspace()
				if err != nil {
					log.Error(err)
				}
				err = niriwindows.MoveWindowToWorkspace(windowID, workspaceOnMonitorID)
				if err != nil {
					log.Error(err)
				}
			}
		case *nirievents.WindowsChangedEvent:
			log.Debugf("Windows changed: %v", e)
			for _, w := range e.Event.Windows {
				if w.IsFloating {
					floatingWindows[w.WindowID] = struct{}{}
				} else {
					delete(floatingWindows, w.WindowID)
				}
			}
		case *nirievents.WindowClosedEvent:
			log.Debugf("Window closed: %v", e)
			if _, ok := floatingWindows[e.Event.WindowID]; ok {
				delete(floatingWindows, e.Event.WindowID)
			}
		case *nirievents.WindowOpenedOrChangedEvent:
			log.Debugf("Window opened or changed: %v", e)
			if e.Event.Window.IsFloating {
				floatingWindows[e.Event.Window.WindowID] = struct{}{}
			} else {
				delete(floatingWindows, e.Event.Window.WindowID)
			}
		}
	}
}
