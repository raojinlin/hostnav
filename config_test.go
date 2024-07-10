package hostnav 

import (
	"fmt"
	"testing"
)

func TestReadConfig(t *testing.T) {
	cfg, err := NewConfig("./config.yaml")
	if err != nil {
		t.Fatal(err)
	}

	if len(cfg.Plugins) == 0 {
		t.Fatal("No plugins")
	}

	fmt.Println(cfg.DefaultPlugins)
	for k, p := range cfg.Plugins {
		fmt.Println(k, p)
	}
}
