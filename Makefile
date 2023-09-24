run: 
	@go run main.go

test:
	@go test -v ./...

test-nocache:
	@go test -count=1 -v ./...

coverage:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

build:
	@go build -o bin/gm -v

deploy:
	@cp bin/gm /usr/local/bin
	@echo "Deployed to /usr/local/bin"
	@rm -rf bin

bad: build deploy

