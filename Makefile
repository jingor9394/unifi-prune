PATH_PRUNE=./cmd/prune
PATH_FILTER=./cmd/filter
BINARY_NAME_PRUNE=prune
BINARY_NAME_FILTER=filter

build:
	cd ${PATH_PRUNE} && go build -o ${BINARY_NAME_PRUNE} *.go && mv ${BINARY_NAME_PRUNE} ../../
	cd ${PATH_FILTER} && go build -o ${BINARY_NAME_FILTER} *.go && mv ${BINARY_NAME_FILTER} ../../

release:
	zip ${BINARY_NAME_PRUNE}-darwin-${VERSION}.zip ./${BINARY_NAME_PRUNE}
	zip ${BINARY_NAME_FILTER}-darwin-${VERSION}.zip ./${BINARY_NAME_FILTER}

run_prune:
	cd ${PATH_PRUNE} && go run *.go

run_filter:
	cd ${PATH_FILTER} && go run *.go

format:
	gofmt -s -w .

clean:
	rm -f ${BINARY_NAME_PRUNE}
	rm -f ${BINARY_NAME_FILTER}
	rm -f ${BINARY_NAME_PRUNE}-darwin-*.zip
	rm -f ${BINARY_NAME_FILTER}-darwin-*.zip
