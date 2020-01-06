package gomux

type killSession struct {
	targetSession string
}

type newSession struct {
	notAttachedToCurrentTerminal bool
	sessionName                  string
	windowName                   string
	workingDir                   string
}

type newWindow struct {
	targetWindow []string
	windowName   string
	workingDir   string
}

type renameWindow struct {
	targetWindow []string
	windowName   string
}

type splitWindow struct {
	horizontalSplit bool
	verticalSplit   bool
	targetPane      string
	workingDir      string
}

type selectWindow struct {
	targetWindow string
}

// This will silently kill a given tmux session
func (session killSession) String() (killCommand []string) {
	killCommand = append(killCommand, "kill-session", "-t", session.targetSession)
	return killCommand
}

func (session newSession) String() (newSessionCommand []string) {
	newSessionCommand = append(newSessionCommand, "new-session")

	if session.notAttachedToCurrentTerminal == true {
		newSessionCommand = append(newSessionCommand, "-d")
	}

	if session.sessionName != "" {
		newSessionCommand = append(newSessionCommand, "-s", session.sessionName)
	}

	if session.windowName != "" {
		newSessionCommand = append(newSessionCommand, "-n", session.windowName)
	}

	if session.workingDir != "" {
		newSessionCommand = append(newSessionCommand, "-c", session.workingDir)
	}

	return newSessionCommand
}

func (window splitWindow) String() (splitWindowCommand []string) {
	splitWindowCommand = append(splitWindowCommand, "split-window")

	if window.horizontalSplit == true {
		splitWindowCommand = append(splitWindowCommand, "-h")
	}

	if window.verticalSplit == true {
		splitWindowCommand = append(splitWindowCommand, "-v")
	}

	if window.targetPane != "" {
		splitWindowCommand = append(splitWindowCommand, "-t", window.targetPane)

	}

	if window.workingDir != "" {
		splitWindowCommand = append(splitWindowCommand, "-c", window.workingDir)
	}

	return splitWindowCommand
}

func (window newWindow) String() (newWindowCommand []string) {
	newWindowCommand = append(newWindowCommand, "new-window")

	if len(window.targetWindow) != 0 {
		for _, v := range window.targetWindow {
			newWindowCommand = append(newWindowCommand, v)
		}
	}

	if window.windowName != "" {
		newWindowCommand = append(newWindowCommand, "-n", window.windowName)
	}

	if window.workingDir != "" {
		newWindowCommand = append(newWindowCommand, "-c", window.workingDir)
	}

	return newWindowCommand
}

func (window renameWindow) String() (renameWindowCommand []string) {
	renameWindowCommand = append(renameWindowCommand, "rename-window")

	if len(window.targetWindow) != 0 {
		for _, v := range window.targetWindow {
			renameWindowCommand = append(renameWindowCommand, v)
		}
	}

	if window.windowName != "" {
		renameWindowCommand = append(renameWindowCommand, window.windowName)
	}

	return renameWindowCommand
}

func (window selectWindow) String() (selectWindowCommand []string) {
	selectWindowCommand = append(selectWindowCommand, "select-window")

	if window.targetWindow != "" {
		selectWindowCommand = append(selectWindowCommand, "-t", window.targetWindow)
	}

	return selectWindowCommand
}
