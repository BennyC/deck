integration:
	go test ./... -run "Integration"
unit:
	go test ./... -run "Unit"
test:
	go test ./...
run:
	go run cmd/api/main.go
build:
	go build cmd/api/main.go