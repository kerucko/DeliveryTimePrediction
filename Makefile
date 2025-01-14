BINARY_NAME=bg

all:

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
