package cache

import (
	"testing"
	"time"

	"github.com/raojinlin/hostnav/pkg/terminal"
)

func TestCache(t *testing.T) {
	filecache := NewFileCache[terminal.Host]("./cache", time.Minute*5)
	filecache.Load()
	filecache.Set("xxxx", terminal.Host{Type: terminal.TerminalTypeContainer, ContainerInfo: terminal.Pod{Name: "test1"}}, time.Minute*5)
	filecache.Set("xxxx1", terminal.Host{Type: terminal.TerminalTypeContainer, ContainerInfo: terminal.Pod{Name: "test2"}}, time.Minute*5)
	filecache.Set("xxxx3", terminal.Host{Type: terminal.TerminalTypeContainer, ContainerInfo: terminal.Pod{Name: "test3"}}, time.Minute*5)
	filecache.Set("xxxx4", terminal.Host{Type: terminal.TerminalTypeContainer, ContainerInfo: terminal.Pod{Name: "test4"}}, time.Minute*5)
	if err := filecache.Save(); err != nil {
		t.Fatal(err)
	}
}

func TestReadCache(t *testing.T) {
	filecache := NewFileCache[terminal.Host]("./cache", time.Minute*5)
	err := filecache.Load()
	if err != nil {
		t.Fatal(err)
	}
	value, err := filecache.Get("xxxx")
	if err != nil {
		t.Fatal(err)
	}
	filecache.Get("xxxx")
	filecache.Get("xxxx")
	filecache.Get("xxxx")

	if value.Type != terminal.TerminalTypeContainer && value.ContainerInfo.Name != "test" {
		t.Fatal("except for container info")
	}
}
