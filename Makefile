BINARY_NAME=unifi-prune

build:
	cd ./cmd/prune && go build -o ${BINARY_NAME} *.go && mv ${BINARY_NAME} ../../

run:
	cd ./cmd/prune && go run *.go

format:
	gofmt -s -w .

clean:
	rm ${BINARY_NAME}