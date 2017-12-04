DEFAULT_ARGS := data/config.json

test:
	go test -v *.go

clean:
	go clean

build: clean
	go build -v -o crawler *.go

run:
	go run agent.go crawler.go util.go main.go ${DEFAULT_ARGS}

run_custom: build
	go run clawler ${ARGS}

start: clean build
	./crawler ${DEFAULT_ARGS}


