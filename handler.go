package jsonpack

import (
	"fmt"
	"reflect"
	"unsafe"

	ibuf "github.com/arloliu/jsonpack/internal/buffer"

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

type structOp struct{}

func (p *structOp) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	var err error
	var filedPtr unsafe.Pointer

	for _, op := range opNode.children {
		filedPtr = op.field.UnsafeGet(ptr)

		if op.isPtrType {
			ptrType := (op.field.Type()).(*reflect2.UnsafePtrType)
			elemType := ptrType.Elem()

			// check if current pointer of field is null
			if *((*unsafe.Pointer)(filedPtr)) == nil {
				//pointer to null, need to allocate memory to hold the value
				newPtr := elemType.UnsafeNew()
				err = op.handler.decodeStruct(buf, op, newPtr)
				if err != nil {
					return err
				}
				// assign new allocated data back to current field pointer
				*((*unsafe.Pointer)(filedPtr)) = newPtr
			} else {
				//reuse existing instance
				err = op.handler.decodeStruct(buf, op, *((*unsafe.Pointer)(filedPtr)))
				if err != nil {
					return err
				}
			}
		} else {
			err = op.handler.decodeStruct(buf, op, filedPtr)
			if err != nil {
				return err
			}
		}
	}

	return nil

}
func (p *structOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	var err error
	var filedPtr unsafe.Pointer

	for _, op := range opNode.children {
		filedPtr = op.field.UnsafeGet(ptr)
		// dereference pointer
		if op.isPtrType {
			// logger.Debugf("encodeStruct: op is pointer type, field name: %v", op.field.Name())
			filedPtr = *((*unsafe.Pointer)(filedPtr))
		}

		err = op.handler.encodeStruct(buf, op, filedPtr)
		if err != nil {
			return err
		}
	}
	return nil
}
func (p *structOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	return &NotImplementedError{"structOp.encodeDynamic"}
}
func (p *structOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return nil, &NotImplementedError{"structOp.decodeDynamic"}
}

type objectOp struct {
}

// decode dynamic map type in struct
func (p *objectOp) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	opInstance := opNode.opInstance

	mapType := opNode.opType.(*reflect2.UnsafeMapType)
	// make a new map if current map hasn't allocated
	if mapType.UnsafeIsNil(ptr) {
		mapType.UnsafeSet(ptr, mapType.UnsafeMakeMap(0))
	}

	_, err := p.decodeDynamic(buf, opInstance, mapType.UnsafeIndirect(ptr))
	return err
}

// encode dynamic map type in struct
func (p *objectOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	// schema := (opNode.handler.(*objectOp)).schema
	// opInstance := opNode.opInstance

	return p.encodeDynamic(buf, opNode.opInstance, opNode.opType.UnsafeIndirect(ptr))
}

func (p *objectOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	// logger.Debugf("objectOp.encodeDynamic type: %T, propName: %s data: %+v", opNode.handler, opNode.propName, data)

	// encode properties in map
	for _, childNode := range opNode.children {
		childData := getPropData(childNode, data)
		// logger.Debugf("objectOp.encodeDynamic, CHILD type: %T, propName: %s data: %+v", childNode.handler, childNode.propName, childData)
		err := childNode.handler.encodeDynamic(buf, childNode, childData)
		if err != nil {
			return err
		}
	}
	return nil
}
func (p *objectOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	var err error

	for _, childNode := range opNode.children {
		m, _ := v.(map[string]interface{})
		childProp := _createNewData(buf, childNode)
		m[childNode.propName], err = childNode.handler.decodeDynamic(buf, childNode, childProp)
		if err != nil {
			return nil, err
		}
	}
	return v, nil

}

type arrayOp struct{}

