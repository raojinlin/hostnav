package hostnav

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Tag struct {
	Name  string `json:"name" yaml:"name"`
	Value string `json:"value" yaml:"value"`
}

type CloudProviderOption struct {
	Endpoint          string   `json:"endpoint" yaml:"endpoint"`
	AccessKey         string   `json:"access_key" yaml:"access_key"`
	AccessKeySecret   string   `json:"access_key_secret" yaml:"access_key_secret"`
	Zones             []string `json:"zones" yaml:"zones"`
	Regions           []string `json:"regions" yaml:"regions"`
	Tags              []Tag    `json:"tags" yaml:"tags"`
	ConnectionOptions struct {
		UseLocalIP bool `json:"use_local_ip" yaml:"use_local_ip"`
	} `json:"connections" yaml:"connections"`
}

type JumpServerOption struct {
	AccessKey       string `json:"access_key" yaml:"access_key"`
	AccessKeySecret string `json:"access_key_secret" yaml:"access_key_secret"`
	Url             string `json:"url" yaml:"url"`
	ApiToken        string `json:"api_token" yaml:"api_token"`
	Search          string `json:"search" yaml:"search"`
}

type SshConfig []string

type DockerOption struct {
	Host    string `json:"host" yaml:"host"`
	Version string `json:"version" yaml:"version"`
}

type KubernetesOption struct {
	KubeConfig string   `json:"kubeconfig" yaml:"kubeconfig"`
	Namespaces []string `json:"namespaces" yaml:"namespaces"`
}

type Config struct {
	Cache struct {
		Directory string `json:"directory" yaml:"directory"`
		Duration  int    `json:"duration" yaml:"duration"`
	} `json:"cache" yaml:"cache"`
	DefaultPlugins []string               `json:"default_plugins" yaml:"default_plugins"`
	Plugins        map[string]interface{} `json:"plugins" yaml:"plugins"`
}

func NewConfig(configfile string) (*Config, error) {
	if configfile == "" {
		return nil, fmt.Errorf("configfile cannot be empty")
	}

	data, err := os.ReadFile(configfile)
	if err != nil {
		return nil, fmt.Errorf("read config file: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("unmarshal config file: %v", err)
	}

	return &config, nil
}
