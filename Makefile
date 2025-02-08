start-forever-program:
	go run cmd/forever/main.go 100

debug-forever-program:
	LOG_ENABLED=1 go run cmd/forever/main.go 100

start-monitor:
	go run cmd/monitor/main.go forever
