package jsonpack

import (
	"errors"
	"fmt"
	"reflect"

	ibuf "github.com/arloliu/jsonpack/internal/buffer"

	"github.com/modern-go/reflect2"
)

/*
Encode returns packed binary data of v with compiled schema definition.
The compiled schema definition is specified by schemaName and
return data will be nil if error occurs.

The type of v can be a map[string]interface{} which presents valid JSON data, or
a struct instance which added by AddSchema function.

The return data contains encoded binary data that can then be decoded by Decode function.

Example of encoding a map type

	data := map[string]interface{} {
		"name": "example name",
		"area": uint32(888),
	}

	jsonPack := jsonpack.NewJSONPack()
	// call jsonPack.AddSchema to register schema
	sch := jsonPack.GetSchema("info")
	result, err := sch.Encode(data)
*/
func (s *Schema) Encode(d interface{}) ([]byte, error) {
	encodeData := make([]byte, s.encodeBufSize)
	err := s.EncodeTo(d, &encodeData)
	return encodeData, err
}

// EncodeTo is similar to Encode function, but passing a pointer to []byte to store
// encoded data instead of returning new allocated []byte encoded data.
func (s *Schema) EncodeTo(d interface{}, dataPtr *[]byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch r := r.(type) {
			case string:
				err = &EncodeError{s.Name, errors.New(r)}
			case error:
				err = &EncodeError{s.Name, r}
			}
		}
	}()

	buf := ibuf.From(*dataPtr)

	switch d := d.(type) {
	case map[string]interface{}:
		err = s.encodeDynamic(buf, s.rootOp, d)

	case *map[string]interface{}:
		err = s.encodeDynamic(buf, s.rootOp, *d)

	case []map[string]interface{}:
		err = s.encodeDynamic(buf, s.rootOp, d)

	case []interface{}:
		err = s.encodeDynamic(buf, s.rootOp, d)

	case *interface{}:
		return s.EncodeTo(*d, dataPtr)

	default:
		dType := reflect2.TypeOf(d)
		switch dType.Kind() {
		case reflect.Struct:
			var sop *structOperation
			sop, err = s.getStructOperation(dType, d)
			if err != nil {
				err = &EncodeError{s.Name, err}
				return
			}
			err = s.encodeStruct(buf, sop, d)

		case reflect.Slice:
			sliceType := dType.(*reflect2.UnsafeSliceType)
			dType = sliceType.Elem()
			if dType.Kind() == reflect.Struct {
				var sop *structOperation
				sop, err = s.getStructOperation(dType, d)
				if err != nil {
					err = &EncodeError{s.Name, err}
					return
				}
				err = s.encodeStruct(buf, sop, d)
			} else if dType.Kind() == reflect.Map {
				err = s.encodeDynamic(buf, s.rootOp, d)
			} else {
				err = &EncodeError{s.Name, &WrongTypeError{dType.String()}}
				return
			}

		case reflect.Array:
			sliceType := dType.(*reflect2.UnsafeArrayType)
			dType = sliceType.Elem()
			if dType.Kind() == reflect.Struct {
				var sop *structOperation
				sop, err = s.getStructOperation(dType, d)
				if err != nil {
					err = &EncodeError{s.Name, err}
					return
				}
				err = s.encodeStruct(buf, sop, d)
			} else if dType.Kind() == reflect.Map {
				err = s.encodeDynamic(buf, s.rootOp, d)
			} else {
				err = &EncodeError{s.Name, &WrongTypeError{dType.String()}}
				return
			}

		case reflect.Ptr:
			elemType := toPtrElemType(dType)
			return s.EncodeTo(elemType.Indirect(d), dataPtr)
		default:
			err = &EncodeError{s.Name, &WrongTypeError{dType.String()}}
			return
		}
	}

	if err != nil {
		err = &EncodeError{s.Name, err}
		return
	}

	// enlarge default encoder buffer allocation with latest data
	s.encodeBufSize = maxInt64(s.encodeBufSize, buf.Offset())

	*dataPtr = buf.Seal()
	return
}

func (s *Schema) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, d interface{}) error {
	if opNode.handler == nil {
		return fmt.Errorf("opearation handler is nil")
	}
	return opNode.handler.encodeStruct(buf, opNode, reflect2.PtrOf(d))
}

func (s *Schema) encodeDynamic(buf *ibuf.Buffer, opNode *operation, d interface{}) error {
	if opNode.handler == nil {
		return fmt.Errorf("opearation handler is nil")
	}
	return opNode.handler.encodeDynamic(buf, opNode, d)
}
