
run:
# Start the Docker containers and run both server and client
	docker compose up -d 

run-server:
# Start only the Go server
	cd server && go run src/main.go

run-client:
# Start only the client
	cd client && npm run dev

run-parallel:
# Start both server and client in parallel using separate terminals
	@echo "Starting services..."
	docker compose up -d
	@echo "Use 'make run-server' in one terminal and 'make run-client' in another"