package plugins

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/raojinlin/jmfzf"
	"github.com/raojinlin/jmfzf/pkg/terminal"
)

type DockerPlugin struct {
	option *jmfzf.DockerOption
	cli    *client.Client
}

func NewDockerPlugin() *DockerPlugin {
	return &DockerPlugin{option: &jmfzf.DockerOption{}}
}

func (p *DockerPlugin) Init(option interface{}) error {
	if err := jmfzf.MapToStruct(option, p.option); err != nil {
		return err
	}


	dockerVersion := p.option.Version
	if dockerVersion == "" {
		dockerVersion = "1.43"
	}

	cli, err := client.NewClient(p.option.Host, p.option.Version, nil, nil)
	if err != nil {
		return err
	}

	p.cli = cli
	return nil
}

func (p *DockerPlugin) Name() string {
	return "docker"
}

func (p *DockerPlugin) List(options *ListOptions) ([]terminal.Host, error) {
	containers, err := p.cli.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		return nil, err
	}

	var result []terminal.Host
	for _, container := range containers {
		result = append(result, terminal.Host{
			Type: terminal.TerminalTypeContainer,
			ContainerInfo: terminal.Pod{
				Name:      "docker",
				Namespace: terminal.NamespaceDocker,
				Container: terminal.Container{
					Id:      container.ID,
					Name:    strings.Trim(container.Names[0], "/"),
					Command: "/bin/sh",
				},
			},
		})
	}

	return result, nil
}

func (p *DockerPlugin) Cache() bool {
	return false
}
