BINARY_NAME=pdf2q.exe
BINARY_DIR=bin
MAIN=cmd

build:
	GOARCH=amd64 GOOS=windows go build -o ./${BINARY_DIR}/${BINARY_NAME} ./${MAIN}

run: build
	./${BINARY_DIR}/${BINARY_NAME}

clean:
	go clean
	rm ./${BINARY_DIR}/${BINARY_NAME}