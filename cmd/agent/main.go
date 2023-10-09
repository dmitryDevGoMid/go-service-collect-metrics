package main

import (
	"fmt"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/app/agent"
)

func main() {
	agent.Run()
	fmt.Println("Hello, World!")
}
