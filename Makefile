BINARY_NAME=unifi-prune

build:
	go build -o ${BINARY_NAME} *.go

build_linux:
	GOOS=linux GOARCH=arm64 go build -o ${BINARY_NAME} *.go

run:
	go run *.go

format:
	gofmt -s -w .

clean:
	rm ${BINARY_NAME}