func (p *arrayOp) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	var err error

	arrayType := opNode.opType.(*reflect2.UnsafeArrayType)
	itemOp := opNode.children[0]

	length, _ := buf.ReadVarUint()

	for i := 0; i < int(length); i++ {
		itemPtr := arrayType.UnsafeGetIndex(ptr, i)
		err = itemOp.handler.decodeStruct(buf, itemOp, itemPtr)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *arrayOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	var err error

	arrayType := opNode.opType.(*reflect2.UnsafeArrayType)
	length := arrayType.Len()
	itemOp := opNode.children[0]

	buf.WriteVarUint(uint64(length))
	for i := 0; i < length; i++ {
		itemPtr := arrayType.UnsafeGetIndex(ptr, i)
		err = itemOp.handler.encodeStruct(buf, itemOp, itemPtr)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *arrayOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	return _encodeSliceTypeDynamic(buf, opNode, data)
}

func (p *arrayOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return _decodeSliceTypeDynamic(buf, opNode, v)
}

type sliceOp struct{}

func (p *sliceOp) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	var err error

	sliceType := opNode.opType.(*reflect2.UnsafeSliceType)
	itemOp := opNode.children[0]

	l, _ := buf.ReadVarUint()
	length := int(l)

	if sliceType.UnsafeIsNil(ptr) {
		// allocate new slice
		sliceType.UnsafeSet(ptr, sliceType.UnsafeMakeSlice(length, length))
	} else if sliceType.UnsafeLengthOf(ptr) < length {
		// grow slice
		sliceType.UnsafeGrow(ptr, length)
	}

	for i := 0; i < length; i++ {
		itemPtr := sliceType.UnsafeGetIndex(ptr, i)

		var newPtr unsafe.Pointer = itemPtr
		if itemOp.isPtrType {
			// allocate new struct instance
			if *((*unsafe.Pointer)(itemPtr)) == nil {
				ptrType := (itemOp.opType).(*reflect2.UnsafePtrType)
				elemType := ptrType.Elem()
				newPtr = elemType.UnsafeNew()
			}
		}
		err = itemOp.handler.decodeStruct(buf, itemOp, newPtr)
		if err != nil {
			return err
		}

		if itemPtr != newPtr {
			*((*unsafe.Pointer)(itemPtr)) = newPtr
		}
	}
	return nil
}
func (p *sliceOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	var err error

	sliceType := opNode.opType.(*reflect2.UnsafeSliceType)
	length := sliceType.UnsafeLengthOf(ptr)
	itemOp := opNode.children[0]

	buf.WriteVarUint(uint64(length))
	for i := 0; i < length; i++ {
		itemPtr := sliceType.UnsafeGetIndex(ptr, i)
		err = itemOp.handler.encodeStruct(buf, itemOp, itemPtr)
		if err != nil {
			return err
		}
	}
	return nil
}
func (p *sliceOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	return _encodeSliceTypeDynamic(buf, opNode, data)
}

func (p *sliceOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return _decodeSliceTypeDynamic(buf, opNode, v)
}

type stringOp struct{}

func (p *stringOp) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	str := buf.ReadString()
	*((*string)(ptr)) = str
	// *((*string)(ptr)) = buf.ReadString()
	return nil
}
func (p *stringOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	d := *((*string)(ptr))
	buf.WriteString(&d)

	return nil
}
func (p *stringOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	d, ok := data.(string)
	if !ok {
		return &TypeAssertionError{data, "string"}
	}
	buf.WriteString(&d)

	return nil
}
func (p *stringOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return buf.ReadString(), nil
}

type booleanOp struct{}

func (p *booleanOp) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	d := buf.ReadByte()
	*((*bool)(ptr)) = (d != 0)
	return nil
}
func (p *booleanOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	d := *((*bool)(ptr))
	if d {
		buf.WriteByte(byte(1))
	} else {
		buf.WriteByte(byte(0))
	}
	return nil
}
func (p *booleanOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	d := data.(bool)
	if d {
		buf.WriteByte(byte(1))
	} else {
		buf.WriteByte(byte(0))
	}
	return nil
}
func (p *booleanOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	d := buf.ReadByte()
	return d != 0, nil
}

/*
 * Signed interger handlers
 */
type int8Op struct{}

