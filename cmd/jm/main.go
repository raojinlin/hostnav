package main

import (
	"fmt"
	"os"

	fzf "github.com/junegunn/fzf/src"
	"github.com/raojinlin/jmfzf/pkg/manager"
)

func main() {
	inputChan := make(chan string)
	pluginManager := manager.New([]string{"ec2", "cvm", "jumpserver"})

	hosts, _ := pluginManager.List(nil)
	go func() {
		for _, host := range hosts {
			name := host.Name + " " + host.PublicIP
			inputChan <- name
		}
		close(inputChan)
	}()

	outputChan := make(chan string)
	go func() {
		for s := range outputChan {
			fmt.Println("Got: " + s)
		}
	}()

	exit := func(code int, err error) {
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		os.Exit(code)
	}

	// Build fzf.Options
	options, err := fzf.ParseOptions(
		true, // whether to load defaults ($FZF_DEFAULT_OPTS_FILE and $FZF_DEFAULT_OPTS)
		[]string{"--multi", "--reverse", "--border"},
	)
	if err != nil {
		exit(fzf.ExitError, err)
	}

	// Set up input and output channels
	options.Input = inputChan
	options.Output = outputChan

	// Run fzf
	code, err := fzf.Run(options)
	exit(code, err)
}
