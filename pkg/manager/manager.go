package manager

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"log"
	"log/slog"
	"time"

	"github.com/raojinlin/hostnav"
	"github.com/raojinlin/hostnav/pkg/cache"
	"github.com/raojinlin/hostnav/pkg/terminal"
	"github.com/raojinlin/hostnav/plugins"
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
	config  *hostnav.Config
}

func New(pluginNames []string, config *hostnav.Config) *Manager {
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
	return &Manager{plugins: plugins, config: config}
}

func (m *Manager) List(options *plugins.ListOptions) ([]terminal.Host, error) {
	result := make([]terminal.Host, 0)
	cacheDuration := time.Duration(m.config.Cache.Duration) * time.Minute
	if cacheDuration == 0 {
		cacheDuration = time.Minute * 30
	}
	filecache := cache.NewFileCache[[]terminal.Host](m.config.Cache.Directory, cacheDuration)

	err := filecache.Load()
	if err != nil {
		slog.Warn("load cache", "error", err.Error())
	}

	defer filecache.Save()
	for _, plugin := range m.plugins {
		pluginOption, _ := json.Marshal(plugin)
		pluginMd5 := md5.New().Sum(pluginOption)
		cacheKey := plugin.Name() + hex.EncodeToString(pluginMd5)

		log.Println("list", "plugin", plugin.Name(), "hosts...")
		var hosts []terminal.Host
		var err error
		if plugin.Cache() {
			cacheContent, err := filecache.Get(cacheKey)
			if err == nil {
				hosts = *cacheContent
			}
		}

		if len(hosts) == 0 {
			hosts, err = plugin.List(options)
			if plugin.Cache() {
				filecache.Set(cacheKey, hosts, 0)
			}
		}
		if err == nil {
			result = append(result, hosts...)
		} else {
			slog.Warn("list", "plugin", plugin.Name(), "error", err.Error())
		}
	}

	return result, nil
}
