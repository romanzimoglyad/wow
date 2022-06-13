run:
	docker-compose up --force-recreate --build server --build client
start-server:
	go run cmd/server/main.go
start-client:
	go run cmd/client/main.go
stop:
	docker-compose down
test:
	go test ./...