func (p *int8Op) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	*((*int8)(ptr)) = buf.ReadInt8()
	return nil
}

func (p *int8Op) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	val := *((*byte)(ptr))
	buf.WriteByte(val)
	return nil
}
func (p *int8Op) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	var val byte
	switch data := data.(type) {
	case float64:
		val = byte(data)
	case int8:
		val = byte(data)
	default:
		return &TypeAssertionError{data, "int8"}
	}

	buf.WriteByte(byte(val))
	return nil
}
func (p *int8Op) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return buf.ReadInt8(), nil
}

type int16LEOp struct{}

func (p *int16LEOp) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	*((*int16)(ptr)) = buf.ReadInt16LE()
	return nil
}

func (p *int16LEOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	val := *((*int16)(ptr))
	buf.WriteInt16LE(val)
	return nil

}
func (p *int16LEOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	var val int16
	switch data := data.(type) {
	case int16:
		val = data
	case float64:
		val = int16(data)
	default:
		return &TypeAssertionError{data, "int16"}
	}
	buf.WriteInt16LE(val)
	return nil
}
func (p *int16LEOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return buf.ReadInt16LE(), nil
}

type int16BEOp struct{}

func (p *int16BEOp) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	*((*int16)(ptr)) = buf.ReadInt16BE()
	return nil
}

func (p *int16BEOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	val := *((*int16)(ptr))
	buf.WriteInt16BE(val)
	return nil
}
func (p *int16BEOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	var val int16
	switch data := data.(type) {
	case int16:
		val = data
	case float64:
		val = int16(data)
	default:
		return &TypeAssertionError{data, "int16"}
	}
	buf.WriteInt16BE(val)
	return nil
}
func (p *int16BEOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return buf.ReadInt16BE(), nil
}

type int32LEOp struct{}

func (p *int32LEOp) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	*((*int32)(ptr)) = buf.ReadInt32LE()
	return nil
}

func (p *int32LEOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	val := *((*int32)(ptr))
	buf.WriteInt32LE(val)
	return nil
}
func (p *int32LEOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	var val int32
	switch data := data.(type) {
	case int32:
		val = data
	case float64:
		val = int32(data)
	default:
		return &TypeAssertionError{data, "int32"}
	}
	buf.WriteInt32LE(val)
	return nil
}
func (p *int32LEOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return buf.ReadInt32LE(), nil
}

type int32BEOp struct{}

func (p *int32BEOp) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	*((*int32)(ptr)) = buf.ReadInt32BE()
	return nil
}

func (p *int32BEOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	val := *((*int32)(ptr))
	buf.WriteInt32BE(val)
	return nil
}
func (p *int32BEOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	var val int32
	switch data := data.(type) {
	case int32:
		val = data
	case float64:
		val = int32(data)
	default:
		return &TypeAssertionError{data, "int32"}
	}
	buf.WriteInt32BE(val)
	return nil
}
func (p *int32BEOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return buf.ReadInt32BE(), nil
}

type int64LEOp struct{}

func (p *int64LEOp) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	*((*int64)(ptr)) = buf.ReadInt64LE()
	return nil
}
func (p *int64LEOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	val := *((*int64)(ptr))
	buf.WriteInt64LE(val)
	return nil
}
func (p *int64LEOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	var val int64
	switch data := data.(type) {
	case int64:
		val = data
	case float64:
		val = int64(data)
	default:
		return &TypeAssertionError{data, "int64"}
	}
	buf.WriteInt64LE(val)
	return nil
}
func (p *int64LEOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return buf.ReadInt64LE(), nil
}

type int64BEOp struct{}

func (p *int64BEOp) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	*((*int64)(ptr)) = buf.ReadInt64BE()
	return nil
}
func (p *int64BEOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	val := *((*int64)(ptr))
	buf.WriteInt64BE(val)
	return nil
}
func (p *int64BEOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	var val int64
	switch data := data.(type) {
	case int64:
		val = data
	case float64:
		val = int64(data)
	default:
		return &TypeAssertionError{data, "int64"}
	}
	buf.WriteInt64BE(val)
	return nil
}
func (p *int64BEOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return buf.ReadInt64BE(), nil
}

