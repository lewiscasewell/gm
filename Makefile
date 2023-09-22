run: 
	@go run main.go

build:
	@go build -o bin/gm -v

deploy:
	@cp bin/gm /usr/local/bin
