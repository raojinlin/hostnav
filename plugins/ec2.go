package plugins

import (
	"github.com/raojinlin/jmfzf"
	"github.com/raojinlin/jmfzf/pkg/terminal"
)

type Ec2Plugin struct {
	option *jmfzf.CloudProviderOption
}

func (p *Ec2Plugin) List(options *ListOptions) ([]terminal.Host, error) {
	// Implement the logic to list EC2 instances
	// Return a slice of Host structs
	return []terminal.Host{}, nil
}

func (e *Ec2Plugin) Name() string {
	return "ec2"
}

func (e *Ec2Plugin) Cache() bool {
	return true
}

func (e *Ec2Plugin) Init(option interface{}) error {
	return nil
}

func NewEc2Plugin() *Ec2Plugin {
	return &Ec2Plugin{}
}
