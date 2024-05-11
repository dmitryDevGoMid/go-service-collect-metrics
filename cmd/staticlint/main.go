// Package staticlint provides a command-line interface for running static analysis checks on Go code.
// Overview
//
// The staticlint package provides a command-line interface for running a set of static analysis checks on Go code.
// The package uses the multichecker tool to run multiple analyzers in parallel and aggregate the results.
//
// # Usage
//
// To run the staticlint analyzer, use the following command:
//
//	go run cmd/staticlint/main.go ./...
//
// This command runs the staticlint analyzer on all Go packages in the current directory and its subdirectories.
// The analyzer reads the configuration from the file cmd/staticlint/statickcheck.toml and prints the results to the console.
//
// # Configuration
//
// The staticlint analyzer uses a configuration file in TOML format to specify which checks to enable.
// The configuration file is located at cmd/staticlint/statickcheck.toml and has the following format:
//
//	checks = ["SA","ST1000","ST1000","ST1005","ST1013"]
//
// The configuration file specifies a list of checks to enable using the Staticcheck tool.
//
// # Analyzers
//
// The staticlint analyzer includes the following analyzers:
//
//   - Staticcheck: a set of checks for style, correctness, and best practices.
//   - Osexit: a check for incorrect usage of os.Exit.
//   - Printf: a check for incorrect usage of fmt.Printf and friends.
//   - Shadow: a check for shadowed variables.
//   - Shift: a check for inefficient shifts.
//   - Structtag: a check for incorrect struct tags.
//
// # Results
//
// The staticlint analyzer prints the results of each check to the console.
// The results include a summary of the number of errors and warnings found, as well as detailed information about each issue.
// If any errors are found, the analyzer returns a non-zero exit code.
package main

import (
	"os"
	"path"
	"strings"

	"github.com/BurntSushi/toml"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/cmd/staticlint/osexit"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/structtag"

	"honnef.co/go/tools/staticcheck"
)

// Config represents the configuration for the static analysis checks.
type Config struct {
	// Staticcheck is a list of checks to enable.
	Staticcheck []string `toml:"checks"`
}

// main is the entry point for the command-line interface. It reads the configuration,
// initializes the analyzers, and runs them using the multichecker.
func main() {
	// Get the current working directory.
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// Read the configuration file.
	var cfg Config
	if _, err := toml.DecodeFile(path.Join(cwd, "/cmd/staticlint/statickcheck.toml"), &cfg); err != nil {
		panic(err)
	}

	// Create a list of analyzers.
	var analyzers []*analysis.Analyzer

	// Add the staticcheck analyzers specified in the configuration.
	for _, check := range cfg.Staticcheck {
		for _, a := range staticcheck.Analyzers {
			if strings.HasPrefix(a.Analyzer.Name, check) || a.Analyzer.Name == check {
				analyzers = append(analyzers, a.Analyzer)
			}
		}
	}

	// Add additional analyzers.
	analyzers = append(
		analyzers,

		// osexit analyzer checks for incorrect usage of os.Exit.
		osexit.Analyzer,

		// printf analyzer checks for incorrect usage of fmt.Printf and friends.
		printf.Analyzer,

		// shadow analyzer checks for shadowed variables.
		shadow.Analyzer,

		// shift analyzer checks for inefficient shifts.
		shift.Analyzer,

		// structtag analyzer checks for incorrect struct tags.
		structtag.Analyzer)

	// Run the multichecker with the specified analyzers.
	// The multichecker runs each analyzer in parallel and aggregates the results.
	// It then prints the results to the console and returns a non-zero exit code if any errors were found.
	multichecker.Main(analyzers...)
}
