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

// NotImplemented means this method/operation is not implemented
type NotImplementedError struct {
	name string
}

func (e *NotImplementedError) Error() string {
	return fmt.Sprintf("'%s' is not implemented", e.name)
}

// WrongTypeError means there has a wrong data type when encoding or decoding data
// with pre-compiled schema definition
type WrongTypeError struct {
	dataType string
}

func (e *WrongTypeError) Error() string {
	return fmt.Sprintf("wrong data type '%s'", e.dataType)
}

// UnknownTypeError means there has a un-supported data type when compiling schema definition
type UnknownTypeError struct {
	dataType string
}

func (e *UnknownTypeError) Error() string {
	return fmt.Sprintf("unknown data type '%s'", e.dataType)
}

// TypeAssertionError indicates an error the that occurred when doing data
// data type asserrion
type TypeAssertionError struct {
	data       interface{}
	expectType string
}

func (e *TypeAssertionError) Error() string {
	return fmt.Sprintf("got data of type %T but wanted %s", e.data, e.expectType)
}

// InvalidPropValueError indicates an error the that occurred when reading or
// writing property got invalid value
type InvalidPropValueError struct {
	name string
	data interface{}
}

func (e *InvalidPropValueError) Error() string {
	return fmt.Sprintf("property %s has invalid value %v", e.name, e.data)
}

// StructFieldNonExistError indicates an error the that occurred when structure
// doesn't contain required field
type StructFieldNonExistError struct {
	name  string
	field string
}

func (e *StructFieldNonExistError) Error() string {
	return fmt.Sprintf("struct %s doesn't contain required field: '%s'", e.name, e.field)
}

// SchemaNonExistError indicates an error that occurred while pre-compiled schema defintion
// does not exist
type SchemaNonExistError struct {
	name string
}

func (e *SchemaNonExistError) Error() string {
	return fmt.Sprintf("schema definition '%s' does not exist", e.name)
}

// CompileError indicates an error that occurred while attempting to compile
// schema definition
type CompileError struct {
	name string
	Err  error
}

func (e *CompileError) Error() string {
	return fmt.Sprintf("compiling '%s' got error: %v", e.name, e.Err.Error())
}

func (e *CompileError) Unwrap() error {
	return e.Err
}

// EncodeError indicates an error that occurred while attempting to encode data with
// pre-compiled schema definition
type EncodeError struct {
	name string
	Err  error
}

func (e *EncodeError) Error() string {
	return fmt.Sprintf("encode with schema definition '%s' got error: %v", e.name, e.Err.Error())
}

func (e *EncodeError) Unwrap() error {
	return e.Err
}

// DecodeError indicates an error that occurred while attempting to decode packed binary data with
// pre-compiled schema definition
type DecodeError struct {
	name string
	Err  error
}

func (e *DecodeError) Error() string {
	return fmt.Sprintf("decode with schema definition '%s' got error: %v", e.name, e.Err.Error())
}

func (e *DecodeError) Unwrap() error {
	return e.Err
}
