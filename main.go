package gomux

import (
	"fmt"
	"os/exec"
	"strings"
)

// Pane is the struct for a single pane nested within a window inside a tmux session
type Pane struct {
	Number   int
	commands []string
	window   *Window
}

// Session represents a tmux session.
// Use the method NewSession to create a Session instance.
type Session struct {
	Name             string
	Directory        string
	Windows          []*Window
	directory        string
	NextWindowNumber int
}

// SessionAttr is an struct that holds attributes for a given tmux session
type SessionAttr struct {
	Name      string
	Directory string
}

// SplitAttr represents the attributes that can be set within a given Split
type SplitAttr struct {
	Directory string
}

// Window represents a single tmux window. You usually should not create an instance of Window directly.
type Window struct {
	Number        int
	Name          string
	Directory     string
	Session       *Session
	Panes         []*Pane
	SplitCommands []string
}

// WindowAttr represents the attributes that can be set for a given window
type WindowAttr struct {
	Name      string
	Directory string
}

// NewPane returns a pane object, which is used to create a pane within a tmux window
func NewPane(number int, window *Window) (newPane *Pane) {
	p := &Pane{
		Number:   number,
		commands: make([]string, 0),
		window:   window,
	}

	return p
}

// Exec executes the provided command in the context of a the current pane
func (p *Pane) Exec(command string) (err error) {
	rawCommand := []string{"send-keys", "-t", p.getTargetName(), strings.Replace(command, "\"", "\\\"", -1), "C-m"}
	tmuxCmd := exec.Command("tmux", rawCommand...)
	return tmuxCmd.Run()
}

// Vsplit horizontally splits the view of the current window to include the current pane
func (p *Pane) Vsplit() (splitPane *Pane) {
	rawCmd := splitWindow{
		horizontalSplit: true,
		targetPane:      p.getTargetName()}

	cmd := exec.Command("tmux", rawCmd.String()...)
	cmd.Start()

	return p.window.AddPane(p.Number + 1)
}

// VsplitWAttr sets attributes of the current pane to be split horizontally
func (p *Pane) VsplitWAttr(attr SplitAttr) (splitPane *Pane, err error) {
	var c string
	if attr.Directory != "" {
		c = attr.Directory
	} else if p.window.Directory != "" {
		c = p.window.Directory
	} else if p.window.Session.Directory != "" {
		c = p.window.Session.Directory
	}

	rawCmd := splitWindow{
		horizontalSplit: true,
		targetPane:      p.getTargetName(),
		workingDir:      c}
	cmd := exec.Command("tmux", rawCmd.String()...)
	err = cmd.Start()

	return p.window.AddPane(p.Number + 1), err
}

// Split vertically splits the view of the current window to include the current pane
func (p *Pane) Split() (splitPane *Pane, err error) {
	rawCmd := splitWindow{
		verticalSplit: true,
		targetPane:    p.getTargetName()}
	cmd := exec.Command("tmux", rawCmd.String()...)
	err = cmd.Start()

	return p.window.AddPane(p.Number + 1), err
}

// SplitWAttr sets attributes of the current pane to be split horizontally
func (p *Pane) SplitWAttr(attr SplitAttr) (splitPane *Pane, err error) {
	var c string
	if attr.Directory != "" {
		c = attr.Directory
	} else if p.window.Directory != "" {
		c = p.window.Directory
	} else if p.window.Session.Directory != "" {
		c = p.window.Session.Directory
	}

	rawCmd := splitWindow{
		verticalSplit: true,
		targetPane:    p.getTargetName(),
		workingDir:    c}
	cmd := exec.Command("tmux", rawCmd.String()...)
	err = cmd.Start()

	return p.window.AddPane(p.Number + 1), err
}

// ResizeRight resizes the current pane to the right
func (p *Pane) ResizeRight(num int) {
	p.resize("R", num)
}

// ResizeLeft resizes the current pane to the left
func (p *Pane) ResizeLeft(num int) {
	p.resize("L", num)
}

// ResizeUp resizes the current pane upward
func (p *Pane) ResizeUp(num int) {
	p.resize("U", num)
}

// ResizeDown resizes the current pane downward
func (p *Pane) ResizeDown(num int) {
	p.resize("D", num)
}

func (p *Pane) resize(prefix string, num int) {
	rawCmd := []string{"resize-pane", "-t", p.getTargetName(), "-" + prefix, fmt.Sprint(num)}
	cmd := exec.Command("tmux", rawCmd...)
	cmd.Start()
}

// SetName changes the name of a given pane to the specified string
func (p *Pane) SetName(name string) (err error) {
	p.window.Exec(fmt.Sprintf("select-pane -t %d -T %s", p.Number, name))
	return
}

func (p *Pane) getTargetName() (targetName string) {
	return fmt.Sprintf("%s:%s.%s", p.window.Session.Name, fmt.Sprint(p.window.Number), fmt.Sprint(p.Number))
}

