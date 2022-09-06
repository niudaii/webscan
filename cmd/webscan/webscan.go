package main

import (
	"github.com/niudaii/webscan/internal/runner"
	"github.com/projectdiscovery/gologger"
	"time"
)

func main() {
	options := runner.ParseOptions()

	start := time.Now()
	gologger.Info().Msgf("当前时间: %v", start.Format("2006-01-02 15:04:05"))

	newRunner, err := runner.NewRunner(options)
	if err != nil {
		gologger.Fatal().Msgf("Could not create runner: %v", err)
	}
	newRunner.Run()

	gologger.Info().Msgf("运行时间: %v", time.Since(start))
}
