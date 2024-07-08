package terminal

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Tmux struct{}

func (t *Tmux) NewWindow(title string, command string) error {
	cmd := t.execTmux("new-window", "-n", title, command)
	if cmd == nil {
		return nil
	}
	return cmd.Run()
}

func (t *Tmux) isWindows() bool {
	return strings.Contains(os.Getenv("OS"), "Windows_NT")
}

func (t *Tmux) execTmux(command ...string) *exec.Cmd {
	if t.isWindows() {
		fmt.Println("tmux", command)
		return nil
	}

	cmd := exec.Command("tmux", command...)
	return cmd
}

func (t *Tmux) SplitWindow(command string) error {
	cmd := t.execTmux("split-window", "-h", command)
	if cmd != nil {
		return cmd.Run()
	}

	return nil
}

func (t *Tmux) Version() string {
	cmd := t.execTmux("-V")
	if cmd == nil {
		return ""
	}

	output, _ := cmd.CombinedOutput()
	return strings.Trim(string(output), "\n")
}
