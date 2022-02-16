package main

import (
	"fmt"
	"github.com/dnsge/orange/internal/memory"
	"github.com/dnsge/orange/internal/vm"
	"os"
)

const (
	stackBottom = 0x7FFF0000
	stackSize   = 0x10000
)

func main() {
	if len(os.Args) != 2 {
		_, _ = fmt.Fprintf(os.Stderr, "usage: %s [input file]\n", os.Args[0])
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

	mem := memory.New()
	sim := vm.NewVirtualMachine(mem)
	_, err = mem.LoadFromReader(0, inputFile)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to load input file into memory: %v\n", err)
		os.Exit(1)
		return
	}

	mem.Alloc(stackBottom+1, stackSize)
	sim.InitStack(stackBottom + stackSize)

	sim.PrintState()
	for !sim.Halted() {
		sim.ExecuteInstruction()
		sim.PrintState()
	}
}
