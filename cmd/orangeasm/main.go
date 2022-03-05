package main

import (
	"flag"
	"fmt"
	"github.com/dnsge/orange/internal/asm"
	"os"
)

var (
	executableFlag = flag.Bool("executable", false, "Compile to executable")
)

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		_, _ = fmt.Fprintf(os.Stderr, "usage: %s [input file] [output file]\n", os.Args[0])
		os.Exit(1)
		return
	}

	inputFile, err := os.Open(args[0])
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to open input file: %v\n", err)
		os.Exit(1)
		return
	}

	defer inputFile.Close()

	outputFile, err := os.Create(args[1])
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to open output file: %v\n", err)
		os.Exit(1)
		return
	}

	defer outputFile.Close()

	if *executableFlag {
		err = asm.AssembleExecutable(inputFile, outputFile)
	} else {
		err = asm.AssembleObjectFile(inputFile, outputFile)
	}

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
		return
	}
}
