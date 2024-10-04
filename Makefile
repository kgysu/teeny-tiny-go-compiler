SHELL=/bin/bash

build:
	go run ./cmd/tiny-compiler/main.go app.teeny
	gcc -o app out.c

run:
	./app


