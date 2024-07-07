package plugins

import (
	"github.com/raojinlin/jmfzf"
)

type JumpServerOptions struct{}

type JumpServerPlugin struct {
	options *JumpServerOptions
}

func NewJumpServerPlugin(options interface{}) (jmfzf.Plugin, error) {
	// if options == nil {
	// 	return nil, fmt.Errorf("options must not be nil")
	// }

	// opt, ok := options.(*JumpServerOptions)
	// if !ok {
	// 	return nil, fmt.Errorf("invalid options type")
	// }

	var opt *JumpServerOptions
	if options != nil {
		opt = options.(*JumpServerOptions) // type assertion
	}

	return &JumpServerPlugin{options: opt}, nil
}

func (plugin *JumpServerPlugin) List(options *jmfzf.ListOptions) ([]jmfzf.Host, error) {
	return []jmfzf.Host{
		{
			Name:         "jumpserver-test1",
			PublicIP:     "192.168.3.2",
			LocalIP:      "172.16.0.1",
			User:         "root",
			IdentityFile: "~/.ssh/id_rsa",
		},
	}, nil
}

func (plugin *JumpServerPlugin) Name() string {
	return "jumpserver"
}
