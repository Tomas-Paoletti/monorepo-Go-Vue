# Atajos del monorepo. Uso: make <target>
# (En Windows con PowerShell, si no tenés `make`, corré los comandos a mano.)

.PHONY: help migrate import api front-dev front-build test build tidy spec

help:
	@echo "Targets:"
	@echo "  migrate      - corre las migraciones (GORM AutoMigrate)"
	@echo "  import       - importa el dataset MovieLens a la DB"
	@echo "  api          - levanta la API en :8080"
	@echo "  front-dev    - Vite dev server en :5173 (proxy a :8080)"
	@echo "  front-build  - build de produccion del frontend"
	@echo "  test         - corre los tests del backend"
	@echo "  build        - compila el binario del backend"
	@echo "  tidy         - go mod tidy en el backend"
	@echo "  spec         - specs: make spec ARGS=\"new checkout --variant=bdd\""

migrate:
	cd backend && go run ./cmd/migrate

import:
	cd backend && go run ./cmd/import

api:
	cd backend && go run ./cmd/api

front-dev:
	cd frontend && npm run dev

front-build:
	cd frontend && npm run build

test:
	cd backend && go test ./...

build:
	cd backend && go build -o bin/api ./cmd/api

tidy:
	cd backend && go mod tidy

spec:
	cd ai-tools && npm run spec -- $(ARGS)
