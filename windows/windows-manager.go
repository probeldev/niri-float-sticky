package windows

type WindowsManager struct {
	floatingWindows map[uint64]struct{}
	manualOverride  map[uint64]bool
	autoStick       bool
}

func NewWindowsManager(autoStick bool) *WindowsManager {
	return &WindowsManager{
		floatingWindows: make(map[uint64]struct{}),
		manualOverride:  make(map[uint64]bool),
		autoStick:       autoStick,
	}
}

func (wm *WindowsManager) GetSticky() []uint64 {
	res := make([]uint64, 0, len(wm.floatingWindows))
	for winId := range wm.floatingWindows {
		if wm.IsSticky(winId) {
			res = append(res, winId)
		}
	}
	return res
}

func (wm *WindowsManager) IsSticky(id uint64) bool {
	if v, ok := wm.manualOverride[id]; ok {
		return v
	}
	_, auto := wm.floatingWindows[id]
	if auto && wm.autoStick {
		return true
	}
	return false
}

func (wm *WindowsManager) SetFloating(id uint64) {
	wm.floatingWindows[id] = struct{}{}
}

func (wm *WindowsManager) ResetFloating() {
	wm.floatingWindows = make(map[uint64]struct{})
}

func (wm *WindowsManager) SetManual(id uint64, sticky bool) {
	current := wm.IsSticky(id)
	if sticky == current {
		return
	}
	wm.manualOverride[id] = sticky
}

func (wm *WindowsManager) Remove(id uint64) {
	delete(wm.floatingWindows, id)
	delete(wm.manualOverride, id)
}
