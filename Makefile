ifeq ($(OS),Windows_NT)
HOSTNAME := $(COMPUTERNAME)
else
HOSTNAME := $(shell hostname)
endif

PORT ?= 3000

ifeq ($(HOSTNAME),KKO11PC)
	DATABASE_DSN ?= host=localhost port=5432 user=postgres password=password dbname=confidant sslmode=disable
else ifeq ($(HOSTNAME),GVS)
	DATABASE_DSN ?= host=localhost port=5433 user=postgres password=password dbname=confidant sslmode=disable
else ifeq ($(HOSTNAME),SURFACE-VEL)
	DATABASE_DSN ?= host=localhost port=5434 user=postgres password=password dbname=confidant sslmode=disable
endif

migrate_up:
ifeq ($(HOSTNAME),KKO11PC)
	@echo "migrating db..."
	migrate -path resources/migration -database "postgres://localhost:5432/confidant?sslmode=disable&user=postgres&password=password" up
	@echo "migrating db completed"
else ifeq ($(HOSTNAME),GVS)
	@echo "migrating db..."
	migrate -path resources/migration -database "postgres://localhost:5433/confidant?sslmode=disable&user=postgres&password=password" up
	@echo "migrating db completed"
else ifeq ($(HOSTNAME),SURFACE-VEL)
	@echo "migrating db..."
	migrate -path resources/migration -database "postgres://localhost:5434/confidant?sslmode=disable&user=postgres&password=password" up
	@echo "migrating db completed"
else
	@echo "no hostname"
endif

migrate_down:
ifeq ($(HOSTNAME),KKO11PC)
	@echo "migrating db..."
	migrate -path resources/migration -database "postgres://localhost:5432/confidant?sslmode=disable&user=postgres&password=password" down
	@echo "migrating db completed"
else ifeq ($(HOSTNAME),GVS)
	@echo "migrating db..."
	migrate -path resources/migration -database "postgres://localhost:5433/confidant?sslmode=disable&user=postgres&password=password" down
	@echo "migrating db completed"
else ifeq ($(HOSTNAME),SURFACE-VEL)
	@echo "migrating db..."
	migrate -path resources/migration -database "postgres://localhost:5434/confidant?sslmode=disable&user=postgres&password=password" down
	@echo "migrating db completed"
else
	@echo "no hostname"
endif

clear_server: migrate_down migrate_up

clear_client:
	@echo "Deleting SQLite DB..."
	if [ -f "confidant_client.db" ]; then \
		rm -v confidant_client.db; \
		echo "confidant_client.db has been deleted"; \
	else \
		echo "confidant_client.db does not exist"; \
	fi
	@echo "Deleting SQLite DB completed"

build_server:
	@echo "Building server..."
	go build ./cmd/confidant_server
	@echo "Building server completed"

build_client:
	@echo "Building client..."
	go build ./cmd/confidant_client
	@echo "Building client completed"

run_server:
	@echo "Running server..."
	start "" "C:\\Program Files\\Git\\bin\\bash.exe" -c "./confidant_server -a ':$(PORT)' -d '$(DATABASE_DSN)'; exec bash"
	@echo "Running server completed"

run_client:
	@echo "Running client..."
	start "" "C:\\Program Files\\Git\\bin\\bash.exe" -c "./confidant_client -a 'localhost:$(PORT)'; exec bash"
	@echo "Running client completed"

build_run_server: build_server run_server

build_run_client: build_client run_client

build_run_server_clear: clear_server build_run_server

build_run_client_clear: clear_client build_run_client

build_run_server_client: build_run_server build_run_client

build_run_server_client_clear: build_run_server_clear build_run_client_clear

run_server_client: run_server run_client

run_server_clear: clear_server run_server

run_client_clear: clear_client run_client
