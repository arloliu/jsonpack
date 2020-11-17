package jsonpack

import (
	"fmt"
	"reflect"
	"unsafe"

	ibuf "github.com/arloliu/jsonpack/buffer"

	"github.com/modern-go/reflect2"
)

type opHandler interface {
	encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error
	decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error
	// write(buf *ibuf.Buffer, data interface{}) error
	encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error
	decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error)
}

type opHandlerType uint8

const (
	nullOpType opHandlerType = iota
	objectOpType
	structOpType
	sliceOpType
	arrayOpType
	booleanOpType
	stringOpType
	int8OpType
	int16LEOpType
	int16BEOpType
	int32LEOpType
	int32BEOpType
	int64LEOpType
	int64BEOpType
	uint8OpType
	uint16LEOpType
	uint16BEOpType
	uint32LEOpType
	uint32BEOpType
	uint64LEOpType
	uint64BEOpType
	float32LEOpType
	float32BEOpType
	float64LEOpType
	float64BEOpType
)

// Null type
type nullOp struct{}

func (p *nullOp) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	return &NotImplementedError{"nullOp.decodeStruct"}
}
func (p *nullOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	return &NotImplementedError{"nullOp.encodeStruct"}
}
func (p *nullOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	return &NotImplementedError{"nullOp.encodeDynamic"}
}
func (p *nullOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return nil, &NotImplementedError{"nullOp.decodeDynamic"}
}

// declare operation variables for type alias
var (
	_nullOp   = &nullOp{}
	_structOp = &structOp{}
	// _objectOp    = &objectOp{}
	_sliceOp     = &sliceOp{}
	_arrayOp     = &arrayOp{}
	_booleanOp   = &booleanOp{}
	_float32LEOp = &float32LEOp{}
	_float32BEOp = &float32BEOp{}
	_float64LEOp = &float64LEOp{}
	_float64BEOp = &float64BEOp{}
)

var builtinTypes = map[string]opHandler{
	"bool":      _booleanOp,
	"boolean":   _booleanOp,
	"string":    &stringOp{},
	"int8":      &int8Op{},
	"int16be":   &int16BEOp{},
	"int16le":   &int16LEOp{},
	"int32be":   &int32BEOp{},
	"int32le":   &int32LEOp{},
	"int64be":   &int64BEOp{},
	"int64le":   &int64LEOp{},
	"uint8":     &uint8Op{},
	"uint16be":  &uint16BEOp{},
	"uint16le":  &uint16LEOp{},
	"uint32be":  &uint32BEOp{},
	"uint32le":  &uint32LEOp{},
	"uint64be":  &uint64BEOp{},
	"uint64le":  &uint64LEOp{},
	"floatbe":   _float32BEOp,
	"floatle":   _float32LEOp,
	"float32be": _float32BEOp,
	"float32le": _float32LEOp,
	"doublebe":  _float64BEOp,
	"doublele":  _float64LEOp,
	"float64be": _float64BEOp,
	"float64le": _float64LEOp,
}

var builtinOpHandlerTypes = map[string]opHandlerType{
	"bool":      booleanOpType,
	"boolean":   booleanOpType,
	"string":    stringOpType,
	"int8":      int8OpType,
	"int16be":   int16BEOpType,
	"int16le":   int16LEOpType,
	"int32be":   int32BEOpType,
	"int32le":   int32LEOpType,
	"int64be":   int64BEOpType,
	"int64le":   int64LEOpType,
	"uint8":     uint8OpType,
	"uint16be":  uint16BEOpType,
	"uint16le":  uint16LEOpType,
	"uint32be":  uint32BEOpType,
	"uint32le":  uint32LEOpType,
	"uint64be":  uint64BEOpType,
	"uint64le":  uint64LEOpType,
	"floatbe":   float32BEOpType,
	"floatle":   float32LEOpType,
	"float32be": float32BEOpType,
	"float32le": float32LEOpType,
	"doublebe":  float64BEOpType,
	"doublele":  float64LEOpType,
	"float64be": float64BEOpType,
	"float64le": float64LEOpType,
}

func isBuiltinType(propType *string) bool {
	_, ok := builtinTypes[*propType]
	return ok
}

func getPropData(opNode *operation, d interface{}) interface{} {
	if len(opNode.propName) == 0 {
		return d
	}
	switch d := d.(type) {
	case map[string]interface{}:
		fVal, ok := d[opNode.propName]
		if !ok {
			return nil
		}
		return fVal

	case *map[string]interface{}:
		return getPropData(opNode, *d)

	case *interface{}:
		return getPropData(opNode, *d)

	default:
		dType := reflect2.TypeOf(d)
		switch dType.Kind() {
		case reflect.Ptr:
			ptrType := dType.(*reflect2.UnsafePtrType)
			return getPropData(opNode, ptrType.Elem().Indirect(d))
		default:
		}
	}

	return nil
}

func derefPtr(ptr unsafe.Pointer) unsafe.Pointer {
	return *(*unsafe.Pointer)(ptr)
}

