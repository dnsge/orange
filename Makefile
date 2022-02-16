.PHONY: asm vm all mult generate

all: asm vm

asm:
	go build -o ./out/orangeasm ./cmd/orangeasm

vm:
	go build -o ./out/orangevm ./cmd/orangevm

mult: all
	./out/orangeasm ./programs/multiplication.orange ./mult.out
	./out/orangevm ./mult.out

generate:
	go generate ./...
