build_server:
	@echo "Build server..."
	go build ./cmd/confidant_server
	@echo "Build server completed"

build_client:
	@echo "Build client..."
	go build ./cmd/confidant_client
	@echo "Build client completed"

run_server:
	@echo "Run server..."
	start "" "C:\\Program Files\\Git\\bin\\bash.exe" -c "./confidant_server; exec bash"
	@echo "Run server completed"

run_client:
	@echo "Run client..."
	start "" "C:\\Program Files\\Git\\bin\\bash.exe" -c "./confidant_client; exec bash"
	@echo "Run client completed"

build_run_server: build_server run_server

build_run_client: build_client run_client

build_run_server_client: build_server run_server build_client run_client