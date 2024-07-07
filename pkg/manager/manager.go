package manager

import (
	"log/slog"

	"github.com/raojinlin/jmfzf"
	"github.com/raojinlin/jmfzf/plugins"
)

var pluginConstructor = map[string]func(option interface{}) (jmfzf.Plugin, error){
	"ec2":        plugins.NewEc2Plugin,
	"cvm":        plugins.NewCVMPlugin,
	"jumpserver": plugins.NewJumpServerPlugin,
}

type Manager struct {
	plugins []jmfzf.Plugin
}

func New(pluginNames []string) *Manager {
	plugins := make([]jmfzf.Plugin, 0)
	for _, pluginName := range pluginNames {
		if constructor, ok := pluginConstructor[pluginName]; ok {
			plugin, err := constructor(nil)
			if err != nil {
				slog.Warn("error creating", "plugin", pluginName, "error", err.Error())
				continue
			}

			plugins = append(plugins, plugin)
		} else {
			slog.Warn("unknown plugin", "pluginName", pluginName)
		}
	}
	return &Manager{plugins: plugins}
}

func (m *Manager) List(options *jmfzf.ListOptions) ([]jmfzf.Host, error) {
	result := make([]jmfzf.Host, 0)
	for _, plugin := range m.plugins {
		hosts, err := plugin.List(options)
		if err == nil {
			result = append(result, hosts...)
		}
	}

	return result, nil
}
