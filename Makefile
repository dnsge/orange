.PHONY: asm vm

asm:
	go build -o ./out/orangeasm ./cmd/orangeasm

vm:
	go build -o ./out/orangevm ./cmd/orangevm
