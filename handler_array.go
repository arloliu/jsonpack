package jsonpack

import (
	"reflect"
	"unsafe"

	ibuf "github.com/arloliu/jsonpack/buffer"
	"github.com/modern-go/reflect2"
)

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

func (p *arrayOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	return encodeSliceTypeDynamic(buf, opNode, data)
}

func (p *arrayOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	if reflect2.TypeOf(v).Kind() == reflect.Slice {
		return _sliceOp.decodeDynamic(buf, opNode, v)
	}
	return decodeArrayTypeDynamic(buf, opNode, v)
}
