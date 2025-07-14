package main

import (
	"flag"
	nirievents "github.com/probeldev/niri-float-sticky/niri-events"
	niriwindows "github.com/probeldev/niri-float-sticky/niri-windows"
	log "github.com/sirupsen/logrus"
)

func main() {
	var debug, allowForeignMonitors bool
	flag.BoolVar(&debug, "debug", false, "enable debug logging")
	flag.BoolVar(&allowForeignMonitors, "allow-moving-to-foreign-monitors", false, "allow moving to foreign monitors")
	flag.Parse()

	if debug {
		log.SetLevel(log.DebugLevel)
	}
	log.Info("Starting niri-float-sticky daemon...")

	events, err := nirievents.GetEventStream()
	if err != nil {
		log.Panic(err)
	}

	floatingWindows := make(map[uint64]struct{})
	workspacesMonitorMap := make(map[uint64]string)
	windowsMonitorMap := make(map[uint64]string)

	for event := range events {
		switch e := event.(type) {
		case *nirievents.WorkspaceActivatedEvent:
			log.Debugf("Workspace activated: %v", e)
			for windowID := range floatingWindows {
				if !allowForeignMonitors && windowsMonitorMap[windowID] != workspacesMonitorMap[e.Event.WorkspaceID] {
					log.Warnf("Ignore moving window %d to foreign monitor %s", windowID, workspacesMonitorMap[e.Event.WorkspaceID])
					continue
				}
				log.Debugf("Moving window %d to workspace %v", windowID, e.Event.WorkspaceID)
				err = niriwindows.MoveWindowToWorkspace(windowID, e.Event.WorkspaceID)
				if err != nil {
					log.Error(err)
				}
			}
		case *nirievents.WorkspacesChangedEvent:
			log.Debugf("Workspaces changed: %v", e)
			workspacesMonitorMap = make(map[uint64]string)
			for _, workspace := range e.Event.Workspaces {
				workspacesMonitorMap[workspace.WorkspaceID] = workspace.Output
			}
		case *nirievents.WindowsChangedEvent:
			log.Debugf("Windows changed: %v", e)
			floatingWindows = make(map[uint64]struct{})
			windowsMonitorMap = make(map[uint64]string)
			for _, window := range e.Event.Windows {
				if window.IsFloating && window.WorkspaceID != nil {
					floatingWindows[window.WindowID] = struct{}{}
					windowsMonitorMap[window.WindowID] = workspacesMonitorMap[*window.WorkspaceID]
				}
			}
		case *nirievents.WindowClosedEvent:
			log.Debugf("Window closed: %v", e)
			delete(floatingWindows, e.Event.WindowID)
		case *nirievents.WindowOpenedOrChangedEvent:
			log.Debugf("Window opened or changed: %v", e)
			if window := e.Event.Window; window.IsFloating && window.WorkspaceID != nil {
				floatingWindows[window.WindowID] = struct{}{}
				windowsMonitorMap[window.WindowID] = workspacesMonitorMap[*window.WorkspaceID]
			}
		}
	}
}
