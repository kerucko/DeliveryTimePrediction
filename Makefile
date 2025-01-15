BINARY_NAME=bg

build:
	go build -o $(BINARY_NAME) ./backend-gateway/cmd

clean:
	go clean
	rm -f $(BINARY_NAME)

run-gateway: build
	./$(BINARY_NAME)


docker-up: docker-compose.yml
	docker compose up -d

docker-down: docker-compose.yml
	docker compose down


deps:
	go install github.com/pressly/goose/v3/cmd/goose@latest


migrate:
	goose -dir ./migrations postgres "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" up

down:
	goose -dir ./migrations postgres "user=postgres dbname=postgres password=postgres host=localhost:5432" down

reset:
	goose -dir ./migrations postgres "user=postgres dbname=postgres password=postgres host=localhost:5432" reset