/*
 * Unsigned interger handlers
 */
type uint8Op struct{}

func (p *uint8Op) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	*((*uint8)(ptr)) = buf.ReadUint8()
	return nil
}
func (p *uint8Op) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	val := *((*byte)(ptr))
	buf.WriteByte(val)
	return nil
}
func (p *uint8Op) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	var val byte
	switch data := data.(type) {
	case uint8:
		val = byte(data)
	case float64:
		val = byte(data)
	default:
		return &TypeAssertionError{data, "uint8"}
	}
	buf.WriteByte(val)
	return nil
}
func (p *uint8Op) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return buf.ReadUint8(), nil
}

type uint16LEOp struct{}

func (p *uint16LEOp) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	*((*uint16)(ptr)) = buf.ReadUint16LE()
	return nil
}
func (p *uint16LEOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	val := *((*uint16)(ptr))
	buf.WriteUint16LE(val)
	return nil
}
func (p *uint16LEOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	var val uint16
	switch data := data.(type) {
	case uint16:
		val = uint16(data)
	case float64:
		val = uint16(data)
	default:
		return &TypeAssertionError{data, "uint16"}
	}
	buf.WriteUint16LE(val)
	return nil
}
func (p *uint16LEOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return buf.ReadUint16LE(), nil
}

type uint16BEOp struct{}

func (p *uint16BEOp) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	*((*uint16)(ptr)) = buf.ReadUint16BE()
	return nil
}
func (p *uint16BEOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	val := *((*uint16)(ptr))
	buf.WriteUint16BE(val)
	return nil
}
func (p *uint16BEOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	var val uint16
	switch data := data.(type) {
	case uint16:
		val = uint16(data)
	case float64:
		val = uint16(data)
	default:
		return &TypeAssertionError{data, "uint16"}
	}
	buf.WriteUint16BE(val)
	return nil
}
func (p *uint16BEOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return buf.ReadUint16BE(), nil
}

type uint32LEOp struct{}

func (p *uint32LEOp) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	*((*uint32)(ptr)) = buf.ReadUint32LE()
	return nil
}
func (p *uint32LEOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	val := *((*uint32)(ptr))
	buf.WriteUint32LE(val)
	return nil
}
func (p *uint32LEOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	var val uint32
	switch data := data.(type) {
	case uint32:
		val = uint32(data)
	case float64:
		val = uint32(data)
	default:
		return &TypeAssertionError{data, "uint32"}
	}
	buf.WriteUint32LE(val)
	return nil
}
func (p *uint32LEOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return buf.ReadUint32LE(), nil
}

type uint32BEOp struct{}

func (p *uint32BEOp) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	*((*uint32)(ptr)) = buf.ReadUint32BE()
	return nil
}

func (p *uint32BEOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	val := *((*uint32)(ptr))
	buf.WriteUint32BE(val)
	return nil
}
func (p *uint32BEOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	var val uint32
	switch data := data.(type) {
	case uint32:
		val = uint32(data)
	case float64:
		val = uint32(data)
	default:
		return &TypeAssertionError{data, "uint32"}
	}
	buf.WriteUint32BE(val)
	return nil
}
func (p *uint32BEOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return buf.ReadUint32BE(), nil
}

type uint64LEOp struct{}

func (p *uint64LEOp) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	*((*uint64)(ptr)) = buf.ReadUint64LE()
	return nil
}

func (p *uint64LEOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	val := *((*uint64)(ptr))
	buf.WriteUint64LE(val)
	return nil
}
func (p *uint64LEOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	var val uint64
	switch data := data.(type) {
	case uint64:
		val = uint64(data)
	case float64:
		val = uint64(data)
	default:
		return &TypeAssertionError{data, "uint64"}
	}
	buf.WriteUint64LE(val)
	return nil
}
func (p *uint64LEOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return buf.ReadUint64LE(), nil
}

