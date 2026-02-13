package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	arrayflag "github.com/probeldev/niri-float-sticky/array-flag"
	"github.com/probeldev/niri-float-sticky/ipc"
	nirievents "github.com/probeldev/niri-float-sticky/niri-events"
	niriwindows "github.com/probeldev/niri-float-sticky/niri-windows"
	"github.com/probeldev/niri-float-sticky/utils"
	"github.com/probeldev/niri-float-sticky/windows"
	log "github.com/sirupsen/logrus"
)

func sendCommand(ipcCmd string) error {
	validCommands := map[string]struct{}{
		"set_sticky":    {},
		"unset_sticky":  {},
		"toggle_sticky": {},
	}

	if _, ok := validCommands[ipcCmd]; !ok {
		return fmt.Errorf("invalid ipc command: %q", ipcCmd)
	}

	focusedWin, err := niriwindows.GetFocusedWindow()
	if err != nil {
		return fmt.Errorf("failed to get focused window: %v", err)
	}

	cmd := ipc.Command{
		Action:   ipcCmd,
		WindowID: focusedWin.WindowID,
	}

	if err := ipc.SendRequest(cmd); err != nil {
		return fmt.Errorf("failed to send command: %v", err)
	}

	log.Infof("Sent command '%s' for window %d", ipcCmd, focusedWin.WindowID)
	return nil

}

func main() {
	var debug, showVersion, allowForeignMonitors, disableAutoStick bool
	var appIds, titles arrayflag.ArrayFlag
	var ipcCmd string
	flag.StringVar(&ipcCmd, "ipc", "", "send IPC command to daemon: set_sticky, unset_sticky, toggle_sticky")
	flag.BoolVar(&disableAutoStick, "disable-auto-stick", false, "disable auto sticking for all windows")
	flag.BoolVar(&debug, "debug", false, "enable debug logging")
	flag.BoolVar(&showVersion, "version", false, "print version and exit")
	flag.BoolVar(&allowForeignMonitors, "allow-moving-to-foreign-monitors", false, "allow moving to foreign monitors")
	flag.Var(&appIds, "app-id", "only move floating windows with app-id matching given patterns")
	flag.Var(&titles, "title", "only move floating windows with title matching this pattern")
	flag.Parse()

	autoStickEnabled := !disableAutoStick

	if showVersion {
		fmt.Println(version)
		os.Exit(0)
	}
	if debug {
		log.SetLevel(log.DebugLevel)
	}
	log.SetOutput(os.Stdout)

	if ipcCmd != "" {
		err := sendCommand(ipcCmd)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	log.Info("Starting niri-float-sticky daemon...")

	events, err := nirievents.GetEventStream()
	if err != nil {
		log.Panic(err)
	}

	appIDPattern := utils.CombinePatterns(appIds)
	titlePattern := utils.CombinePatterns(titles)

	workspacesMonitorMap := make(map[uint64]string)
	windowsMonitorMap := make(map[uint64]string)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	windowsMgr := windows.NewWindowsManager(autoStickEnabled)
	cmdChan := make(chan ipc.Command)
	err = ipc.StartIPC(ctx, cmdChan)
	if err != nil {
		log.Fatalf("error to start ipc: %v", err)
	}

	for {
		select {
		case event := <-events:
			switch e := event.(type) {
			case *nirievents.WorkspaceActivatedEvent:
				log.Debugf("Workspace %d activated", e.Event.WorkspaceID)
				for _, windowID := range windowsMgr.GetSticky() {
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
				windowsMgr.ResetFloating()
				log.Debug("Floating windows cache have been reset")
				windowsMonitorMap = make(map[uint64]string)
				log.Debug("Windows to monitor bindings have been reset")
				for _, win := range e.Event.Windows {
					if win.IsFloating && win.WorkspaceID != nil && appIDPattern.MatchString(win.AppID) && titlePattern.MatchString(win.Title) {
						windowsMgr.SetFloating(win.WindowID)
						windowsMonitorMap[win.WindowID] = workspacesMonitorMap[*win.WorkspaceID]
						logf := log.WithFields(log.Fields{"app_id": win.AppID, "output": windowsMonitorMap[win.WindowID]})
						logf.Debugf("Window %d is now floating on %d workspace", win.WindowID, *win.WorkspaceID)
					}
				}
			case *nirievents.WindowClosedEvent:
				log.Debugf("Window %d is closed", e.Event.WindowID)
				windowsMgr.Remove(e.Event.WindowID)
			case *nirievents.WindowOpenedOrChangedEvent:
				win := e.Event.Window
				if win.IsFloating && win.WorkspaceID != nil && appIDPattern.MatchString(win.AppID) && titlePattern.MatchString(win.Title) {
					windowsMgr.SetFloating(win.WindowID)
					windowsMonitorMap[win.WindowID] = workspacesMonitorMap[*win.WorkspaceID]
					logf := log.WithFields(log.Fields{"app_id": win.AppID, "output": windowsMonitorMap[win.WindowID]})
					logf.Debugf("Window %d is now floating on %d workspace", win.WindowID, *win.WorkspaceID)
				} else if !win.IsFloating {
					windowsMgr.Remove(win.WindowID)
					log.WithField("app_id", win.AppID).Debugf("Window %d is now tiled mode", win.WindowID)
				}
			}
		case cmd := <-cmdChan:
			current := windowsMgr.IsSticky(cmd.WindowID)
			switch cmd.Action {
			case "set_sticky":
				windowsMgr.SetManual(cmd.WindowID, true)
			case "unset_sticky":
				windowsMgr.SetManual(cmd.WindowID, false)
			case "toggle_sticky":
				windowsMgr.SetManual(cmd.WindowID, !current)
			}
			isSticky := windowsMgr.IsSticky(cmd.WindowID)
			log.Infof("Window %d sticky state changed. Is sticky now: %v", cmd.WindowID, isSticky)
		case <-ctx.Done():
			log.Info("shutdown")
			return
		}
	}
}
