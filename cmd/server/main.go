// Package server provides a command-line interface for the metrics server.
package main

import (
	"fmt"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/app/server"
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
	server.Run()
	//Fixed: call to os.Exit found at /Users/aleksandrserbakov/cloneservicecollect/go-service-collect-metrics/cmd/server/main.go:23:2
	//os.Exit(1)
}
