package plugins

import "github.com/raojinlin/jmfzf"

type Ec2PluginOptions struct {
	AccessKey       string `json:"access_key" yaml:"access_key"`
	AccessKeySecret string `json:"access_key_secret" yaml:"access_key_secret"`
	Zone            string `json:"zone" yaml:"zone"`
}

type Ec2Plugin struct {
	options *Ec2PluginOptions
}

func (e *Ec2Plugin) List(option *jmfzf.ListOptions) ([]jmfzf.Host, error) {
	// Implement the logic to list EC2 instances
	// Return a slice of Host structs
	return []jmfzf.Host{
		{
			Name:         "test",
			PublicIP:     "223.5.5.5",
			Port:         22,
			User:         "root",
			IdentityFile: "~/.ssh/id_rsa",
			LocalIP:      "127.0.0.1",
		},
		{
			Name:         "ec2-test2",
			PublicIP:     "223.5.5.5",
			Port:         22,
			User:         "root",
			IdentityFile: "~/.ssh/id_rsa",
			LocalIP:      "127.0.0.1",
		},
	}, nil
}

func (e *Ec2Plugin) Name() string {
	return "ec2"
}

func NewEc2Plugin(options interface{}) (jmfzf.Plugin, error) {
	var opt *Ec2PluginOptions
	if options != nil {
		opt = options.(*Ec2PluginOptions)
	}
	return &Ec2Plugin{options: opt}, nil
}
