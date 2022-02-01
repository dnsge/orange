package asm

import (
	"bufio"
	"encoding/binary"
	"io"
)

func AssembleStream(inputFile io.Reader, outputFile io.Writer) error {
	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		line := scanner.Text()
		if assembled, err := ParseAssembly(line); err != nil {
			return err
		} else if err = binary.Write(outputFile, binary.LittleEndian, assembled); err != nil {
			return err
		}
	}
	return nil
}
