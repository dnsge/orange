package main

import (
	"fmt"
	"github.com/dnsge/orange/internal/asm"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		_, _ = fmt.Fprintf(os.Stderr, "usage: %s [input file] [output file]\n", os.Args[0])
		os.Exit(1)
		return
	}

	inputFile, err := os.Open(os.Args[1])
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to open input file: %v\n", err)
		os.Exit(1)
		return
	}

	defer inputFile.Close()

	outputFile, err := os.Create(os.Args[2])
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to open output file: %v\n", err)
		os.Exit(1)
		return
	}

	defer outputFile.Close()

	err = asm.AssembleStream(inputFile, outputFile)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
		return
	}
}
