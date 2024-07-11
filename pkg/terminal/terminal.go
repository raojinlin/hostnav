package terminal

import (
	"fmt"
)

const (
	NamespaceDocker = "__docker__"
)

type ConnectionOption struct {
	Tmux struct {
		NewWindow bool
	} `json:"tmux"`
}

// Terminal represents a terminal that can be used to interact
type Terminal interface {
	Connect(option *ConnectionOption) error
}

type TerminalType = string

var tmux = &Tmux{}

const (
	// TerminalTypeHost represents a host that can be used to interact with ssh commands
	TerminalTypeHost TerminalType = "host"
	// TerminalTypeContainer represents a container that can be used to interact with docker/kubectl commands
	TerminalTypeContainer TerminalType = "container"
)

type SSHInfo struct {
	Name         string `json:"name" yaml:"name"`
	PublicIP     string `json:"public_ip" yaml:"public_ip"`
	LocalIP      string `json:"local_ip" yaml:"local_ip"`
	Port         int    `json:"port" yaml:"port"`
	User         string `json:"user" yaml:"user"`
	IdentityFile string `json:"identity_file" yaml:"identity_file"`
	UseLocalIP   bool   `json:"use_local_ip" yaml:"use_local_ip"`
}

func (s *SSHInfo) Connect(option *ConnectionOption) error {
	opts := "-o StrictHostKeyChecking=no"
	if s.IdentityFile != "" {
		opts += " -i " + s.IdentityFile
	}

	dest := s.PublicIP
	if s.UseLocalIP && s.LocalIP != "" {
		dest = s.LocalIP
	}

	user := s.User
	if user == "" {
		user = "root"
	}

	port := s.Port
	if port == 0 {
		port = 22
	}

	command := fmt.Sprintf("ssh %s %s@%s -p%d", opts, user, dest, port)
	if option.Tmux.NewWindow {
		return tmux.NewWindow(s.Name, command)
	}

	return tmux.SplitWindow(command)
}

func (s *SSHInfo) String() string {
	return fmt.Sprintf("%s %s %s", s.Name, s.PublicIP, s.LocalIP)
}

type Host struct {
	Type          TerminalType `json:"type" yaml:"type"`
	SSHInfo       SSHInfo      `json:"ssh_info" yaml:"ssh_info"`
	ContainerInfo Pod          `json:"container_info" yaml:"container_info"`
}

func (h *Host) String() string {
	if h.Type == TerminalTypeContainer {
		return h.ContainerInfo.String()
	}

	return h.SSHInfo.String()
}

func (h *Host) Connect(option *ConnectionOption) error {
	if h.Type == TerminalTypeHost {
		return h.SSHInfo.Connect(option)
	}

	if h.Type == TerminalTypeContainer {
		return h.ContainerInfo.Connect(option)
	}

	return fmt.Errorf("unkown terminal type: %s", h.Type)
}

type Pod struct {
	Name       string    `json:"name" yaml:"name"`
	KubeConfig string    `json:"kube_config" yaml:"kube_config"`
	Namespace  string    `json:"namespace" yaml:"namespace"`
	Container  Container `json:"container" yaml:"container"`
}

func (p *Pod) String() string {
	if p.Namespace == NamespaceDocker {
		return fmt.Sprintf("Container: %s", p.Container.Name)
	}

	return fmt.Sprintf("POD: %s/%s/%s", p.Namespace, p.Name, p.Container.Name)
}

func (p *Pod) Connect(option *ConnectionOption) error {
	var command string
	var windowName string
	if p.Namespace == NamespaceDocker {
		command = fmt.Sprintf("docker exec -it %s %s", p.Container.Name, p.Container.Command)
		windowName = "Container: " + p.Container.Name
	} else {
		command = fmt.Sprintf("kubectl --kubeconfig %s exec -n %s -it %s -c %s -- %s", p.KubeConfig, p.Namespace, p.Name, p.Container.Name, p.Container.Command)
		windowName = fmt.Sprintf("Pod: %s/%s", p.Name, p.Container.Name)
	}

	if option.Tmux.NewWindow {
		return tmux.NewWindow(windowName, command)
	}

	return tmux.SplitWindow(command)
}

type Container struct {
	Id      string `json:"id" yaml:"id"`
	Name    string `json:"name" yaml:"name"`
	Command string `json:"command" yaml:"command"`
}

func (c *Container) String() string {
	return fmt.Sprintf("container: %s", c.Name)
}
