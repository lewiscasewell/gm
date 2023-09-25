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
	@GOOS=linux GOARCH=amd64 go build -o bin/gm -v
	@GOOS=windows GOARCH=amd64 go build -o bin/gm.exe -v
	@GOOS=darwin GOARCH=amd64 go build -o bin/gm-mac -v

move-to-bin:
	@cp bin/gm /usr/local/bin

delete-bin:
	@rm -rf bin

it: build move-to-bin delete-bin

tar:
	@tar -czvf gm.tar.gz bin/gm

zip:
	@zip -r gm.zip bin/gm.exe

install:
	@go install
