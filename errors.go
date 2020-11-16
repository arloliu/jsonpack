package jsonpack

import (
	"errors"
	"fmt"
)

var (
	errPropertiesProp = errors.New("'properties' property non-exist or invalid")
	errItemsProp      = errors.New("'items' property non-exist or invalid")
	errOrderProp      = errors.New("'order' property non-exist or invalid")
	errTypeProp       = errors.New("'type' property non-exist or invalid")
)

// NotImplementedError is returned when the operation in handler doesn't implemented.
type NotImplementedError struct {
	Name string // operation name that not implemented by handler
}

func (e *NotImplementedError) Error() string {
	return fmt.Sprintf("'%s' is not implemented", e.Name)
}

// WrongTypeError is returned when wrong data type found in data, it happens in Encode/Decode methods.
type WrongTypeError struct {
	DataType string // wrong type of data
}

func (e *WrongTypeError) Error() string {
	return fmt.Sprintf("wrong data type '%s'", e.DataType)
}

// UnknownTypeError is returned when un-supported data type found, it happends in AddSchema method.
type UnknownTypeError struct {
	DataType string // data type that un-supported by schema defintion
}

func (e *UnknownTypeError) Error() string {
	return fmt.Sprintf("unknown data type '%s'", e.DataType)
}

// TypeAssertionError indicates an error that occurs when doing data type asserrion.
type TypeAssertionError struct {
	Data         interface{} // the data that failed on type assertion
	ExpectedType string      // expected data type
}

func (e *TypeAssertionError) Error() string {
	return fmt.Sprintf("got data of type %T but wanted %s", e.Data, e.ExpectedType)
}

// StructFieldNonExistError indicates an error that occurs when structure doesn't have required field.
type StructFieldNonExistError struct {
	Name  string // name of structure
	Field string // field of structure that expects to be existed
}

func (e *StructFieldNonExistError) Error() string {
	return fmt.Sprintf("struct %s doesn't contain required field: '%s'", e.Name, e.Field)
}

// SchemaNonExistError indicates an error that occurs when pre-compiled schema defintion does not exist.
type SchemaNonExistError struct {
	Name string // schema name
}

func (e *SchemaNonExistError) Error() string {
	return fmt.Sprintf("schema definition '%s' does not exist", e.Name)
}

// CompileError represents an error from calling AddSchema method, it indicates there has an error occurs
// in compiling procedure of schema definition.
type CompileError struct {
	Name string // schema name
	Err  error  // actual error
}

func (e *CompileError) Error() string {
	return fmt.Sprintf("compiling '%s' got error: %v", e.Name, e.Err.Error())
}

// Unwrap returns the underlying error.
func (e *CompileError) Unwrap() error {
	return e.Err
}

// EncodeError represents an error from calling Encode or Marshal methods.
type EncodeError struct {
	Name string // schema name
	Err  error  // actual error
}

func (e *EncodeError) Error() string {
	return fmt.Sprintf("encode with schema definition '%s' got error: %v", e.Name, e.Err.Error())
}

// Unwrap returns the underlying error.
func (e *EncodeError) Unwrap() error {
	return e.Err
}

// DecodeError represents an error from calling Decode or Unmarshal methods.
type DecodeError struct {
	Name string // schema name
	Err  error  // actual error
}

func (e *DecodeError) Error() string {
	return fmt.Sprintf("decode with schema definition '%s' got error: %v", e.Name, e.Err.Error())
}

// Unwrap returns the underlying error.
func (e *DecodeError) Unwrap() error {
	return e.Err
}
