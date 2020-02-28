gomux
=====

Go wrapper to create tmux sessions, windows and panes. This is a minimal fork of `wricardo/gomux` to meet our specific needs.

### Example
example.go:
```go
package main
import (
	"os"
	"github.com/disneystreaming/gomux"
)
func main() {
	sessionName := "SESSION_NAME"

	s, _ := gomux.NewSession(sessionName)

	//WINDOW 1
	w1, _ := s.AddWindow("Monitoring")

	w1p0 := w1.Pane(0)
	w1p0.Exec("htop")

	w1p1, _ := w1.Pane(0).Split()
	w1p1.Exec("tail -f /var/log/syslog")

	//WINDOW 2
	w2 := s.AddWindow("Vim")
	w2p0 := w2.Pane(0)

	w2p0.Exec("echo \"this is to vim\" | vim -")

	w2p1 := w2p0.Vsplit()
	w2p1.Exec("cd /tmp/")
	w2p1.Exec("ls -la")

	w2p0.ResizeRight(30)
	w1.Select()
}
```

To create and attach to the tmux session:
```
go run example.go
tmux attach -t SESSION_NAME
```
