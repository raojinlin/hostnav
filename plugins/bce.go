package plugins

import (
	"fmt"

	"github.com/baidubce/bce-sdk-go/services/bcc"
	"github.com/baidubce/bce-sdk-go/services/bcc/api"
	"github.com/raojinlin/jmfzf"
	"github.com/raojinlin/jmfzf/pkg/terminal"
)

var bceRegionEndpoints = map[string]string{
	"cn-bj":  "bcc.bj.baidubce.com",
	"cn-bd":  "bcc.bd.baidubce.com",
	"cn-gz":  "bcc.gz.baidubce.com",
	"cn-su":  "bcc.su.baidubce.com",
	"cn-hkg": "bcc.hkg.baidubce.com",
	"cn-fwh": "bcc.fwh.baidubce.com",
	"cn-cd":  "bcc.cd.baidubce.com",
	"cn-fsh": "bcc.fsh.baidubce.com",
}

type BcePlugin struct {
	options       *jmfzf.CloudProviderConfig
	regionClients map[string]*bcc.Client
}

func NewBcePlugin(options interface{}) (jmfzf.Plugin, error) {
	var opt jmfzf.CloudProviderConfig
	if options != nil {
		if err := jmfzf.MapToStruct(options, &opt); err != nil {
			return nil, err
		}
	}

	regionClients := make(map[string]*bcc.Client)
	for _, region := range opt.Regions {
		endpoint, ok := bceRegionEndpoints[region]
		if !ok || endpoint == "" {
			return nil, fmt.Errorf("invalid region: %s", region)
		}

		bce, err := bcc.NewClient(opt.AccessKey, opt.AccessKeySecret, endpoint)
		if err != nil {
			return nil, fmt.Errorf("create bce %s client: %v", region, err)
		}

		regionClients[region] = bce
	}

	return &BcePlugin{options: &opt, regionClients: regionClients}, nil
}

func (plugin *BcePlugin) List(options *jmfzf.ListOptions) ([]terminal.Host, error) {
	var instances []api.InstanceModel
	for _, client := range plugin.regionClients {
		resp, err := client.ListInstances(&api.ListInstanceArgs{})
		if err != nil {
			return nil, fmt.Errorf("list bce instances: %v", err)
		}

		instances = append(instances, resp.Instances...)
	}

	var result []terminal.Host
	for _, instance := range instances {
		result = append(
			result,
			terminal.Host{
				Type: terminal.TerminalTypeHost,
				SSHInfo: terminal.SSHInfo{
					Name:     fmt.Sprintf("%s(%s): %s", plugin.Name(), instance.ZoneName, instance.InstanceName),
					PublicIP: instance.PublicIP,
					LocalIP:  instance.InternalIP,
					User:     "root",
					Port:     22,
				},
			},
		)
	}

	return result, nil
}

func (plugin *BcePlugin) Name() string {
	return "bce"
}
