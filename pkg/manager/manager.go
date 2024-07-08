package manager

import (
	"log"
	"log/slog"

	"github.com/raojinlin/jmfzf"
	"github.com/raojinlin/jmfzf/plugins"
)

var pluginConstructor = map[string]func(option interface{}) (jmfzf.Plugin, error){
	"ec2":        plugins.NewEc2Plugin,
	"cvm":        plugins.NewCVMPlugin,
	"jumpserver": plugins.NewJumpServerPlugin,
	"bce":        plugins.NewBcePlugin,
}

type Manager struct {
	plugins []jmfzf.Plugin
}

func New(pluginNames []string, config *jmfzf.Config) *Manager {
	plugins := make([]jmfzf.Plugin, 0)
	for _, pluginName := range pluginNames {
		if constructor, ok := pluginConstructor[pluginName]; ok {
			pluginConfig := config.Plugins[pluginName]
			plugin, err := constructor(pluginConfig)
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
		log.Println("fetching", "plugin", plugin.Name(), "hosts")
		hosts, err := plugin.List(options)
		if err == nil {
			result = append(result, hosts...)
		} else {
			slog.Warn("fetcing", "plugin", plugin.Name(), "error", err)
		}
	}

	return result, nil
}
