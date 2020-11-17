package jsonpack

import (
	"reflect"
	"unsafe"

	ibuf "github.com/arloliu/jsonpack/buffer"
	"github.com/modern-go/reflect2"
)

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
			ptrType := (itemOp.opType).(*reflect2.UnsafePtrType)
			// always allocate new struct instance
			elemType := ptrType.Elem()
			newPtr = elemType.UnsafeNew()
		}
		err = itemOp.handler.decodeStruct(buf, itemOp, newPtr)
		if err != nil {
			return err
		}

		if itemPtr != newPtr {
			assignPtr(itemPtr, newPtr)
		}
	}
	return nil
}
func (p *sliceOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	var err error

	sliceType := opNode.opType.(*reflect2.UnsafeSliceType)
	length := sliceType.UnsafeLengthOf(ptr)
	itemOp := opNode.children[0]

	if sliceType.UnsafeIsNil(ptr) {
		// allocate new slice
		sliceType.UnsafeSet(ptr, sliceType.UnsafeMakeSlice(length, length))
	} else if sliceType.UnsafeLengthOf(ptr) < length {
		// grow slice
		sliceType.UnsafeGrow(ptr, length)
	}

	buf.WriteVarUint(uint64(length))
	for i := 0; i < length; i++ {
		itemPtr := sliceType.UnsafeGetIndex(ptr, i)
		// dereference pointer
		if itemOp.isPtrType {
			itemPtr = derefPtr(itemPtr)
		}

		err = itemOp.handler.encodeStruct(buf, itemOp, itemPtr)
		if err != nil {
			return err
		}
	}
	return nil
}
func (p *sliceOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	return encodeSliceTypeDynamic(buf, opNode, data)
}

func (p *sliceOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	if reflect2.TypeOf(v).Kind() == reflect.Array {
		return _arrayOp.decodeDynamic(buf, opNode, v)
	}
	return decodeSliceTypeDynamic(buf, opNode, v)
}
