# MetaGo
Utilities by meta programming.

## cmd/ggs
A cli tool for getter/setter generation.

## pkg/codegen
This library allows for the output of struct definitions using the `Format` method. When using `%#v`, if the type is a pointer, it results in the memory address being outputted. By utilizing this library and specifying `%#g`, it is possible to output the values of pointer types as their corresponding struct representations, enabling outputs that can be directly used as definitions within the source code.

