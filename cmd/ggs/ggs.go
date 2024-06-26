package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func parseArgs() (
	inFilename string,
) {

	inputFilePtr := flag.String("input", "", "go input file path")
	flag.Parse()
	inFilename = *inputFilePtr
	println("Input file:", inFilename)
	return
}

func main() {
	filePath := parseArgs()
	outFilePath := strings.TrimSuffix(filePath, ".go") + "_gen.go"
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: failed to read file %s: %v\n", filePath, err)
		os.Exit(1)
	}
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	sp := NewStructParser(fset, fileContent)

	bld := strings.Builder{}
	bld.WriteString(GeneratePackage(astFile))
	bld.WriteString(GenerateImports(astFile))

	ast.Inspect(astFile, func(n ast.Node) bool {
		ts, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}
		st, ok := ts.Type.(*ast.StructType)
		if !ok {
			return true
		}

		structName := ts.Name.Name
		for _, field := range st.Fields.List {
			fieldTypeText := sp.fieldTypeText(field)
			for _, fieldName := range field.Names {
				structField := StructField{
					StructName:    structName,
					FieldName:     fieldName.Name,
					FieldTypeText: fieldTypeText,
				}
				if sp.fieldGetter(field) {
					bld.WriteString(structField.GenerateGS())
				}
			}
		}
		return true
	})

	result := bld.String()
	if err = os.WriteFile(outFilePath, []byte(result), os.FileMode(0644)); err != nil {
		panic(err)
	}
}

func NewStructParser(fileSet *token.FileSet, fileContent []byte) StructParser {
	return StructParser{
		fileSet:                fileSet,
		fileContent:            fileContent,
		whitespaceRegexp:       regexp.MustCompile(`\s+`),
		flagGetterSetterRegexp: regexp.MustCompile(`\b+ggs\b`),
	}
}

type StructParser struct {
	fileSet                *token.FileSet
	fileContent            []byte
	whitespaceRegexp       *regexp.Regexp
	flagGetterSetterRegexp *regexp.Regexp
}

func (sp *StructParser) fieldGetter(field *ast.Field) bool {
	return sp.flagGetterSetterRegexp.MatchString(field.Comment.Text())
}

func (sp *StructParser) fieldTypeText(field *ast.Field) string {
	begin := sp.fileSet.Position(field.Type.Pos()).Offset
	end := sp.fileSet.Position(field.Type.End()).Offset
	return sp.whitespaceRegexp.ReplaceAllString(string(sp.fileContent[begin:end]), " ")
}

type StructField struct {
	StructName    string
	FieldName     string
	FieldTypeText string
}

func GeneratePackage(astFile *ast.File) string {
	bld := &strings.Builder{}
	bld.WriteString("// Code generated by ggs; DO NOT EDIT.\n\n")
	bld.WriteString(fmt.Sprintf("package %s\n\n", astFile.Name.Name))
	return bld.String()
}

func GenerateImports(astFile *ast.File) string {
	bld := &strings.Builder{}
	bld.WriteString("import (\n")
	for _, i := range astFile.Imports {
		bld.WriteString(fmt.Sprintf("\t%s\n", i.Path.Value))
	}
	bld.WriteString(")\n\n")
	return bld.String()
}

func (sf *StructField) GenerateGS() string {
	upper := cases.Title(language.English).String(sf.FieldName)
	return fmt.Sprintf(`
func Get%s(v *%s) %s {
	return v.%s
}

func Set%s(v *%s, value %s) {
	v.%s = value
}
`, upper, sf.StructName, sf.FieldTypeText, sf.FieldName,
		upper, sf.StructName, sf.FieldTypeText, sf.FieldName)
}
