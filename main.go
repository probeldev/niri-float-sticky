package main

import (
	"log"
	"time"

	niriwindows "github.com/probeldev/niri-float-sticky/niri-float-sticky/niri-windows"
	niriworkspaces "github.com/probeldev/niri-float-sticky/niri-float-sticky/niri-workspaces"
)

func main() {

	for {
		windows, err := niriwindows.GetFloatWinwods()
		log.Println(err)
		log.Println(windows)

		workspace, err := niriworkspaces.GetCurrentWorkspace()

		log.Println(err)
		log.Println(workspace)

		for _, w := range windows {
			err := niriwindows.MoveWindowToWorkspace(w.ID, workspace)
			if err != nil {
				log.Println(err)
			}
		}
		time.Sleep(200 * time.Millisecond)

	}

}
