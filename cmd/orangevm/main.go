package main

import (
	"flag"
	"fmt"
	"github.com/dnsge/orange/internal/memory"
	"github.com/dnsge/orange/internal/vm"
	"os"
)

const (
	stackBottom = 0x7FFF0000
	stackSize   = 0x10000
)

var (
	quietFlag = flag.Bool("quiet", false, "Disable printing state")
)

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		_, _ = fmt.Fprintf(os.Stderr, "usage: %s [input file]\n", os.Args[0])
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

	mem := memory.New()
	sim := vm.NewVirtualMachine(mem, *quietFlag)
	_, err = mem.LoadFromReader(0, inputFile)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to load input file into memory: %v\n", err)
		os.Exit(1)
		return
	}

	mem.Alloc(stackBottom+1, stackSize)
	sim.InitStack(stackBottom + stackSize)

	if !*quietFlag {
		sim.PrintState()
	}

	for !sim.Halted() {
		sim.ExecuteInstruction()
		if !*quietFlag {
			sim.PrintState()
		}
	}
}
