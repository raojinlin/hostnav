package main

import (
	"fmt"
	"os"

	fzf "github.com/junegunn/fzf/src"
	"github.com/raojinlin/jmfzf"
	"github.com/raojinlin/jmfzf/pkg/manager"
)

func main() {
	cfg, err := jmfzf.NewConfig("./config.yaml")
	if err != nil {
		panic(err)
	}

	inputChan := make(chan string)
	pluginManager := manager.New([]string{"cvm", "bce", "jumpserver"}, cfg)

	hosts, _ := pluginManager.List(nil)
	go func() {
		for _, host := range hosts {
			name := host.Name + " " + host.PublicIP
			inputChan <- name
		}
		close(inputChan)
	}()

	outputChan := make(chan string)
	done := make(chan struct{}) // 用来通知任务完成

	go func() {
		defer close(outputChan)
		output := <-outputChan
		fmt.Println("got output: ", output)
		fmt.Println("ddd")
		tmux := &jmfzf.Tmux{}
		tmux.SplitWindow("docker exec -it hanxiucao_mysql /bin/bash")
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

	fmt.Println("Main program done")
}
