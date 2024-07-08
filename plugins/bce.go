package plugins

import (
	"fmt"

	"github.com/baidubce/bce-sdk-go/services/bcc"
	"github.com/baidubce/bce-sdk-go/services/bcc/api"
	"github.com/raojinlin/jmfzf"
)

type BcePlugin struct {
	options *jmfzf.CloudProviderConfig
	bce     *bcc.Client
}

func NewBcePlugin(options interface{}) (jmfzf.Plugin, error) {
	var opt jmfzf.CloudProviderConfig
	if options != nil {
		if err := jmfzf.MapToStruct(options, &opt); err != nil {
			return nil, err
		}
	}

	bce, err := bcc.NewClient(opt.AccessKey, opt.AccessKeySecret, opt.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("create bce client: %v", err)
	}
	return &BcePlugin{options: &opt, bce: bce}, nil
}

func (plugin *BcePlugin) List(options *jmfzf.ListOptions) ([]jmfzf.Host, error) {
	resp, err := plugin.bce.ListInstances(&api.ListInstanceArgs{})
	if err != nil {
		return nil, fmt.Errorf("list bce instances: %v", err)
	}

	var result []jmfzf.Host
	for _, instance := range resp.Instances {
		result = append(
			result,
			jmfzf.Host{
				Name:     plugin.Name() + ": " + instance.InstanceName,
				PublicIP: instance.PublicIP,
				LocalIP:  instance.InternalIP,
				User:     "root",
				Port:     22,
			},
		)
	}

	return result, nil
}

func (plugin *BcePlugin) Name() string {
	return "bce"
}
