package plugins

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/raojinlin/jmfzf"
	"github.com/raojinlin/jmfzf/pkg/terminal"
	"gopkg.in/twindagger/httpsig.v1"
)

// MetaData represents the metadata structure in the JSON
type MetaData struct {
	Type string `json:"type"`
	Data struct {
		ID        string   `json:"id"`
		Hostname  string   `json:"hostname"`
		IP        string   `json:"ip"`
		Protocols []string `json:"protocols"`
		Platform  string   `json:"platform"`
		OrgName   string   `json:"org_name"`
	} `json:"data"`
}

// Asset represents the main structure of the JSON
type Asset struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Title       string   `json:"title"`
	PId         string   `json:"pId"`
	IsParent    bool     `json:"isParent"`
	Open        bool     `json:"open"`
	IconSkin    string   `json:"iconSkin"`
	ChkDisabled bool     `json:"chkDisabled"`
	Meta        MetaData `json:"meta"`
}

type SigAuth struct {
	KeyID    string
	SecretID string
}

func (auth *SigAuth) Sign(r *http.Request) error {
	headers := []string{"(request-target)", "date"}
	signer, err := httpsig.NewRequestSigner(auth.KeyID, auth.SecretID, "hmac-sha256")
	if err != nil {
		return err
	}
	return signer.SignRequest(r, headers, nil)
}

type JumpServerPlugin struct {
	option *jmfzf.JumpServerOption
}

func NewJumpServerPlugin() *JumpServerPlugin {
	return &JumpServerPlugin{option: &jmfzf.JumpServerOption{}}
}

func (p *JumpServerPlugin) Init(option interface{}) error {
	return jmfzf.MapToStruct(option, p.option)
}

func (p *JumpServerPlugin) doRequest(method string, url string, headers map[string]string, body io.Reader) (*http.Response, error) {
	gmtFmt := "Mon, 02 Jan 2006 15:04:05 GMT"
	req, err := http.NewRequest(method, p.option.Url+url, body)
	req.Header.Add("Date", time.Now().Format(gmtFmt))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-JMS-ORG", "00000000-0000-0000-0000-000000000002")

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	if err != nil {
		return nil, err
	}

	auth := SigAuth{KeyID: p.option.AccessKey, SecretID: p.option.AccessKeySecret}
	if err := auth.Sign(req); err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (p *JumpServerPlugin) getUserPermsAssets() ([]Asset, error) {
	resp, err := p.doRequest("GET", "/api/v1/perms/users/assets/tree/", nil, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var assets []Asset
	err = json.Unmarshal(data, &assets)
	if err != nil {
		return nil, err
	}

	return assets, nil
}

func (p *JumpServerPlugin) List(options *ListOptions) ([]terminal.Host, error) {
	assets, err := p.getUserPermsAssets()
	if err != nil {
		return nil, err
	}

	var hosts []terminal.Host
	for _, asset := range assets {
		hosts = append(hosts, terminal.Host{
			Type: terminal.TerminalTypeHost,
			SSHInfo: terminal.SSHInfo{
				PublicIP: asset.Meta.Data.IP,
				Name:     p.Name() + ": " + asset.Name,
				Port:     22,
				User:     "root",
			},
		})
	}

	return hosts, nil
}

func (p *JumpServerPlugin) Name() string {
	return "jumpserver"
}

func (p *JumpServerPlugin) Cache() bool {
	return true
}
