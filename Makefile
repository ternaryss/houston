app = houston
dsn = ./data/app.db

all: hotreload

build: clean volumes
	@echo "--- Building ---"
	@npm install
	@npx tailwindcss -i ./views/tailwind.css -o ./public/css/styles.css
	@templ generate
	@go build -o ./bin/$(app) ./cmd/$(app)/main.go

clean:
	@echo "--- Cleaning ---"
	@rm -rf ./public
	@rm -rf ./bin

down: volumes
	@echo "--- Migrations down ---"
	@GOOSE_DRIVER=sqlite3 GOOSE_DBSTRING=$(dsn) goose -dir=./migrations down

help:
	@echo "all - runs application in Hot Reload mode"
	@echo "build - build application"
	@echo "clean - clean filesystem after application build"
	@echo "down - rollback last database migration"
	@echo "hotreload - runs application in Hot Reload mode"
	@echo "run - run application"
	@echo "status - check database migrations status"
	@echo "up - apply databases migrations"
	@echo "volumes - prepare filesystem to build & run application"

hotreload: clean volumes up
	@echo "--- Running in Hot Reload mode ---"
	@air

run: build up
	@echo "--- Running ---"
	@./bin/$(app)

status: volumes
	@echo "--- Migrations status ---"
	@GOOSE_DRIVER=sqlite3 GOOSE_DBSTRING=$(dsn) goose -dir=./migrations status

up: volumes
	@echo "--- Migrations up ---"
	@GOOSE_DRIVER=sqlite3 GOOSE_DBSTRING=$(dsn) goose -dir=./migrations up

volumes:
	@echo "--- Volumes ---"
	@mkdir -p ./bin
	@mkdir -p ./data
	@touch ./data/app.db
	@mkdir -p ./public
	@cp -r ./static/. ./public
