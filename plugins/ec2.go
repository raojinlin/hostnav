package plugins

import (
	"github.com/raojinlin/jmfzf"
	"github.com/raojinlin/jmfzf/pkg/terminal"
)

type Ec2Plugin struct {
	options *jmfzf.CloudProviderConfig
}

func (e *Ec2Plugin) List(option *jmfzf.ListOptions) ([]terminal.Host, error) {
	// Implement the logic to list EC2 instances
	// Return a slice of Host structs
	return []terminal.Host{}, nil
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
