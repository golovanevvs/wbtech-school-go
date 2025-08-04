build-server:
	@echo "Building server..."
	@go build ./cmd/order-service-mainServer
	@echo "Building server completed"

run-server:
	@echo "Running server..."
	@start "" "/C/Program Files/Git/git-bash.exe" -c "./order-service-mainServer"
	@echo "Running server completed"

build-run-server: build-server run-server