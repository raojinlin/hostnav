package plugins

import (
	"fmt"

	"github.com/raojinlin/hostnav"
	"github.com/raojinlin/hostnav/pkg/terminal"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

type CVMPlugin struct {
	Option        *hostnav.CloudProviderOption
	regionClients map[string]*cvm.Client
}

func NewCVMPlugin() *CVMPlugin {
	return &CVMPlugin{Option: &hostnav.CloudProviderOption{}}
}

func (p *CVMPlugin) Init(option interface{}) error {
	err := hostnav.MapToStruct(option, p.Option)
	if err != nil {
		return err
	}

	cert := common.NewCredential(p.Option.AccessKey, p.Option.AccessKeySecret)
	profile := profile.NewClientProfile()
	regionClients := make(map[string]*cvm.Client)
	for _, region := range p.Option.Regions {
		client, err := cvm.NewClient(cert, region, profile)
		if err != nil {
			return fmt.Errorf("failed to create %s client: %v", region, err)
		}

		regionClients[region] = client
	}

	p.regionClients = regionClients
	return nil
}

func (p *CVMPlugin) Name() string {
	return "cvm"
}

func (p *CVMPlugin) listInstances(client *cvm.Client) ([]*cvm.Instance, error) {
	var instances []*cvm.Instance
	var offset int64
	var limit int64 = 100
	var total int64 = limit

	for offset <= total {
		req := cvm.NewDescribeInstancesRequest()
		req.Limit = common.Int64Ptr(limit)
		req.Offset = common.Int64Ptr(offset)
		resp, err := client.DescribeInstances(req)
		if err != nil {
			return nil, fmt.Errorf("failed to describe instances: %v", err)
		}
		instances = append(instances, resp.Response.InstanceSet...)
		total = *resp.Response.TotalCount
		offset += limit
	}

	return instances, nil
}

func (p *CVMPlugin) List(options *ListOptions) ([]terminal.Host, error) {
	// Implement the logic to list CVM instances
	// Return a slice of Host structs
	var instances []*cvm.Instance

	for _, client := range p.regionClients {
		result, err := p.listInstances(client)
		if err != nil {
			return nil, fmt.Errorf("failed to list instances: %v", err)
		}

		instances = append(instances, result...)
	}

	var result []terminal.Host
	for _, instance := range instances {
		result = append(result, terminal.Host{
			Type: terminal.TerminalTypeHost,
			SSHInfo: terminal.SSHInfo{
				Name:       fmt.Sprintf("%s(%s): %s", p.Name(), *instance.Placement.Zone, *instance.InstanceName),
				PublicIP:   *instance.PublicIpAddresses[0],
				Port:       22,
				User:       "root",
				LocalIP:    *instance.PrivateIpAddresses[0],
				UseLocalIP: p.Option.ConnectionOptions.UseLocalIP,
			},
		})
	}

	return result, nil
}

func (p *CVMPlugin) Cache() bool {
	return true
}
