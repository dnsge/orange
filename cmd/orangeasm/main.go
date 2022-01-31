package main

import (
	"github.com/dnsge/orange/internal/asm"
	"log"
)

func main() {
	i, err := asm.ParseAssembly("MOVZ r1, #10")
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Printf("%032b\n", i)
	}
}
