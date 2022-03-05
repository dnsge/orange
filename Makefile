.PHONY: asm vm linker all mult generate

all: asm vm linker

clean:
	rm *.obj *.out

asm:
	go build -o ./out/orangeasm ./cmd/orangeasm

vm:
	go build -o ./out/orangevm ./cmd/orangevm

linker:
	go build -o ./out/orangelinker ./cmd/orangelinker

mult: all
	./out/orangeasm ./programs/multiplication.orange ./mult.out
	./out/orangevm ./mult.out

generate:
	go generate ./...

%.orange: all
	./out/orangeasm --executable ./programs/$*.orange ./$*.out
	./out/orangevm ./$*.out

link: asm linker
	./out/orangeasm ./programs/link/main.orange ./main.obj
	./out/orangeasm ./programs/link/strlen.orange ./strlen.obj
	./out/orangelinker ./main.obj ./strlen.obj ./strlen_main.out

greet: asm linker stdlib
	./out/orangeasm ./programs/greet/greet.orange ./greet.obj
	./out/orangelinker ./greet.obj ./std_strio.obj ./greet.out

stdlib:
	./out/orangeasm ./programs/std/strio.orange ./std_strio.obj
