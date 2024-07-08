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
	cli *client.Client
}

func NewDockerPlugin(option interface{}) (jmfzf.Plugin, error) {
	var opt jmfzf.DockerConfig
	if option != nil {
		err := jmfzf.MapToStruct(option, &opt)
		if err != nil {
			return nil, err
		}
	}

	dockerVersion := opt.Version
	if dockerVersion == "" {
		dockerVersion = "1.43"
	}

	cli, err := client.NewClient(opt.Host, opt.Version, nil, nil)
	if err != nil {
		return nil, err
	}

	return &DockerPlugin{cli: cli}, nil
}

func (plugin *DockerPlugin) Name() string {
	return "docker"
}

func (plugin *DockerPlugin) List(options *jmfzf.ListOptions) ([]terminal.Host, error) {
	containers, err := plugin.cli.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		return nil, err
	}

	var result []terminal.Host
	for _, container := range containers {
		result = append(result, terminal.Host{
			Type: terminal.TerminalTypeContainer,
			ContainerInfo: terminal.Container{
				Name:    strings.Trim(container.Names[0], "/"),
				Id:      container.ID,
				Command: "/bin/sh",
			},
		})
	}

	return result, nil
}
