package terminal

import (
	"fmt"
	"strings"
)

type Terminal interface {
	Connect() error
}

type TerminalType = string

var tmux = &Tmux{}

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

func (s *SSHInfo) Connect() error {
	return tmux.NewWindow(s.Name, strings.Join([]string{"ssh", fmt.Sprintf("%s@%s", s.User, s.PublicIP), "-p", fmt.Sprintf("%d", s.Port)}, " "))
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
		return h.SSHInfo.Connect()
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
	command := fmt.Sprintf("kubectl --config %s exec -n %s -it %s -- /bin/sh", p.KubeConfig, p.Namespace, p.Name)
	return tmux.NewWindow("Pod: "+p.Name, command)
}

type Container struct {
	Id      string `json:"id" yaml:"id"`
	Name    string `json:"name" yaml:"name"`
	Command string `json:"command" yaml:"command"`
	Pod     *Pod   `json:"pod" yaml:"pod"`
}

func (c *Container) String() string {
	if c.Pod != nil {
		return "pod:" + c.Pod.String() + " " + c.Name
	}

	return fmt.Sprintf("docker: %s", c.Name)
}

func (c *Container) Connect() error {
	if c.Pod != nil {
		return c.Pod.Connect()
	}

	return tmux.NewWindow("container: "+c.Name, fmt.Sprintf("docker exec -it %s %s", c.Name, c.Command))
}
