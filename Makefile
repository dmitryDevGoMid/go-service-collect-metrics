.PHONY: server
server:
	go build -o cmd/server/server cmd/server/main.go
.PHONY: servrun
servrun:
	go run cmd/server/main.go
.PHONY: agent
agent:
	go build -o cmd/agent/server cmd/agent/main.go
.PHONY: agentrun
agentrun:
	go run cmd/agent/main.go