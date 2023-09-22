run: 
	@go run main.go

build:
	@go build -o bin/gm -v

deploy:
	@cp bin/gm /usr/local/bin
	@echo "Deployed to /usr/local/bin"
	@rm -rf bin

bad: build deploy

