package plugins

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/raojinlin/hostnav"
	"github.com/raojinlin/hostnav/pkg/terminal"
)

type Ec2Plugin struct {
	Option        *hostnav.CloudProviderOption
	regionClients map[string]*ec2.EC2
}

func (p *Ec2Plugin) List(options *ListOptions) ([]terminal.Host, error) {
	// Implement the logic to list EC2 instances
	// Return a slice of Host structs
	var result []terminal.Host
	for _, ec2Client := range p.regionClients {
		output, err := ec2Client.DescribeInstances(nil)
		if err != nil {
			return nil, err
		}

		for _, reservation := range output.Reservations {
			for _, instance := range reservation.Instances {
				result = append(result, terminal.Host{
					Type: terminal.TerminalTypeHost,
					SSHInfo: terminal.SSHInfo{
						Name:       *instance.PublicDnsName,
						PublicIP:   *instance.PublicIpAddress,
						LocalIP:    *instance.PrivateIpAddress,
						User:       "root",
						Port:       22,
						UseLocalIP: p.Option.ConnectionOptions.UseLocalIP,
					},
				})
			}
		}
	}

	return result, nil
}

func (e *Ec2Plugin) Name() string {
	return "ec2"
}

func (e *Ec2Plugin) Cache() bool {
	return true
}

func (e *Ec2Plugin) Init(option interface{}) error {
	if err := hostnav.MapToStruct(option, e.Option); err != nil {
		return err
	}

	os.Setenv("AWS_ACCESS_KEY_ID", e.Option.AccessKey)
	os.Setenv("AWS_SECRET_ACCESS_KEY", e.Option.AccessKeySecret)

	for _, region := range e.Option.Regions {
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String(region),
		})

		if err != nil {
			return err
		}

		e.regionClients[region] = ec2.New(sess)
	}

	return nil
}

func NewEc2Plugin() *Ec2Plugin {
	return &Ec2Plugin{Option: &hostnav.CloudProviderOption{}, regionClients: make(map[string]*ec2.EC2)}
}
