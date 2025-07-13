package main

import (
	"flag"
	nirievents "github.com/probeldev/niri-float-sticky/niri-events"
	niriwindows "github.com/probeldev/niri-float-sticky/niri-windows"
	log "github.com/sirupsen/logrus"
)

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "enable debug logging")
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

	for event := range events {
		switch e := event.(type) {
		case *nirievents.WorkspaceActivatedEvent:
			log.Debugf("Workspace activated: %v", e)
			for windowID := range floatingWindows {
				log.Debugf("Moving window %d to workspace %v", windowID, e.Event.WorkspaceID)
				err = niriwindows.MoveWindowToWorkspace(windowID, e.Event.WorkspaceID)
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
			delete(floatingWindows, e.Event.WindowID)
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
