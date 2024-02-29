include .env.example
-include .env
export

default: help

## help: prints this help message
.PHONY: help
help: # Show help for each of the Makefile recipes.
	@echo 'Usage:'
	@grep -E '^##'  Makefile | sort | while read -r l; do printf "  \033[1;32m$$(echo $$l | cut -b 4- | cut -f 1 -d':')\033[00m:$$(echo $$l | cut -f 2- -d':')\n"; done


# ---------------------------------------------------
# Development
# ---------------------------------------------------

## run/setup: run required commands in order to start development
.PHONY: run/setup
run/setup:
	@if [ ! -e ".env" ]; then echo 'Copying .env.example under name .env.' && cp .env.example .env; fi
	@if [ ! -e "bin" ]; then echo 'Creating bin directory.' && mkdir bin; fi
	@if [ ! -e "bin/migrate" ]; then echo 'Downloading migrate in ./bin ...'; curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz -C bin migrate; fi

## docker/start: Starts the application through docker
.PHONY: 
run/dev: run/setup
	docker compose up

## run/clear: Removed any existing instance and volume, rebuilds the project without any cache
.PHONY: run/clear
run/clear:
	docker compose down -v
	docker compose build --no-cache

# ---------------------------------------------------
# Database
# ---------------------------------------------------

## db/conenct: Connect to database
.PHONY: db/connect
db/connect: run/setup
	docker container exec -it movie_night-database-1 psql ${MOVIENIGHT_DB_DSN}

## db/migrations/up: Run migrations
.PHONY: db/migrations/up
db/migrations/up: run/setup
	docker container exec -it movie_nights_app ./bin/migrate -path ./migrations -database ${MOVIENIGHT_DB_DSN} -verbose up

## db/migrations/down: Rollback migrations
.PHONY: db/migrations/down
db/migrations/down: run/setup
	docker container exec -it movie_nights_app ./bin/migrate -path ./migrations -database ${MOVIENIGHT_DB_DSN} -verbose down

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new: run/setup
	@echo 'Creating migration files for ${name}'
	@./bin/migrate create -seq -ext=.sql -dir=./migrations ${name}
