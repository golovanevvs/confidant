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
ifeq ($(HOSTNAME),SURFACE-VEL)
	@echo "migrating db..."
	migrate -path resources/migration -database "postgres://localhost:5434/confidant?sslmode=disable&user=postgres&password=password" up
	@echo "migrating db completed"
else
	@echo "no hostname"
endif
endif

migrate_down:
ifeq ($(HOSTNAME),GVS)
	@echo "migrating db..."
	migrate -path resources/migration -database "postgres://localhost:5433/confidant?sslmode=disable&user=postgres&password=password" down
	@echo "migrating db completed"
else
ifeq ($(HOSTNAME),SURFACE-VEL)
	@echo "migrating db..."
	migrate -path resources/migration -database "postgres://localhost:5434/confidant?sslmode=disable&user=postgres&password=password" down
	@echo "migrating db completed"
else
	@echo "no hostname"
endif
endif

deldb:
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
	start "" "C:\\Program Files\\Git\\bin\\bash.exe" -c "./confidant_server; exec bash"
# start "./confidant_server"
	@echo "Running server completed"

run_client:
	@echo "Running client..."
	start "" "C:\\Program Files\\Git\\bin\\bash.exe" -c "./confidant_client; exec bash"
# start "./confidant_client"
	@echo "Running client completed"

build_run_server: build_server run_server

build_run_client: build_client run_client

build_run_client_clear: deldb build_client run_client

run_client_clear: deldb run_client

build_run_server_client_clear: migrate_down migrate_up build_server run_server deldb build_client run_client

run_server_client: run_server run_client