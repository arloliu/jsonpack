package jsonpack

import (
	"reflect"

	"github.com/pkg/errors"

	ibuf "github.com/arloliu/jsonpack/buffer"

	"github.com/modern-go/reflect2"
)

/*
Encode returns packed binary data of v with compiled schema definition.
The compiled schema definition is specified by schemaName and
return data will be nil if error occurs.

The type of v can be a map[string]interface{} which represents valid JSON data, or
a struct instance which added by AddSchema function.

The return data contains encoded binary data that can then be decoded by Decode function.

Before encode data into jsonpack encoding format, we need to create a jsonpack instance and create a schema.
	// the Info struct we want to decode
	type Info struct {
		Name string `json:"name"`
		Area uint32 `json:"area"`
		// omit this field
		ExcludeField string `-`
	}
	// create a new jsonpack instance
	jsonPack := jsonpack.NewJSONPack()
	// create schema with Info struct
	sch, err := jsonPack.AddSchema("Info", Info{}, jsonpack.LittleEndian)

If we want to encode a map with Info schema.
	data := map[string]interface{} {
		"name": "example name",
		"area": uint32(888),
	}

	// encodes data into encodedResult
	encodedResult, err := sch.Encode(data)

Or if we want to encode a Info struct instance with Info schema.
	data := &Info{
		Name: "example name",
		Area: 888,
	}

	// encodes data into encodedResult
	encodedResult, err := sch.Encode(infoStruct)
*/
func (s *Schema) Encode(d interface{}) ([]byte, error) {
	encodeData := make([]byte, s.encodeBufSize)
	err := s.EncodeTo(d, &encodeData)
	return encodeData, err
}

// Marshal is an alias to Encode function, provides familiar interface of json package
func (s *Schema) Marshal(d interface{}) ([]byte, error) {
	return s.Encode(d)
}

// EncodeTo is similar to Encode function, but passing a pointer to []byte to store
// encoded data instead of returning new allocated []byte encoded data.
//
// This method is useful with buffer pool for saving memory allocation usage and improving performance.
//
// Caution: the encoder might re-allocate and grow the slice if necessary, the length and capacity of slice might be changed.
func (s *Schema) EncodeTo(d interface{}, dataPtr *[]byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch r := r.(type) {
			case string:
				err = errors.WithStack(&EncodeError{s.Name, errors.New(r)})
			case error:
				err = errors.WithStack(&EncodeError{s.Name, r})
			}
		}
	}()

	buf := ibuf.From(*dataPtr)

	switch d := d.(type) {
	// fast path: use type assertion, it's faster then reflection
	case map[string]interface{}:
		err = s.encodeDynamic(buf, s.rootOp, d)

	case *map[string]interface{}:
		err = s.encodeDynamic(buf, s.rootOp, *d)

	case []map[string]interface{}:
		err = s.encodeDynamic(buf, s.rootOp, d)

	case *[]map[string]interface{}:
		err = s.encodeDynamic(buf, s.rootOp, *d)

	case []interface{}:
		err = s.encodeDynamic(buf, s.rootOp, d)

	case *[]interface{}:
		err = s.encodeDynamic(buf, s.rootOp, *d)

	case *interface{}:
		return s.EncodeTo(*d, dataPtr)

	default:
		// slow path: use reflection to check type
		dType := reflect2.TypeOf(d)
		switch dType.Kind() {
		case reflect.Struct:
			var sop *structOperation
			sop, err = s.getStructOperation(dType, d)
			if err != nil {
				err = errors.WithStack(&EncodeError{s.Name, err})
				return
			}
			err = s.encodeStruct(buf, sop, d)

		case reflect.Slice:
			sliceType := dType.(*reflect2.UnsafeSliceType)
			dType = sliceType.Elem()
			switch dType.Kind() {
			case reflect.Struct:
				var sop *structOperation
				sop, err = s.getStructOperation(dType, d)
				if err != nil {
					err = errors.WithStack(&EncodeError{s.Name, err})
					return
				}
				err = s.encodeStruct(buf, sop, d)

			case reflect.Map, reflect.Interface:
				err = s.encodeDynamic(buf, s.rootOp, d)

			default:
				return errors.WithStack(&EncodeError{s.Name, &WrongTypeError{dType.String()}})
			}

		case reflect.Array:
			arrayType := dType.(*reflect2.UnsafeArrayType)
			dType = arrayType.Elem()
			switch dType.Kind() {
			case reflect.Struct:
				var sop *structOperation
				sop, err = s.getStructOperation(dType, d)
				if err != nil {
					err = errors.WithStack(&EncodeError{s.Name, err})
					return
				}
				err = s.encodeStruct(buf, sop, d)

			case reflect.Map, reflect.Interface:
				err = s.encodeDynamic(buf, s.rootOp, d)

			default:
				return errors.WithStack(&EncodeError{s.Name, &WrongTypeError{dType.String()}})
			}

		case reflect.Ptr:
			elemType := toPtrElemType(dType)
			return s.EncodeTo(elemType.Indirect(d), dataPtr)
		default:
			err = errors.WithStack(&EncodeError{s.Name, &WrongTypeError{dType.String()}})
			return
		}
	}

	if err != nil {
		err = errors.WithStack(&EncodeError{s.Name, err})
		return
	}

	// enlarge default encoder buffer allocation with latest encoded result
	s.encodeBufSize = maxInt64(s.encodeBufSize, buf.Offset())

	*dataPtr = buf.Seal()
	return
}

func (s *Schema) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, d interface{}) error {
	if opNode.handler == nil {
		return errors.Errorf("opearation handler is nil")
	}
	return opNode.handler.encodeStruct(buf, opNode, reflect2.PtrOf(d))
}

func (s *Schema) encodeDynamic(buf *ibuf.Buffer, opNode *operation, d interface{}) error {
	if opNode.handler == nil {
		return errors.Errorf("opearation handler is nil")
	}
	return opNode.handler.encodeDynamic(buf, opNode, d)
}
