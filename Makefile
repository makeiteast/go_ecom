compose-up: ### Run docker-compose
	docker-compose up --build -d postgres redis swagger-ui swagger-editor

compose-down: ### Down docker-compose
	docker-compose down --remove-orphans

run: ### Run docker-compose
	go run main.go	

test: ### run test
	go test -v -cover -race ./...