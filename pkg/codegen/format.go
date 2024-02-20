package codegen

import (
	"fmt"
	"reflect"
	"strings"
)

type Fmt[T any] struct {
	Original T
}

func Format[T any](t T) Fmt[T] {
	return Fmt[T]{Original: t}
}

// implement fmt.Formatter interface
func (t Fmt[T]) Format(f fmt.State, c rune) {
	if c == 'g' && f.Flag('#') {
		fmt.Fprintf(f, "%s", StructCodeGen(t.Original))
	} else {
		fmt.Fprintf(f, "%#v", t.Original)
	}
}

func switchCodeGen[T any](t T) string {
	val := reflect.ValueOf(t)
	// typ := reflect.TypeOf(t)
	switch val.Kind() {
	case reflect.Struct:
		return StructCodeGen(t)
	default:
		return fmt.Sprintf("%#v", t)
	}
}

func StructCodeGen[T any](t T) string {
	var builder strings.Builder

	val := reflect.ValueOf(t)
	typ := reflect.TypeOf(t)

	builder.WriteString(fmt.Sprintf("%s{\n", typ.Name()))

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		switch field.Kind() {
		case reflect.Struct:
			builder.WriteString(fmt.Sprintf("\t%s: %s,\n", fieldType.Name, StructCodeGen(field.Interface())))
		case reflect.Interface:
			builder.WriteString(fmt.Sprintf("\t%s: %s,\n", fieldType.Name, StructCodeGen(field.Interface())))
		case reflect.Ptr:
			if field.IsNil() {
				builder.WriteString(fmt.Sprintf("\t%s: nil,\n", fieldType.Name))
				continue
			}
			value := reflect.Indirect(field)
			builder.WriteString(fmt.Sprintf("\t%s: &%s,\n", fieldType.Name, switchCodeGen(value.Interface())))
		default:
			builder.WriteString(fmt.Sprintf("\t%s: %#v,\n", fieldType.Name, field.Interface()))
		}
	}
	builder.WriteString("}")
	return builder.String()
}
