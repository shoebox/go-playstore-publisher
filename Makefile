APP=go-play-publisher

deps:
	go mod download

build:
	go build -o ${APP} cmd/gpp/main.go

install:
	CGO_ENABLED=0 go install ./cmd/gpp/main.go

run:
	go run -race cmd/gpp/main.go

clean:
	rm -rf ${APP}
