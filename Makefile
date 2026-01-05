.PHONY: install build run clean

install:
	@echo "Installing dependencies..."
	go mod download

build: install
	@echo "Building go server..."
	go build -o main.exe main.go

start: install build
	@echo "Starting go server..."
	chmod +x release.sh
	./release.sh

clean:
	@echo "Cleaning up..."
	rm -f main.exe
	docker compose down -v
	rm -rf pgdata mongodata