func createWindow(number int, attr WindowAttr, session *Session) (window *Window, err error) {
	window = &Window{
		Name:          attr.Name,
		Directory:     attr.Directory,
		Number:        number,
		Session:       session,
		Panes:         make([]*Pane, 0),
		SplitCommands: make([]string, 0),
	}

	window.AddPane(0)

	if number != 0 {
		rawCmd := newWindow{
			targetWindow: window.getTargetWindow(),
			windowName:   window.Name,
			workingDir:   attr.Directory}
		cmd := exec.Command("tmux", rawCmd.String()...)
		err = cmd.Run()
	}

	if err != nil {
		return window, err
	}

	rawCmd := renameWindow{
		targetWindow: window.getTargetWindow(),
		windowName:   window.Name}

	cmd := exec.Command("tmux", rawCmd.String()...)
	err = cmd.Run()
	return window, err
}

func (window *Window) getTargetWindow() (targetWindow []string) {
	args := []string{"-t", fmt.Sprintf("%s:%s", window.Session.Name, fmt.Sprint(window.Number))}
	return args
}

// KillPane removes a pane with the specified pane number in the context of the current window
func (window *Window) KillPane(withNumber int) (err error) {
	window.Exec(fmt.Sprint("kill-pane -t ", withNumber))
	return
}

// AddPane creates a new Pane and adds to the current window
func (window *Window) AddPane(withNumber int) (newPane *Pane) {
	pane := NewPane(withNumber, window)
	window.Panes = append(window.Panes, pane)

	return pane
}

// Pane returns the current Pane object by its index in the Panes slice
func (window *Window) Pane(number int) (currentPane *Pane) {
	return window.Panes[number]
}

// Exec executes a command on the first pane of this window
func (window *Window) Exec(command string) (err error) {
	rawCommand := strings.Split(command, " ")
	tmuxCmd := exec.Command("tmux", rawCommand...)
	return tmuxCmd.Run()
}

// SetConfig aliases window.Exec() to more clearly indicate that window configuration commands are being done
func (window *Window) SetConfig(configCommand string) (err error) {
	return window.Exec(configCommand)
}

// Select changes the active tmux window to the current window object
func (window *Window) Select() (err error) {
	rawCmd := selectWindow{
		targetWindow: window.Session.Name + ":" + fmt.Sprint(window.Number)}
	cmd := exec.Command("tmux", rawCmd.String()...)

	return cmd.Run()
}

// NewSession creates a new tmux session. It will kill any existing session with the provided name.
func NewSession(name string) (session *Session, err error) {
	p := SessionAttr{
		Name: name,
	}

	err = KillSession(p.Name)
	if err != nil {
		return nil, err
	}
	return NewSessionAttr(p)
}

// NewSessionAttr creates a new tmux session based on the provided SessionAttr object. It will kill any existing session with the provided name.
func NewSessionAttr(p SessionAttr) (session *Session, err error) {
	s := &Session{
		Name:      p.Name,
		Directory: p.Directory,
		Windows:   make([]*Window, 0),
	}

	rawCmd := newSession{
		notAttachedToCurrentTerminal: true,
		sessionName:                  p.Name,
		workingDir:                   p.Directory,
		windowName:                   "tmp"}
	cmd := exec.Command("tmux", rawCmd.String()...)
	err = cmd.Run()

	return s, err
}

// CheckSessionExists runs `tmux ls` in order to verify that a given session exists before trying to delete it
func CheckSessionExists(name string) (exists bool, err error) {
	cmd := exec.Command("tmux", "ls", "-F", "'#{session_name}'")
	lsOut, err := cmd.CombinedOutput()
	if err != nil {
		// We may have goofed
		if _, ok := err.(*exec.ExitError); ok {
			// This is fine! tmux will still exit with return code 1, but it's fine.
			if strings.Contains(string(lsOut), "no server running") || strings.Contains(string(lsOut), "no such file or directory") {
				return false, nil
			}
		}
		return false, err
	}

	if strings.Contains(string(lsOut), name) {
		return true, nil
	}
	return false, nil
}

// KillSession sends a command to kill the tmux session
func KillSession(name string) (err error) {

	// See if session exists first
	if exists, err := CheckSessionExists(name); exists == false {
		// If it doesn't exist, check for err
		if err != nil {
			return err
		}

		// It actually doesn't exist!
		return nil
	}

	sess := killSession{
		targetSession: name}

	rawCmd := sess.String()
	cmd := exec.Command("tmux", rawCmd...)

	return cmd.Run()
}

// AddWindow creates a window with provided name for this session
func (s *Session) AddWindow(name string) (window *Window, err error) {
	attr := WindowAttr{
		Name: name,
	}

	return s.AddWindowAttr(attr)
}

// AddWindowAttr creates a window with provided name and WindowAttr properties for this session
func (s *Session) AddWindowAttr(attr WindowAttr) (window *Window, err error) {
	w, err := createWindow(s.NextWindowNumber, attr, s)
	s.Windows = append(s.Windows, w)
	s.NextWindowNumber = s.NextWindowNumber + 1

	return w, err
}
