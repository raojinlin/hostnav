package jmfzf

import "github.com/raojinlin/jmfzf/pkg/terminal"

// ListOptions represents a list of options for a plugin
type ListOptions struct {
	Order   int    `json:"order" yaml:"order"`
	OrderBy string `json:"order_by" yaml:"order_by"`
}

type Plugin interface {
	List(option *ListOptions) ([]terminal.Host, error)
	Name() string
}
