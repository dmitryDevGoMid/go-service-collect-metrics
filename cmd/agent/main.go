// Package agent provides a command-line interface for the metrics agent.
package main

import (
	"fmt"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/app/agent"
)

// Build variables
var buildVersion string
var buildDate string
var buildCommit string

func init() {
	// Print build information
	fmt.Println("Build version:", buildVersion)
	fmt.Println("Build date:", buildDate)
	fmt.Println("Build commit:", buildCommit)
}

func main() {
	agent.Run()
	//Fixed: call to os.Exit found at /Users/aleksandrserbakov/cloneservicecollect/go-service-collect-metrics/cmd/agent/main.go:23:2
	//os.Exit(1)
}
