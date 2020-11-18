package jsonpack

import (
	"reflect"

	"github.com/pkg/errors"

	ibuf "github.com/arloliu/jsonpack/buffer"

	"github.com/modern-go/reflect2"
)

/*
Decode reads encoded data with compiled schema definition and stores the result
in the value pointed to v.

If type of v is not a pointer type that pointed to a map or struct, Decode
function will return DecodeError.

The valid type of v is either a *map[string]interface{} or a pointer to the struct
which added by AddSchema function.
*/
func (s *Schema) Decode(data []byte, v interface{}) (err error) {
	return s.decode(data, v, true)
}

// Unmarshal is an alias to Decode function, provides familiar interface of json package
func (s *Schema) Unmarshal(data []byte, v interface{}) (err error) {
	return s.Decode(data, v)
}

func (s *Schema) decode(data []byte, v interface{}, checkPtrType bool) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch r := r.(type) {
			case string:
				err = errors.WithStack(&DecodeError{s.Name, errors.New(r)})
			case error:
				err = errors.WithStack(&DecodeError{s.Name, r})
			}
		}
	}()

	if v == nil {
		return errors.WithStack(&DecodeError{s.Name, errors.New("target of decoding is nil")})
	}

	vType := reflect2.TypeOf(v)
	if vType == nil {
		return errors.WithStack(&DecodeError{s.Name, errors.New("invalid type of target")})
	}

	vKind := vType.Kind()

	if reflect2.IsNil(v) {
		err = errors.WithStack(&DecodeError{s.Name, &WrongTypeError{vType.String()}})
		return
	}

	if checkPtrType && vKind != reflect.Ptr && vKind != reflect.Map {
		err = errors.WithStack(&DecodeError{s.Name, &WrongTypeError{vType.String()}})
		return
	}

	buf := ibuf.From(data)
	switch d := v.(type) {
	case *map[string]interface{}:
		_, err = s.decodeDynamic(buf, s.rootOp, *d)
		if err != nil {
			return err
		}

	case *[]map[string]interface{}:
		_, err = _sliceOp.decodeDynamic(buf, s.rootOp, d)
		if err != nil {
			return err
		}

	case *[]interface{}:
		// pass pointer of slice instead of dereferenced value due to slice is not reference type
		_, err = _sliceOp.decodeDynamic(buf, s.rootOp, d)
		if err != nil {
			return err
		}

	case *interface{}:
		return s.decode(data, *d, false)

	default:
		switch vKind {
		case reflect.Struct:
			var sop *structOperation
			sop, err = s.getStructOperation(vType, v)
			if err != nil {
				return errors.WithStack(&DecodeError{s.Name, err})
			}
			return s.decodeStruct(buf, sop, v)

		case reflect.Slice:
			sliceType := vType.(*reflect2.UnsafeSliceType)
			vType = sliceType.Elem()
			if vType.Kind() == reflect.Struct {
				var sop *structOperation
				sop, err = s.getStructOperation(vType, v)
				if err != nil {
					return errors.WithStack(&DecodeError{s.Name, err})
				}
				return s.decodeStruct(buf, sop, v)
			} else if vType.Kind() == reflect.Map || vType.Kind() == reflect.Interface {
				_, err = s.decodeDynamic(buf, s.rootOp, v)
			} else {
				return errors.WithStack(&DecodeError{s.Name, &WrongTypeError{vType.String()}})
			}

		case reflect.Array:
			sliceType := vType.(*reflect2.UnsafeArrayType)
			vType = sliceType.Elem()
			if vType.Kind() == reflect.Struct {
				var sop *structOperation
				sop, err = s.getStructOperation(vType, v)
				if err != nil {
					return errors.WithStack(&DecodeError{s.Name, err})
				}
				return s.decodeStruct(buf, sop, v)
			} else if vType.Kind() == reflect.Map || vType.Kind() == reflect.Interface {
				_, err = s.decodeDynamic(buf, s.rootOp, v)
				if err != nil {
					return err
				}
			} else {
				return errors.WithStack(&DecodeError{s.Name, &WrongTypeError{vType.String()}})
			}

		case reflect.Ptr:
			elemType := toPtrElemType(vType)
			if elemType.Kind() == reflect.Array {
				_, err = _arrayOp.decodeDynamic(buf, s.rootOp, v)
			} else {
				return s.decode(data, elemType.Indirect(v), false)
			}

		default:
			return errors.WithStack(&DecodeError{s.Name, &WrongTypeError{vType.String()}})
		}
	}
	if err != nil {
		err = errors.WithStack(&DecodeError{s.Name, err})
	}
	return
}

func (s *Schema) decodeDynamic(buf *ibuf.Buffer, opNode *operation, d interface{}) (interface{}, error) {
	return opNode.handler.decodeDynamic(buf, opNode, d)
}

func (s *Schema) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, v interface{}) error {
	if opNode.handler == nil {
		return errors.Errorf("opearation handler is nil")
	}
	return opNode.handler.decodeStruct(buf, opNode, reflect2.PtrOf(v))
}
