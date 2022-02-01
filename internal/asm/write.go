package asm

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"
)

type labelMap map[string]uint32

type assemblyContext struct {
	labels   labelMap
	currLine uint32
}

func AssembleStream(inputFile io.Reader, outputFile io.Writer) error {
	aContext := assemblyContext{
		labels:   make(labelMap),
		currLine: 0,
	}

	rawData, err := io.ReadAll(inputFile)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(rawData)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " \t")
		if line[0] == '.' && line[len(line)-1] == ':' { // label
			labelName := line[1 : len(line)-1]
			if _, ok := aContext.labels[labelName]; ok {
				return fmt.Errorf("duplicate label definition for %q", labelName)
			}
			aContext.labels[labelName] = aContext.currLine
		} else {
			aContext.currLine++
		}
	}

	fmt.Printf("%#v\n", aContext.labels)

	_, _ = reader.Seek(0, io.SeekStart)
	scanner = bufio.NewScanner(reader)

	aContext.currLine = 0
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " \t")
		if line[0] == '.' {
			continue
		}
		if assembled, err := aContext.ParseAssembly(line); err != nil {
			return err
		} else if err = binary.Write(outputFile, binary.LittleEndian, assembled); err != nil {
			return err
		}
		aContext.currLine++
	}
	return nil
}
