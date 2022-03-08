package main

import (
	"bytes"
	"embed"
	"fmt"
	"go/format"
	"io"
	"log"
	"os"
	"strings"
	"text/template"
	"time"
)

//go:embed output.go.tmpl
var outputTemplateFiles embed.FS

type outputContext struct {
	Time            string
	TokenEnumLines  []string
	OpTokens        []string
	AllInstructions []*Instruction
	Categories      []*categoryOutput
}

type categoryOutput struct {
	Title string
	ID    string
}

func main() {
	var outputFile io.Writer
	if len(os.Args) == 2 {
		outputFileName := os.Args[1]
		of, err := os.Create(outputFileName)
		if err != nil {
			log.Fatalf("Failed to open output file: %v\n", err)
		}
		outputFile = of
	} else {
		outputFile = os.Stdout
	}

	allTemplates := template.Must(template.ParseFS(outputTemplateFiles, "*.go.tmpl"))

	categories, categoryMap := mapToCategories(instructionTokens)
	enumLines, opTokens := generateEnumOutput(categoryMap)

	outputCtx := &outputContext{
		Time:            time.Now().Format(time.RFC3339),
		TokenEnumLines:  enumLines,
		AllInstructions: instructionTokens,
		Categories:      categories,
		OpTokens:        opTokens,
	}

	buf := new(bytes.Buffer)

	err := allTemplates.ExecuteTemplate(buf, "output.go.tmpl", outputCtx)
	if err != nil {
		panic(err)
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		log.Println(buf.String())
		log.Fatalf("Failed to format produced source: %v\n", err)
	}

	reader := bytes.NewReader(formatted)
	io.Copy(outputFile, reader)
}

func mapToCategories(in []*Instruction) ([]*categoryOutput, map[string][]*Instruction) {
	var resCat []*categoryOutput
	resMap := make(map[string][]*Instruction)
	for _, instruction := range in {
		lst, ok := resMap[instruction.TokenCategory]
		if !ok {
			lst = nil
			if instruction.TokenCategory != NoCategory {
				resCat = append(resCat, &categoryOutput{
					Title: strings.Title(instruction.TokenCategory),
					ID:    instruction.TokenCategory,
				})
			}
		}
		resMap[instruction.TokenCategory] = append(lst, instruction)
	}
	return resCat, resMap
}

func generateEnumOutput(categoryMap map[string][]*Instruction) ([]string, []string) {
	usedCategories := make(map[string]struct{})
	used := make(map[string]struct{})

	var lines []string
	var opStrings []string

	emitAllForCategory := func(category string) {
		iList := categoryMap[category]

		if category != NoCategory {
			lines = append(lines, fmt.Sprintf("_%sStart", category))
		}
		for _, ins := range iList {
			name := ins.EnumName()

			if _, ok := used[name]; ok {
				// already dealt with
				continue
			}
			used[name] = struct{}{}

			lines = append(lines, ins.EnumName())
			if ins.IsRealOp() {
				opStrings = append(opStrings, ins.EnumName())
			}
		}
		if category != NoCategory {
			lines = append(lines, fmt.Sprintf("_%sEnd", category))
		}
	}

	lines = append(lines, "_invalid TokenKind = iota")
	for _, orderedName := range enumOrder {
		emitAllForCategory(orderedName)
		usedCategories[orderedName] = struct{}{}
	}

	for category := range categoryMap {
		if _, ok := usedCategories[category]; ok {
			continue
		}
		emitAllForCategory(category)
	}

	return lines, opStrings
}
