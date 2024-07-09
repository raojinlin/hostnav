package manager

import (
	"log"
	"log/slog"

	"github.com/raojinlin/jmfzf"
	"github.com/raojinlin/jmfzf/pkg/terminal"
	"github.com/raojinlin/jmfzf/plugins"
)

var supportPlugins = map[string]plugins.Plugin{
	"ec2":        plugins.NewEc2Plugin(),
	"cvm":        plugins.NewCVMPlugin(),
	"jumpserver": plugins.NewJumpServerPlugin(),
	"bce":        plugins.NewBcePlugin(),
	"docker":     plugins.NewDockerPlugin(),
	"kubernetes": plugins.NewKubernetesPlugin(),
}

type Manager struct {
	plugins []plugins.Plugin
}

func New(pluginNames []string, config *jmfzf.Config) *Manager {
	plugins := make([]plugins.Plugin, 0)
	for _, pluginName := range pluginNames {
		if plugin, ok := supportPlugins[pluginName]; ok {
			if err := plugin.Init(config.Plugins[pluginName]); err != nil {
				slog.Warn("error init", "plugin", pluginName, "error", err.Error())
				continue
			}

			plugins = append(plugins, plugin)
		} else {
			slog.Warn("unknown plugin", "pluginName", pluginName)
		}
	}
	return &Manager{plugins: plugins}
}

func (m *Manager) List(options *plugins.ListOptions) ([]terminal.Host, error) {
	result := make([]terminal.Host, 0)
	for _, plugin := range m.plugins {
		log.Println("list", "plugin", plugin.Name(), "hosts...")
		hosts, err := plugin.List(options)
		if err == nil {
			result = append(result, hosts...)
		} else {
			slog.Warn("list", "plugin", plugin.Name(), "error", err.Error())
		}
	}

	return result, nil
}
