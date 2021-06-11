build-web:
	echo "::: Compiling web project"
	go build -o ./cmd/web/gh-web ./cmd/web/main.go

build-cli:
	echo "::: Compiling CLI project"
	go build -o ./cmd/cli/gh-cli ./cmd/cli/main.go


run-web:
	echo "::: Running web project"
	go run ./cmd/web

run-cli:
	echo "::: Running CLI project"
	go run ./cmd/cli -sid=$$sid

all-web: build-web run-web
all-cli: build-cli run-cli