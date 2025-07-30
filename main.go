package main

import (
	"flag"
	"fmt"
	"os"

	nirievents "github.com/probeldev/niri-float-sticky/niri-events"
	niriwindows "github.com/probeldev/niri-float-sticky/niri-windows"
	log "github.com/sirupsen/logrus"
)

func main() {
	var debug, showVersion, allowForeignMonitors bool
	flag.BoolVar(&debug, "debug", false, "enable debug logging")
	flag.BoolVar(&showVersion, "version", false, "print version and exit")
	flag.BoolVar(&allowForeignMonitors, "allow-moving-to-foreign-monitors", false, "allow moving to foreign monitors")
	flag.Parse()

	if showVersion {
		fmt.Println(version)
		os.Exit(0)
	}
	if debug {
		log.SetLevel(log.DebugLevel)
	}
	log.SetOutput(os.Stdout)
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
			log.Debugf("Workspace %d activated", e.Event.WorkspaceID)
			for windowID := range floatingWindows {
				if !allowForeignMonitors && windowsMonitorMap[windowID] != workspacesMonitorMap[e.Event.WorkspaceID] {
					log.Warnf("Ignore moving window %d to foreign monitor %s", windowID, workspacesMonitorMap[e.Event.WorkspaceID])
					continue
				}
				log.Debugf("Moving window %d to workspace %d", windowID, e.Event.WorkspaceID)
				err = niriwindows.MoveWindowToWorkspace(windowID, e.Event.WorkspaceID)
				if err != nil {
					log.Error(err)
				}
			}
		case *nirievents.WorkspacesChangedEvent:
			log.Debug("Workspaces to monitor bindings have been reset")
			workspacesMonitorMap = make(map[uint64]string)
			for _, workspace := range e.Event.Workspaces {
				workspacesMonitorMap[workspace.WorkspaceID] = workspace.Output
				log.WithField("output", workspace.Output).Debugf("Workspace %d binded to monitor", workspace.WorkspaceID)
			}
		case *nirievents.WindowsChangedEvent:
			floatingWindows = make(map[uint64]struct{})
			log.Debug("Floating windows cache have been reset")
			windowsMonitorMap = make(map[uint64]string)
			log.Debug("Windows to monitor bindings have been reset")
			for _, win := range e.Event.Windows {
				if win.IsFloating && win.WorkspaceID != nil {
					floatingWindows[win.WindowID] = struct{}{}
					windowsMonitorMap[win.WindowID] = workspacesMonitorMap[*win.WorkspaceID]
					logf := log.WithFields(log.Fields{"app_id": win.AppID, "output": windowsMonitorMap[win.WindowID]})
					logf.Debugf("Window %d is now floating on %d workspace", win.WindowID, *win.WorkspaceID)
				}
			}
		case *nirievents.WindowClosedEvent:
			log.Debugf("Window %d is closed", e.Event.WindowID)
			delete(floatingWindows, e.Event.WindowID)
		case *nirievents.WindowOpenedOrChangedEvent:
			win := e.Event.Window
			if win.IsFloating && win.WorkspaceID != nil {
				floatingWindows[win.WindowID] = struct{}{}
				windowsMonitorMap[win.WindowID] = workspacesMonitorMap[*win.WorkspaceID]
				logf := log.WithFields(log.Fields{"app_id": win.AppID, "output": windowsMonitorMap[win.WindowID]})
				logf.Debugf("Window %d is now floating on %d workspace", win.WindowID, *win.WorkspaceID)
			} else if !win.IsFloating {
				delete(floatingWindows, win.WindowID)
				log.WithField("app_id", win.AppID).Debugf("Window %d is now tiled mode", win.WindowID)
			}
		}
	}
}
