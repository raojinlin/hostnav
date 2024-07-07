package plugins

import "github.com/raojinlin/jmfzf"

type CVMPluginOptions struct{}

type CVMPlugin struct {
	options *CVMPluginOptions
}

func NewCVMPlugin(option interface{}) (jmfzf.Plugin, error) {
	var opt *CVMPluginOptions
	if option != nil {
		opt = option.(*CVMPluginOptions) // type assertion
	}
	return &CVMPlugin{options: opt}, nil
}

func (p *CVMPlugin) Name() string {
	return "cvm"
}

func (p *CVMPlugin) List(options *jmfzf.ListOptions) ([]jmfzf.Host, error) {
	// Implement the logic to list CVM instances
	// Return a slice of Host structs
	return []jmfzf.Host{
		{
			Name:         "cvm-test",
			PublicIP:     "192.168.1.1",
			Port:         22,
			User:         "root",
			IdentityFile: "~/.ssh/id_rsa",
			LocalIP:      "127.0.0.1",
		},
		{
			Name:         "cvm-test2",
			PublicIP:     "192.168.1.2",
			Port:         22,
			User:         "root",
			IdentityFile: "~/.ssh/id_rsa",
			LocalIP:      "127.0.0.1",
		},
	}, nil
}
