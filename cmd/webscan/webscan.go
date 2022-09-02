package main

import (
	"github.com/niudaii/webscan/internal/runner"
	"github.com/projectdiscovery/gologger"
)

func main() {
	options := runner.ParseOptions()
	newRunner, err := runner.NewRunner(options)
	if err != nil {
		gologger.Fatal().Msgf("Could not create runner: %v", err)
	}
	newRunner.Run()
}
