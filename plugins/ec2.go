package plugins

import (
	"github.com/raojinlin/jmfzf"
)

type Ec2Plugin struct {
	options *jmfzf.CloudProviderConfig
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
	var opt jmfzf.CloudProviderConfig
	if options != nil {
		err := jmfzf.MapToStruct(options, &opt)
		if err != nil {
			return nil, err
		}
	}
	return &Ec2Plugin{options: &opt}, nil
}
