package terminal

import (
	"fmt"
	"os"
	"os/exec"
)

type Terminal interface {
	Connect() error
}

type TerminalType = string

const (
	TerminalTypeHost      TerminalType = "host"
	TerminalTypeContainer TerminalType = "container"
)

type SSHInfo struct {
	Name         string `json:"name" yaml:"name"`
	PublicIP     string `json:"public_ip" yaml:"public_ip"`
	LocalIP      string `json:"local_ip" yaml:"local_ip"`
	Port         int    `json:"port" yaml:"port"`
	User         string `json:"user" yaml:"user"`
	IdentityFile string `json:"identity_file" yaml:"identity_file"`
}

type Host struct {
	Type          TerminalType `json:"type" yaml:"type"`
	SSHInfo       SSHInfo      `json:"ssh_info" yaml:"ssh_info"`
	ContainerInfo Container    `json:"container_info" yaml:"container_info"`
}

func (h *Host) String() string {
	if h.Type == TerminalTypeContainer {
		return h.ContainerInfo.String()
	}

	return fmt.Sprintf("%s %s", h.SSHInfo.Name, h.SSHInfo.PublicIP)
}

func (h *Host) Connect() error {
	if h.Type == TerminalTypeHost {
		cmd := exec.Command("ssh", fmt.Sprintf("%s@%s", h.SSHInfo.User, h.SSHInfo.PublicIP), "-p", fmt.Sprintf("%d", h.SSHInfo.Port))
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		return cmd.Run()
	}

	if h.Type == TerminalTypeContainer {
		return h.ContainerInfo.Connect()
	}

	return fmt.Errorf("unkown terminal type: %s", h.Type)
}

type Pod struct {
	Name       string `json:"name" yaml:"name"`
	KubeConfig string `json:"kube_config" yaml:"kube_config"`
	Namespace  string `json:"namespace" yaml:"namespace"`
}

func (p *Pod) String() string {
	return fmt.Sprintf("%s/%s", p.Namespace, p.Name)
}

func (p *Pod) Connect() error {
	cmd := exec.Command("kubectl", "--config", p.KubeConfig, "exec", "-n", p.Namespace, "-it", p.Name, "--", "/bin/sh")
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

type Container struct {
	Id      string `json:"id" yaml:"id"`
	Name    string `json:"name" yaml:"name"`
	Command string `json:"command" yaml:"command"`
	Pod     *Pod   `json:"pod" yaml:"pod"`
}

func (c *Container) String() string {
	if c.Pod != nil {
		return fmt.Sprintf("%s %s", c.Pod.Name, c.Name)
	}

	return fmt.Sprintf("conatiner: %s", c.Name)
}

func (c *Container) Connect() error {
	if c.Pod != nil {
		return c.Pod.Connect()
	}

	cmd := exec.Command("docker", "exec", "-it", c.Name, c.Command)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
