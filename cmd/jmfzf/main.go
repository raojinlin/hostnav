package main

import (
	"fmt"
	"log/slog"
	"os"

	fzf "github.com/junegunn/fzf/src"
	"github.com/raojinlin/jmfzf"
	"github.com/raojinlin/jmfzf/pkg/manager"
	"github.com/raojinlin/jmfzf/pkg/terminal"
)

func main() {
	cfg, err := jmfzf.NewConfig("./config.yaml")
	if err != nil {
		panic(err)
	}

	inputChan := make(chan string)
	hostManager := manager.New([]string{"docker", "bce", "cvm", "jumpserver"}, cfg)

	hosts, _ := hostManager.List(nil)
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
	done := make(chan struct{}) // 用来通知任务完成

	go func() {
		defer close(outputChan)
		output := <-outputChan
		host := indexedHosts[output]
		err := host.Connect()
		if err != nil {
			slog.Error("Error connecting to host", "error", err)
		}
		done <- struct{}{}
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

	code, err := fzf.Run(options)
	if err != nil {
		exit(code, err)
	}

	// 等待输出协程完成
	<-done
}
