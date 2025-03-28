ifeq ($(OS),Windows_NT)
HOSTNAME := $(COMPUTERNAME)
else
HOSTNAME := $(shell hostname)
endif

migrate_up:
ifeq ($(HOSTNAME),GVS)
	@echo "migrating db..."
	migrate -path resources/migration -database "postgres://localhost:5433/confidant?sslmode=disable&user=postgres&password=password" up
	@echo "migrating db completed"
else
	@echo "no hostname"
endif

migrate_down:
ifeq ($(HOSTNAME),GVS)
	@echo "migrating db..."
	migrate -path resources/migration -database "postgres://localhost:5433/confidant?sslmode=disable&user=postgres&password=password" down
	@echo "migrating db completed"
else
	@echo "no hostname"
endif

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
	start "" "C:\\Program Files\\Git\\bin\\bash.exe" -c "./confidant_server; exec bash"
	@echo "Running server completed"

run_client:
	@echo "Running client..."
	start "" "C:\\Program Files\\Git\\bin\\bash.exe" -c "./confidant_client; exec bash"
	@echo "Running client completed"

build_run_server: build_server run_server

build_run_client: build_client run_client

build_run_server_client: build_server run_server build_client run_client