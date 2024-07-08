package plugins

import (
	"fmt"

	"github.com/raojinlin/jmfzf"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

type CVMPlugin struct {
	options       *jmfzf.CloudProviderConfig
	regionClients map[string]*cvm.Client
}

func NewCVMPlugin(options interface{}) (jmfzf.Plugin, error) {
	var opt jmfzf.CloudProviderConfig
	var err error
	if options != nil {
		err = jmfzf.MapToStruct(options, &opt)
		if err != nil {
			return nil, err
		}
	}

	cert := common.NewCredential(opt.AccessKey, opt.AccessKeySecret)
	profile := profile.NewClientProfile()
	regionClients := make(map[string]*cvm.Client)
	for _, region := range opt.Regions {
		client, err := cvm.NewClient(cert, region, profile)
		if err != nil {
			return nil, fmt.Errorf("failed to create %s client: %v", region, err)
		}

		regionClients[region] = client
	}

	return &CVMPlugin{options: &opt, regionClients: regionClients}, nil
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

func (p *CVMPlugin) List(options *jmfzf.ListOptions) ([]jmfzf.Host, error) {
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

	var result []jmfzf.Host
	for _, instance := range instances {
		result = append(result, jmfzf.Host{
			Name:     fmt.Sprintf("%s(%s): %s", p.Name(), *instance.Placement.Zone, *instance.InstanceName),
			PublicIP: *instance.PublicIpAddresses[0],
			Port:     22,
			User:     "root",
			LocalIP:  *instance.PrivateIpAddresses[0],
		})
	}

	return result, nil
}