type uint64BEOp struct{}

func (p *uint64BEOp) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	*((*uint64)(ptr)) = buf.ReadUint64BE()
	return nil
}
func (p *uint64BEOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	val := *((*uint64)(ptr))
	buf.WriteUint64BE(val)
	return nil
}
func (p *uint64BEOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	var val uint64
	switch data := data.(type) {
	case uint64:
		val = uint64(data)
	case float64:
		val = uint64(data)
	default:
		return &TypeAssertionError{data, "uint64"}
	}
	buf.WriteUint64BE(val)
	return nil
}
func (p *uint64BEOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return buf.ReadUint64BE(), nil
}

/*
 * floating number handlers
 */
type float32LEOp struct{}

func (p *float32LEOp) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	*((*float32)(ptr)) = buf.ReadFloat32LE()
	return nil
}
func (p *float32LEOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	val := *((*float32)(ptr))
	buf.WriteFloat32LE(val)
	return nil
}
func (p *float32LEOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	var val float32
	switch data := data.(type) {
	case float32:
		val = float32(data)
	case float64:
		val = float32(data)
	default:
		return &TypeAssertionError{data, "float32"}
	}
	buf.WriteFloat32LE(val)
	return nil
}
func (p *float32LEOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return buf.ReadFloat32LE(), nil
}

type float32BEOp struct{}

func (p *float32BEOp) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	*((*float32)(ptr)) = buf.ReadFloat32BE()
	return nil
}

func (p *float32BEOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	val := *((*float32)(ptr))
	buf.WriteFloat32BE(val)
	return nil
}
func (p *float32BEOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	var val float32
	switch data := data.(type) {
	case float32:
		val = float32(data)
	case float64:
		val = float32(data)
	default:
		return &TypeAssertionError{data, "float32"}
	}
	buf.WriteFloat32BE(val)
	return nil
}
func (p *float32BEOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return buf.ReadFloat32BE(), nil
}

type float64LEOp struct{}

func (p *float64LEOp) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	*((*float64)(ptr)) = buf.ReadFloat64LE()
	return nil
}

func (p *float64LEOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	val := *((*float64)(ptr))
	buf.WriteFloat64LE(val)
	return nil
}
func (p *float64LEOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	var val float64
	switch data := data.(type) {
	case float64:
		val = data
	default:
		return &TypeAssertionError{data, "float64"}
	}
	buf.WriteFloat64LE(val)
	return nil
}
func (p *float64LEOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return buf.ReadFloat64LE(), nil
}

type float64BEOp struct{}

func (p *float64BEOp) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	*((*float64)(ptr)) = buf.ReadFloat64BE()
	return nil
}

func (p *float64BEOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	val := *((*float64)(ptr))
	buf.WriteFloat64BE(val)
	return nil
}
func (p *float64BEOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	var val float64
	switch data := data.(type) {
	case float64:
		val = data
	default:
		return &TypeAssertionError{data, "float64"}
	}
	buf.WriteFloat64BE(val)
	return nil
}
func (p *float64BEOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return buf.ReadFloat64BE(), nil
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

func _encodeSliceTypeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
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

func _decodeSliceTypeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	var err error
	var m []interface{}
	var mPtr *[]interface{}

	vType := reflect2.TypeOf(v)
	if vType.Kind() == reflect.Ptr {
		mPtr = v.(*[]interface{})
		m = *mPtr
	} else {
		m = v.([]interface{})
		mPtr = &m
	}
	size, _ := buf.ReadVarUint()

	if cap(m) < int(size) {
		m = make([]interface{}, size)
	} else {
		m = m[:size]
	}
	childNode := opNode.children[0]
	for i := uint64(0); i < size; i++ {
		childProp := _createNewData(buf, childNode)
		m[i] = childProp
		m[i], err = childNode.handler.decodeDynamic(buf, childNode, m[i])
		if err != nil {
			return nil, err
		}
	}

	*mPtr = m

	return m, nil
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
