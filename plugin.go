package jmfzf

type Host struct {
	Name         string `json:"name" yaml:"name"`
	PublicIP     string `json:"public_ip" yaml:"public_ip"`
	LocalIP      string `json:"local_ip" yaml:"local_ip"`
	Port         int    `json:"port" yaml:"port"`
	User         string `json:"user" yaml:"user"`
	IdentityFile string `json:"identity_file" yaml:"identity_file"`
	Enabled      bool   `json:"enabled" yaml:"enabled"`
}

// ListOptions represents a list of options for a plugin
type ListOptions struct {
	Order   int    `json:"order" yaml:"order"`
	OrderBy string `json:"order_by" yaml:"order_by"`
}

type Plugin interface {
	List(option *ListOptions) ([]Host, error)
	Name() string
}
