package plugins

import (
	"fmt"

	"github.com/baidubce/bce-sdk-go/services/bcc"
	"github.com/baidubce/bce-sdk-go/services/bcc/api"
	"github.com/raojinlin/hostnav"
	"github.com/raojinlin/hostnav/pkg/terminal"
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
	option        *hostnav.CloudProviderOption
	regionClients map[string]*bcc.Client
}

func NewBcePlugin() *BcePlugin {
	return &BcePlugin{option: &hostnav.CloudProviderOption{}}
}

func (p *BcePlugin) Init(option interface{}) error {
	if err := hostnav.MapToStruct(option, p.option); err != nil {
		return err
	}

	regionClients := make(map[string]*bcc.Client)
	for _, region := range p.option.Regions {
		endpoint, ok := bceRegionEndpoints[region]
		if !ok || endpoint == "" {
			return fmt.Errorf("invalid region: %s", region)
		}

		bce, err := bcc.NewClient(p.option.AccessKey, p.option.AccessKeySecret, endpoint)
		if err != nil {
			return fmt.Errorf("create bce %s client: %v", region, err)
		}

		regionClients[region] = bce
	}

	p.regionClients = regionClients
	return nil
}

func (p *BcePlugin) List(options *ListOptions) ([]terminal.Host, error) {
	var instances []api.InstanceModel
	for _, client := range p.regionClients {
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
					Name:       fmt.Sprintf("%s(%s): %s", p.Name(), instance.ZoneName, instance.InstanceName),
					PublicIP:   instance.PublicIP,
					LocalIP:    instance.InternalIP,
					User:       "root",
					Port:       22,
					UseLocalIP: p.option.ConnectionOptions.UseLocalIP,
				},
			},
		)
	}

	return result, nil
}

func (p *BcePlugin) Name() string {
	return "bce"
}

func (p *BcePlugin) Cache() bool {
	return true
}
