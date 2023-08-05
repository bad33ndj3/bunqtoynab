include .env
export

test:
	go run gotest.tools/gotestsum@latest ./...

run:
	@echo "Running the application..."
	go run main.go

deploy:
	cp .env package/cmd/ynabtobunq/
	#doctl serverless deploy ../bunqparser --remote-build
