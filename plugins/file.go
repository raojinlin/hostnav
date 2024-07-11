package plugins

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/raojinlin/hostnav"
	"github.com/raojinlin/hostnav/pkg/terminal"
	"gopkg.in/yaml.v3"
)

// FilePlugin represents a plugin that load hosts from a given file(json|yaml|text)
type FilePlugin struct {
	Option *hostnav.FileOption
}

func NewFilePlugin() *FilePlugin {
	return &FilePlugin{Option: &hostnav.FileOption{}}
}

func (p *FilePlugin) Init(option interface{}) error {
	if err := hostnav.MapToStruct(option, p.Option); err != nil {
		return err
	}

	return nil
}

func (p *FilePlugin) Cache() bool { return false }

func (p *FilePlugin) Name() string {
	return "file"
}

func (p *FilePlugin) List(options *ListOptions) ([]terminal.Host, error) {
	var result []terminal.Host
	for _, file := range *p.Option {
		f, err := os.Open(file)
		if err != nil {
			return nil, err
		}

		data, err := io.ReadAll(f)
		if err != nil {
			return nil, err
		}

		parts := strings.Split(file, ".")
		suffix := parts[len(parts)-1]

		var hosts []terminal.Host
		switch suffix {
		case "json":
			err = json.Unmarshal(data, &hosts)
			break
		case "yaml", "yml":
			fmt.Println(suffix, hosts)
			err = yaml.Unmarshal(data, &hosts)
		default:
			err = fmt.Errorf("unknown file type")
		}

		if err != nil {
			return nil, err
		}

		result = append(result, hosts...)
	}

	return result, nil
}