func assignPtr(dst, src unsafe.Pointer) {
	if dst != src {
		*(*unsafe.Pointer)(dst) = src
	}
}

func encodeSliceTypeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	var err error

	switch items := data.(type) {
	// slice type fast path
	case []interface{}:
		// logger.Tracef("arrayOp.encodeDynamic FASTPATH, type: %T, data: %+v", data, data)
		itemLen := len(items)
		itemOpNode := opNode.children[0]
		buf.WriteVarUint((uint64(itemLen)))
		for i := 0; i < itemLen; i++ {
			err = itemOpNode.handler.encodeDynamic(buf, itemOpNode, items[i])
			if err != nil {
				return err
			}
		}
	case []map[string]interface{}:
		// logger.Tracef("arrayOp.encodeDynamic FASTPATH, type: %T, data: %+v", data, data)
		itemLen := len(items)
		itemOpNode := opNode.children[0]
		buf.WriteVarUint((uint64(itemLen)))
		for i := 0; i < itemLen; i++ {
			err = itemOpNode.handler.encodeDynamic(buf, itemOpNode, items[i])
			if err != nil {
				return err
			}
		}

	default:
		// logger.Tracef("arrayOp.encodeDynamic SLOWPATH, type: %T, data: %+v", data, data)
		fVal := reflect.ValueOf(data)
		fKind := fVal.Type().Kind()
		if fKind == reflect.Slice || fKind == reflect.Array {
			itemLen := fVal.Len()
			itemOpNode := opNode.children[0]
			buf.WriteVarUint((uint64(itemLen)))
			for i := 0; i < itemLen; i++ {
				err = itemOpNode.handler.encodeDynamic(buf, itemOpNode, fVal.Index(i).Interface())
				if err != nil {
					return err
				}
			}
		} else {
			return &WrongTypeError{fmt.Sprintf("%T", data)}
		}
	}

	return nil
}

func decodeSliceAnyDynamic(buf *ibuf.Buffer, opNode *operation, ptr *[]interface{}) (interface{}, error) {
	var err error
	var m []interface{} = *ptr
	size, _ := buf.ReadVarUint()

	if cap(m) < int(size) {
		m = make([]interface{}, size)
	} else {
		m = m[:size]
	}
	childNode := opNode.children[0]
	for i := uint64(0); i < size; i++ {
		childProp := _createNewData(buf, childNode)
		m[i], err = childNode.handler.decodeDynamic(buf, childNode, childProp)
		if err != nil {
			return nil, err
		}
	}

	*ptr = m

	return m, nil
}

func decodeSliceMapDynamic(buf *ibuf.Buffer, opNode *operation, ptr *[]map[string]interface{}) (interface{}, error) {
	var m []map[string]interface{} = *ptr
	var ok bool

	size, _ := buf.ReadVarUint()

	if cap(m) < int(size) {
		m = make([]map[string]interface{}, size)
	} else {
		m = m[:size]
	}
	childNode := opNode.children[0]
	for i := uint64(0); i < size; i++ {
		childProp := _createNewData(buf, childNode)
		result, err := childNode.handler.decodeDynamic(buf, childNode, childProp)
		if err != nil {
			return nil, err
		}
		m[i], ok = result.(map[string]interface{})
		if !ok {
			return nil, &TypeAssertionError{Data: result, ExpectedType: "map[string]interface{}"}
		}
	}

	*ptr = m

	return m, nil
}

func decodeSliceTypeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	switch ptr := v.(type) {
	case *[]interface{}:
		return decodeSliceAnyDynamic(buf, opNode, ptr)
	case []interface{}:
		return decodeSliceAnyDynamic(buf, opNode, &ptr)
	case *[]map[string]interface{}:
		return decodeSliceMapDynamic(buf, opNode, ptr)
	case []map[string]interface{}:
		return decodeSliceMapDynamic(buf, opNode, &ptr)
	}
	return nil, &WrongTypeError{reflect.TypeOf(v).Name()}
}

func decodeArrayTypeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	var err error
	var m []interface{}
	var ptr *[]interface{}

	size, _ := buf.ReadVarUint()

	val := reflect.ValueOf(v).Elem()
	if val.Len() < int(size) {
		return nil, fmt.Errorf("the capacity of array of decode target less than data")
	}

	m = val.Slice(0, int(size)).Interface().([]interface{})
	ptr = &m
	childNode := opNode.children[0]
	for i := 0; i < int(size); i++ {
		childProp := _createNewData(buf, childNode)
		m[i], err = childNode.handler.decodeDynamic(buf, childNode, childProp)
		if err != nil {
			return nil, err
		}
	}
	*ptr = m

	return v, nil
}

func _createNewData(buf *ibuf.Buffer, opNode *operation) interface{} {
	switch opNode.handlerType {
	case objectOpType:
		return make(map[string]interface{})
	case arrayOpType, sliceOpType:
		size, n := buf.ReadVarUint()
		buf.Seek(-int64(n), true)
		return make([]interface{}, size)
	default:
		return nil
	}
}
