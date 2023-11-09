PATH_PRUNE=./cmd/prune
PATH_FILTER=./cmd/filter
BINARY_NAME_PRUNE=prune
BINARY_NAME_FILTER=filter

build_prune:
	cd ${PATH_PRUNE} && go build -o ${BINARY_NAME_PRUNE} *.go && mv ${BINARY_NAME_PRUNE} ../../

build_filter:
	cd ${PATH_FILTER} && go build -o ${BINARY_NAME_FILTER} *.go && mv ${BINARY_NAME_FILTER} ../../

run_prune:
	cd ${PATH_PRUNE} && go run *.go

run_filter:
	cd ${PATH_FILTER} && go run *.go

format:
	gofmt -s -w .

clean:
	rm -f ${BINARY_NAME_PRUNE}
	rm -f ${BINARY_NAME_FILTER}
