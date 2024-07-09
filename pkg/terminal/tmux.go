package terminal

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mattn/go-shellwords"
)

type Tmux struct{}

func (t *Tmux) NewWindow(title string, command string) error {
	// fall back to the default window
	if !t.inTmux() || !t.hasTmux() {
		return t.exec(command)
	}

	// use tmux if available
	cmd := t.execTmux("new-window", "-n", title, command)
	if cmd == nil {
		return nil
	}
	return cmd.Run()
}

func (t *Tmux) isWindows() bool {
	return os.Getenv("OS") == "Windows_NT"
}

func (t *Tmux) inTmux() bool {
	return os.Getenv("TMUX") != ""
}

func (t *Tmux) hasTmux() bool {
	return !t.isWindows() && t.Version() != ""
}

func (t *Tmux) exec(command string) error {
	var cmd *exec.Cmd
	parser := shellwords.NewParser()
	words, err := parser.Parse(command)
	if err != nil {
		return fmt.Errorf("error parsing %s: %v", command, err)
	}

	cmd = exec.Command(words[0], words[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (t *Tmux) execTmux(command ...string) *exec.Cmd {
	cmd := exec.Command("tmux", command...)
	return cmd
}

func (t *Tmux) SplitWindow(command string) error {
	if t.inTmux() && t.hasTmux() {
		cmd := t.execTmux("split-window", "-h", command)
		if cmd != nil {
			return cmd.Run()
		}
		return nil
	}

	return t.exec(command)

}

func (t *Tmux) Version() string {
	cmd := t.execTmux("-V")
	if cmd == nil {
		return ""
	}

	output, _ := cmd.CombinedOutput()
	return strings.Trim(string(output), "\n")
}
