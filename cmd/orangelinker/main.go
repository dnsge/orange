package main

import (
	"flag"
	"fmt"
	"github.com/dnsge/orange/internal/linker"
	"io"
	"os"
)

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		_, _ = fmt.Fprintf(os.Stderr, "usage: %s [input files...] [output file]\n", os.Args[0])
		os.Exit(1)
		return
	}

	inputFiles := make([]io.Reader, len(args)-1)
	defer func() {
		for _, f := range inputFiles {
			if f != nil {
				if fc, ok := f.(io.ReadCloser); ok {
					_ = fc.Close()
				}
			}
		}
	}()

	for i := range inputFiles {
		f, err := os.Open(args[i])
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to open input file: %v\n", err)
			os.Exit(1)
			return
		}
		inputFiles[i] = f
	}

	outputFile, err := os.Create(args[len(args)-1])
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to open output file: %v\n", err)
		os.Exit(1)
		return
	}

	defer outputFile.Close()

	err = linker.Link(inputFiles, outputFile)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
		return
	}
}
