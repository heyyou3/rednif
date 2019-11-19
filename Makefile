.PHONY: install build
install:
	go install ./cmd/rednif/
build:
	go build -o ./rednif ./cmd/rednif/main.go
