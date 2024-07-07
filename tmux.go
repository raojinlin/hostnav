package jmfzf

import (
	"fmt"
	"os/exec"
	"time"
)

type Tmux struct{}

func (t *Tmux) NewWindow(title string, command string) error {
	cmd := exec.Command("tmux", "new-window", "-n", title, command)
	return cmd.Run()
}

func (t *Tmux) SplitWindow(command string) error {
	cmd := exec.Command("tmux", "split-window", "-h", command)
	return cmd.Run()
}

func (t *Tmux) Version() string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in Version", r)
		}
	}()
	fmt.Println("tmux.Version called")
	time.Sleep(500 * time.Microsecond)
	fmt.Println("tmux.Version completed")
	return "xxx"
	// cmd := exec.Command("tmux", "-V")
	// output, _ := cmd.CombinedOutput()
	// return strings.Trim(string(output), "\n")
}
