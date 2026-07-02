# ==========================================
# Go Eventos - Makefile
# ==========================================

APP_NAME := go-eventos
MAIN_PATH := cmd/api/main.go
BIN_DIR := bin
BIN_PATH := $(BIN_DIR)/$(APP_NAME)

# Carrega variáveis do .env, se existir (usado pelos targets de migration)
ifneq (,$(wildcard .env))
	include .env
	export
endif

MIGRATIONS_PATH := migrations

.PHONY: help run build test test-cover tidy fmt vet lint clean \
        swag docker-build docker-run \
        db-up db-down db-logs \
        migrate-up migrate-down migrate-force migrate-version migrate-create \
        db-reset

## help: lista todos os comandos disponíveis
help:
	@echo "Comandos disponiveis:"
	@grep -E '^## ' Makefile | sed 's/## /  /'

## run: executa a aplicacao localmente
run:
	go run $(MAIN_PATH)

## build: compila o binario para bin/
build:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_PATH) $(MAIN_PATH)
	@echo "Binario gerado em $(BIN_PATH)"

## test: executa os testes
test:
	go test ./... -v

## test-cover: executa os testes com relatorio de cobertura
test-cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Relatorio gerado em coverage.html"

## tidy: organiza e baixa as dependencias do go.mod
tidy:
	go mod tidy

## fmt: formata o codigo
fmt:
	go fmt ./...

## vet: analisa o codigo em busca de erros comuns
vet:
	go vet ./...

## lint: roda fmt e vet
lint: fmt vet

## clean: remove binarios e arquivos gerados
clean:
	rm -rf $(BIN_DIR) coverage.out coverage.html logs.json

## swag: gera a documentacao Swagger em docs/
swag:
	swag init -g $(MAIN_PATH) -o docs

## docker-build: builda a imagem docker da aplicacao
docker-build:
	docker build -t $(APP_NAME) .

## docker-run: sobe a aplicacao via docker
docker-run:
	docker run --rm -p 8080:8080 --env-file .env $(APP_NAME)

## db-up: sobe o container do PostgreSQL via docker-compose
db-up:
	docker compose up -d db

## db-down: derruba o container do PostgreSQL
db-down:
	docker compose down

## db-logs: mostra os logs do container do PostgreSQL
db-logs:
	docker compose logs -f db

## migrate-up: aplica todas as migrations pendentes
migrate-up:
	migrate -database "$(DATABASE_URL)" -path $(MIGRATIONS_PATH) up

## migrate-down: reverte a ultima migration aplicada
migrate-down:
	migrate -database "$(DATABASE_URL)" -path $(MIGRATIONS_PATH) down 1

## migrate-version: mostra a versao atual das migrations
migrate-version:
	migrate -database "$(DATABASE_URL)" -path $(MIGRATIONS_PATH) version

## migrate-force: forca a versao das migrations (uso: make migrate-force VERSION=1)
migrate-force:
	migrate -database "$(DATABASE_URL)" -path $(MIGRATIONS_PATH) force $(VERSION)

## migrate-create: cria um novo par de arquivos de migration (uso: make migrate-create NAME=nome_da_migration)
migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(NAME)

## db-reset: reverte todas as migrations e aplica novamente
db-reset:
	migrate -database "$(DATABASE_URL)" -path $(MIGRATIONS_PATH) down -all
	migrate -database "$(DATABASE_URL)" -path $(MIGRATIONS_PATH) up
