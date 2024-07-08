package plugins

import (
	"fmt"

	"github.com/raojinlin/jmfzf"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

type CVMPlugin struct {
	options *jmfzf.CloudProviderConfig
	cvm     *cvm.Client
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
	client, err := cvm.NewClient(cert, "ap-shanghai", profile.NewClientProfile())
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}

	return &CVMPlugin{options: &opt, cvm: client}, nil
}

func (p *CVMPlugin) Name() string {
	return "cvm"
}

func (p *CVMPlugin) List(options *jmfzf.ListOptions) ([]jmfzf.Host, error) {
	// Implement the logic to list CVM instances
	// Return a slice of Host structs
	req := cvm.NewDescribeInstancesRequest()
	req.Limit = common.Int64Ptr(100)
	resp, err := p.cvm.DescribeInstances(req)
	if err != nil {
		return nil, fmt.Errorf("failed to describe instances: %v", err)
	}
	var result []jmfzf.Host
	fmt.Println(*resp.Response.TotalCount)
	for _, instance := range resp.Response.InstanceSet {
		result = append(result, jmfzf.Host{
			Name:     p.Name() + ": " + *instance.InstanceName,
			PublicIP: *instance.PublicIpAddresses[0],
			Port:     22,
			User:     "root",
			LocalIP:  *instance.PrivateIpAddresses[0],
		})
	}

	return result, nil
}
