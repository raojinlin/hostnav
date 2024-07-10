package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path"
	"strings"

	fzf "github.com/junegunn/fzf/src"
	"github.com/raojinlin/hostnav"
	"github.com/raojinlin/hostnav/pkg/manager"
	"github.com/raojinlin/hostnav/pkg/terminal"
)

const quitSignalCode int = 130

func exit(code int, err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	os.Exit(code)
}

func main() {
	homedir, _ := os.UserHomeDir()
	var configfile string = path.Join(homedir, ".hostnav.yaml")
	// plugins to use, mutiple plugins coma separated
	var plugins string
	flag.StringVar(&configfile, "config", configfile, "path to configuration file")
	flag.StringVar(&plugins, "plugins", plugins, "plugin to use, mutiple plugins comma separated")

	flag.Parse()
	cfg, err := hostnav.NewConfig(configfile)
	if err != nil {
		exit(1, err)
	}

	inputChan := make(chan string)
	enabledPlugins := []string{}
	if len(plugins) == 0 {
		enabledPlugins = cfg.DefaultPlugins
	} else {
		for _, plugin := range strings.Split(plugins, ",") {
			if strings.Trim(plugin, ", ") != "" {
				enabledPlugins = append(enabledPlugins, plugin)
			}
		}
	}

	hostManager := manager.New(enabledPlugins, cfg)

	hosts, _ := hostManager.List(nil)
	if len(hosts) == 0 {
		slog.Warn("no host information")
		return
	}
	indexedHosts := make(map[string]terminal.Host)
	go func() {
		for _, host := range hosts {
			name := host.String()
			indexedHosts[name] = host
			inputChan <- name
		}
		close(inputChan)
	}()

	outputChan := make(chan string)
	done := make(chan struct{})

	go func() {
		defer close(outputChan)
		output := <-outputChan
		host := indexedHosts[output]
		slog.Info("connect to", "host", host.String())
		err := host.Connect()
		if err != nil {
			slog.Error("Error connecting to host", "error", err)
		}
		done <- struct{}{}
	}()

	// Build fzf.Options
	options, err := fzf.ParseOptions(
		true, // whether to load defaults ($FZF_DEFAULT_OPTS_FILE and $FZF_DEFAULT_OPTS)
		[]string{"--reverse"},
	)
	if err != nil {
		exit(fzf.ExitError, err)
	}

	// Set up input and output channels
	options.Input = inputChan
	options.Output = outputChan

	go func() {

		code, err := fzf.Run(options)
		if err != nil {
			exit(code, err)
		}

		// quit code
		if code == quitSignalCode {
			done <- struct{}{}
		}
	}()

	<-done
}
