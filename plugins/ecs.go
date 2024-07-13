package plugins

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	ecs20140526 "github.com/alibabacloud-go/ecs-20140526/v4/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/raojinlin/hostnav"
	"github.com/raojinlin/hostnav/pkg/terminal"
)

var regionMap = map[string]string{
	"cn-qingdao":       "ecs.cn-qingdao.aliyuncs.com",
	"cn-beijing":       "ecs.cn-beijing.aliyuncs.com",
	"cn-zhangjiakou":   "ecs.cn-zhangjiakou.aliyuncs.com",
	"cn-huhehaote":     "ecs.cn-huhehaote.aliyuncs.com",
	"cn-wulanchabu":    "ecs.cn-wulanchabu.aliyuncs.com",
	"cn-hangzhou":      "ecs-cn-hangzhou.aliyuncs.com",
	"cn-shanghai":      "ecs.cn-shanghai.aliyuncs.com",
	"cn-nanjing":       "ecs.cn-nanjing.aliyuncs.com",
	"cn-fuzhou":        "ecs.cn-fuzhou.aliyuncs.com",
	"cn-shenzhen":      "ecs.cn-shenzhen.aliyuncs.com",
	"cn-heyuan":        "ecs.cn-heyuan.aliyuncs.com",
	"cn-guangzhou":     "ecs.cn-guangzhou.aliyuncs.com",
	"cn-wuhan-lr":      "ecs.cn-wuhan-lr.aliyuncs.com",
	"ap-southeast-2":   "ecs.ap-southeast-2.aliyuncs.com",
	"ap-southeast-6":   "ecs.ap-southeast-6.aliyuncs.com",
	"ap-northeast-2":   "ecs.ap-northeast-2.aliyuncs.com",
	"ap-southeast-3":   "ecs.ap-southeast-3.aliyuncs.com",
	"ap-northeast-1":   "ecs.ap-northeast-1.aliyuncs.com",
	"ap-southeast-7":   "ecs.ap-southeast-7.aliyuncs.com",
	"cn-chengdu":       "ecs.cn-chengdu.aliyuncs.com",
	"ap-southeast-1":   "ecs.ap-southeast-1.aliyuncs.com",
	"ap-southeast-5":   "ecs.ap-southeast-5.aliyuncs.com",
	"cn-zhengzhou-jva": "ecs.cn-zhengzhou-jva.aliyuncs.com",
	"cn-hongkong":      "ecs.cn-hongkong.aliyuncs.com",
	"eu-central-1":     "ecs.eu-central-1.aliyuncs.com",
	"us-east-1":        "ecs.us-east-1.aliyuncs.com",
	"us-west-1":        "ecs.us-west-1.aliyuncs.com",
	"eu-west-1":        "ecs.eu-west-1.aliyuncs.com",
	"me-east-1":        "ecs.me-east-1.aliyuncs.com",
	"me-central-1":     "ecs.me-central-1.aliyuncs.com",
	"ap-south-1":       "ecs.ap-south-1.aliyuncs.com",
}

func createEcsClient(endpoint string, accessKey, accessKeySecret string) (client *ecs20140526.Client, err error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(accessKey),
		AccessKeySecret: tea.String(accessKeySecret),
	}
	config.Endpoint = tea.String(endpoint)
	client = &ecs20140526.Client{}
	client, err = ecs20140526.NewClient(config)
	return client, err
}

type EcsPlugin struct {
	Option        hostnav.CloudProviderOption
	regionsClient map[string]*ecs20140526.Client
}

func NewEcsPlugin() *EcsPlugin {
	return &EcsPlugin{regionsClient: make(map[string]*ecs20140526.Client)}
}

func (p *EcsPlugin) Init(option interface{}) error {
	if err := hostnav.MapToStruct(option, &p.Option); err != nil {
		return err
	}

	for _, region := range p.Option.Regions {
		client, err := createEcsClient(regionMap[region], p.Option.AccessKey, p.Option.AccessKeySecret)
		if err != nil {
			return err
		}

		p.regionsClient[region] = client
	}

	return nil
}

func (p *EcsPlugin) Name() string {
	return "ecs"
}

func (p *EcsPlugin) Cache() bool {
	return true
}

func (p *EcsPlugin) list(client *ecs20140526.Client, region string) ([]terminal.Host, error) {
	var result []terminal.Host
	var page int32 = 1
	var pageSize int32 = 100
	var nextToken string
	var firstPage bool = true

	for nextToken != "" || firstPage {
		response, err := client.DescribeInstances(&ecs20140526.DescribeInstancesRequest{
			PageNumber: tea.Int32(page),
			PageSize:   tea.Int32(pageSize),
			NextToken:  &nextToken,
			RegionId:   tea.String(region),
		})

		if err != nil {
			return nil, err
		}

		firstPage = false
		nextToken = *response.Body.NextToken

		for _, instance := range response.Body.Instances.Instance {
			result = append(result, terminal.Host{
				Type: terminal.TerminalTypeContainer,
				SSHInfo: terminal.SSHInfo{
					Name:     *instance.InstanceName,
					PublicIP: *instance.EipAddress.IpAddress,
					LocalIP:  *instance.PublicIpAddress.IpAddress[0],
					User:     "root",
					Port:     22,
				},
			})
		}
	}

	return result, nil

}

func (p *EcsPlugin) List(options *ListOptions) ([]terminal.Host, error) {
	var result []terminal.Host
	for region, client := range p.regionsClient {
		regionInstances, err := p.list(client, region)
		if err != nil {
			return nil, err
		}

		result = append(result, regionInstances...)
	}

	return result, nil
}
