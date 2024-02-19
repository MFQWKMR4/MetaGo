package codegen_test

import (
	"fmt"
	"testing"

	"github.com/MFQWKMR4/MetaGo/pkg/codegen"
)

type SimpleStruct struct {
	Field1 int
	Field2 string
	Field3 []int
	Field4 map[string]int
	Field5 *string
}

type SimpleStruct2 struct {
	Field1 int
	Field2 string
	Field3 []int
	Field4 map[string]int
	Field5 *string
}

func (s SimpleStruct2) DoSomething() {
	fmt.Println("DoSomething")
}

type NestedStruct struct {
	Num            int
	Text           string
	StructField    SimpleStruct
	StructFieldPtr *SimpleStruct
}

type SampleInterface interface {
	DoSomething()
}

type NestedStruct2 struct {
	Num            int
	Text           string
	StructField    SimpleStruct
	StructFieldPtr *SimpleStruct
	InterfaceField SampleInterface
}

func TestCustomPrintf(t *testing.T) {
	tests := []struct {
		name  string
		caseN int
		one   codegen.Fmt[NestedStruct]
		two   codegen.Fmt[NestedStruct2]
	}{
		{
			name:  "test",
			caseN: 1,
			one: codegen.Format(NestedStruct{
				Num:  1,
				Text: "test",
				StructField: SimpleStruct{
					Field1: 1,
					Field2: "test",
					Field3: []int{1, 2, 3},
					Field4: map[string]int{"test": 1},
					Field5: nil,
				},
				StructFieldPtr: &SimpleStruct{
					Field1: 1,
					Field2: "test",
					Field3: []int{1, 2, 3},
					Field4: map[string]int{"test": 1},
					Field5: nil,
				},
			}),
		},
		{
			name:  "test",
			caseN: 2,
			two: codegen.Format(NestedStruct2{
				Num:  1,
				Text: "test",
				StructField: SimpleStruct{
					Field1: 1,
					Field2: "test",
					Field3: []int{1, 2, 3},
					Field4: map[string]int{"test": 1},
					Field5: nil,
				},
				StructFieldPtr: &SimpleStruct{
					Field1: 1,
					Field2: "test",
					Field3: []int{1, 2, 3},
					Field4: map[string]int{"test": 1},
					Field5: nil,
				},
				InterfaceField: SimpleStruct2{
					Field1: 1,
					Field2: "test",
					Field3: []int{1, 2, 3},
					Field4: map[string]int{"test": 1},
					Field5: nil,
				},
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.caseN {
			case 1:
				fmt.Printf("\n %#g \n", tt.one)
			case 2:
				fmt.Printf("\n %#g \n", tt.two)
			}
		})
	}
}
