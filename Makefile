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
	docker container exec -it movie_night-database-1 psql "postgres://$$MOVIE_NIGHTS_DB_USER:$$(printf %s $$MOVIE_NIGHTS_DB_PASS | jq -sRr @uri)@$$MOVIE_NIGHTS_DB_HOST:$$MOVIE_NIGHTS_DB_PORT/$$MOVIE_NIGHTS_DB_NAME?$$MOVIE_NIGHTS_DB_ARGS"

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

# ---------------------------------------------------
# Deploy
# ---------------------------------------------------
.PHONY: app/deploy/build
app/deploy/build:
	@echo "Compiling Templ templates"
	@templ generate
	@echo "Templ templates compiled successfully"
	@echo "Building project to ./bin/movie_nights"
	@CGO_ENABLED=0 go build -o ./bin/movie_nights
	@echo "Build was successful"
	@echo "Copying over binary to server"
	@rsync -rP --delete ./bin/movie_nights movie_nights@movie-nights.enthys.com:~
	@echo "Binary upload to server was a success"

.PHONY: app/deploy/migrations
app/deploy/migrations:
	@echo "Copying over migration files to server"
	@rsync -rP --delete ./migrations movie_nights@movie-nights.enthys.com:~
	@echo "Running migration files on server"
	@ssh -t movie_nights@movie-nights.enthys.com 'migrate -path ~/migrations -database "postgres://$$MOVIE_NIGHTS_DB_USER:$$(printf %s $$MOVIE_NIGHTS_DB_PASS | jq -sRr @uri)@$$MOVIE_NIGHTS_DB_HOST:$$MOVIE_NIGHTS_DB_PORT/$$MOVIE_NIGHTS_DB_NAME?$$MOVIE_NIGHTS_DB_ARGS" up'
	@echo "Migration was successful"

.PHONY: app/deploy/app.service
app/deploy/app.service:
	@echo "Copying service configuration file to server"
	@rsync -P ./remote/production/movie_nights.service movie_nights@movie-nights.enthys.com:~
	@echo "Applying configuration file and restarting service"
	@ssh -t movie_nights@movie-nights.enthys.com 'sudo mv ~/movie_nights.service /etc/systemd/system/ && sudo systemctl enable movie_nights && sudo systemctl restart movie_nights'
	@echo "Application service restarted"

.PHONY: app/deploy/static
app/deploy/static:
	@echo "Copying over static files"
	@rsync -rP --delete ./assets movie_nights@movie-nights.enthys.com:~
	@echo "Static files have been copied over"

.PHONY: app/deploy
app/deploy: app/deploy/build app/deploy/migrations app/deploy/app.service app/deploy/static

.PHONY: app/deploy
app/deploy/caddyfile/init:
	@rsync -P ./remote/setup/caddy.sh root@movie-nights.enthys.com:~
	@ssh -t root@movie-nights.enthys.com 'sh caddy.sh && rm caddy.sh && systemctl restart caddy'

.PHONY: app/deploy/caddyfile
app/deploy/caddyfile:
	@echo "Uploading remote/production/movie-nights.Caddyfile to server"
	@rsync -P ./remote/production/movie-nights.Caddyfile root@movie-nights.enthys.com:/etc/caddy/
	@echo "Reloading caddy"
	@ssh -t movie_nights@movie-nights.enthys.com 'sudo systemctl restart caddy'
	@echo "Caddy has been updated"