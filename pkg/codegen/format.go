package codegen

import (
	"fmt"
	"os"
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
		fmt.Fprintf(f, "%s", switchCodeGen(t.Original))
	} else {
		fmt.Fprintf(f, "%#v", t.Original)
	}
}

func WriteFile[T any](filePath string, content T) error {
	str := fmt.Sprintf("%#g\n", Format(content))
	err := os.WriteFile(filePath, []byte(str), 0644)
	if err != nil {
		return err
	}
	return nil
}

func switchCodeGen[T any](t T) string {
	val := reflect.ValueOf(t)
	switch val.Kind() {
	case reflect.Struct:
		return StructCodeGen(t)
	case reflect.Interface:
		return StructCodeGen(t)
	case reflect.Ptr:
		if val.IsNil() {
			return "nil,"
		}
		value := reflect.Indirect(val)
		return switchCodeGen(value)
	case reflect.Slice:
		return SliceCodeGen(t)
	default:
		return fmt.Sprintf("%#v", t)
	}
}

func SliceCodeGen[T any](t T) string {

	rv := reflect.ValueOf(t)

	// if rv.Len() == 0 {
	// 	return "nil"
	// }

	var builder strings.Builder

	builder.WriteString("[]")

	tt := rv.Type().Elem()
	elemType := tt.String()
	builder.WriteString(elemType)
	builder.WriteString("{\n")

	for i := 0; i < rv.Len(); i++ {
		field := rv.Index(i)

		if field.CanInterface() {
			switch field.Kind() {
			case reflect.Struct:
				builder.WriteString(fmt.Sprintf("\t%s,\n", StructCodeGen(field.Interface())))
			case reflect.Interface:
				builder.WriteString(fmt.Sprintf("\t%s,\n", StructCodeGen(field.Interface())))
			case reflect.Ptr:
				if field.IsNil() {
					continue
				}
				value := reflect.Indirect(field)
				v2 := value.Interface()
				// val := reflect.ValueOf(v2)
				typ := reflect.TypeOf(v2)

				// typeName := reflect.TypeOf(value.Interface()).Name()
				// if typeName == "" {
				// 	typeName = reflect.TypeOf(value.Interface()).Elem().Name()
				// }

				builder.WriteString(fmt.Sprintf("\tlo.ToPtr[%s](%s),\n", typ.String(), switchCodeGen(v2)))
			case reflect.Slice:
				builder.WriteString(fmt.Sprintf("\t%s,\n", switchCodeGen(field.Interface())))
			case reflect.Array:
				builder.WriteString(fmt.Sprintf("\t%s,\n", switchCodeGen(field.Interface())))
			default:
				builder.WriteString(fmt.Sprintf("\t%#v,\n", field.Interface()))
			}
		} else {
			builder.WriteString("\tlo.ToPtr({}),\n")
		}

	}
	builder.WriteString("\t}")
	return builder.String()
}

func StructCodeGen[T any](t T) string {
	var builder strings.Builder

	val := reflect.ValueOf(t)
	typ := reflect.TypeOf(t)

	builder.WriteString(fmt.Sprintf("%s{\n", typ.String()))

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		if field.CanInterface() {
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

				v2 := value.Interface()
				// val := reflect.ValueOf(v2)
				typ := reflect.TypeOf(v2)

				builder.WriteString(fmt.Sprintf("\t%s: lo.ToPtr[%s](%s),\n", fieldType.Name, typ.String(), switchCodeGen(value.Interface())))
			case reflect.Slice:
				builder.WriteString(fmt.Sprintf("\t%s: %s,\n", fieldType.Name, SliceCodeGen(field.Interface())))
			case reflect.Array:
				builder.WriteString(fmt.Sprintf("\t%s: %s,\n", fieldType.Name, SliceCodeGen(field.Interface())))
			default:
				builder.WriteString(fmt.Sprintf("\t%s: %#v,\n", fieldType.Name, field.Interface()))
			}
		} else {
			builder.WriteString("")
		}
	}
	builder.WriteString("}")
	return builder.String()
}
