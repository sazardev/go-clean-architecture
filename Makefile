# HR API Makefile para Windows PowerShell

.PHONY: help build run test clean deps docker-build docker-run

# Variables
APP_NAME = hr-api
DOCKER_IMAGE = hr-api:latest
PORT = 8080

help: ## Mostrar esta ayuda
	@echo "Comandos disponibles:"
	@echo "  build        - Compilar la aplicación"
	@echo "  run          - Ejecutar la aplicación"
	@echo "  test         - Ejecutar tests"
	@echo "  clean        - Limpiar archivos compilados"
	@echo "  deps         - Descargar dependencias"
	@echo "  docker-build - Construir imagen Docker"
	@echo "  docker-run   - Ejecutar contenedor Docker"

build: ## Compilar la aplicación
	go build -o $(APP_NAME).exe cmd/server/main.go

run: ## Ejecutar la aplicación
	go run cmd/server/main.go

test: ## Ejecutar tests
	go test -v ./...

test-coverage: ## Ejecutar tests con coverage
	go test -v -cover ./...

clean: ## Limpiar archivos compilados
	if exist $(APP_NAME).exe del $(APP_NAME).exe

deps: ## Descargar dependencias
	go mod download
	go mod tidy

docker-build: ## Construir imagen Docker
	docker build -t $(DOCKER_IMAGE) .

docker-run: ## Ejecutar contenedor Docker
	docker run -p $(PORT):$(PORT) --env-file .env $(DOCKER_IMAGE)

dev: ## Ejecutar en modo desarrollo con hot reload
	go run cmd/server/main.go
