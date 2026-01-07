.PHONY: install build run start clean

install:
	@echo "Installing dependencies..."
	go mod download

build: install
	@echo "Building go server..."
	go build -o main main.go

run: build
	@echo "Running go server..."
	./main

start: install build
	@echo "Starting go server..."
	chmod +x release.sh
	./release.sh

clean:
	@echo "Cleaning up..."
	rm -f main
	docker compose down -v
	rm -rf pgdata mongodata