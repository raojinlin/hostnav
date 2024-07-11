package plugins

import (
	"fmt"

	"github.com/baidubce/bce-sdk-go/services/bcc"
	"github.com/baidubce/bce-sdk-go/services/bcc/api"
	"github.com/raojinlin/hostnav"
	"github.com/raojinlin/hostnav/pkg/terminal"
)

var bccRegionEndpoints = map[string]string{
	"cn-bj":  "bcc.bj.baidubce.com",
	"cn-bd":  "bcc.bd.baidubce.com",
	"cn-gz":  "bcc.gz.baidubce.com",
	"cn-su":  "bcc.su.baidubce.com",
	"cn-hkg": "bcc.hkg.baidubce.com",
	"cn-fwh": "bcc.fwh.baidubce.com",
	"cn-cd":  "bcc.cd.baidubce.com",
	"cn-fsh": "bcc.fsh.baidubce.com",
}

type BccPlugin struct {
	Option        *hostnav.CloudProviderOption
	regionClients map[string]*bcc.Client
}

func NewBccPlugin() *BccPlugin {
	return &BccPlugin{Option: &hostnav.CloudProviderOption{}}
}

func (p *BccPlugin) Init(option interface{}) error {
	if err := hostnav.MapToStruct(option, p.Option); err != nil {
		return err
	}

	regionClients := make(map[string]*bcc.Client)
	for _, region := range p.Option.Regions {
		endpoint, ok := bccRegionEndpoints[region]
		if !ok || endpoint == "" {
			return fmt.Errorf("invalid region: %s", region)
		}

		bcc, err := bcc.NewClient(p.Option.AccessKey, p.Option.AccessKeySecret, endpoint)
		if err != nil {
			return fmt.Errorf("create bcc %s client: %v", region, err)
		}

		regionClients[region] = bcc
	}

	p.regionClients = regionClients
	return nil
}

func (p *BccPlugin) List(options *ListOptions) ([]terminal.Host, error) {
	var instances []api.InstanceModel
	for _, client := range p.regionClients {
		resp, err := client.ListInstances(&api.ListInstanceArgs{})
		if err != nil {
			return nil, fmt.Errorf("list bcc instances: %v", err)
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
					UseLocalIP: p.Option.ConnectionOptions.UseLocalIP,
				},
			},
		)
	}

	return result, nil
}

func (p *BccPlugin) Name() string {
	return "bcc"
}

func (p *BccPlugin) Cache() bool {
	return true
}
