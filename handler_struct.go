package jsonpack

import (
	"unsafe"

	ibuf "github.com/arloliu/jsonpack/buffer"
	"github.com/modern-go/reflect2"
	"github.com/pkg/errors"
)

type structOp struct{}

func (p *structOp) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	var err error
	var fieldPtr unsafe.Pointer

	for _, op := range opNode.children {
		fieldPtr = op.field.UnsafeGet(ptr)

		if op.isPtrType {
			ptrType := (op.field.Type()).(*reflect2.UnsafePtrType)
			elemType := ptrType.Elem()

			// check if current pointer of field is null
			if derefPtr(fieldPtr) == nil {
				//pointer to null, need to allocate memory to hold the value
				newPtr := elemType.UnsafeNew()
				err = op.handler.decodeStruct(buf, op, newPtr)
				if err != nil {
					return err
				}
				// assign new allocated data back to current field pointer
				assignPtr(fieldPtr, newPtr)
			} else {
				//reuse existing instance
				err = op.handler.decodeStruct(buf, op, derefPtr(fieldPtr))
				if err != nil {
					return err
				}
			}
		} else {
			err = op.handler.decodeStruct(buf, op, fieldPtr)
			if err != nil {
				return err
			}
		}
	}

	return nil

}
func (p *structOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	var err error
	var fieldPtr unsafe.Pointer

	for _, op := range opNode.children {
		fieldPtr = op.field.UnsafeGet(ptr)
		// dereference pointer
		if op.isPtrType {
			// logger.Debugf("encodeStruct: op is pointer type, field name: %v", op.field.Name())
			fieldPtr = derefPtr(fieldPtr)
		}

		err = op.handler.encodeStruct(buf, op, fieldPtr)
		if err != nil {
			return err
		}
	}
	return nil
}
func (p *structOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	return errors.WithStack(&NotImplementedError{"structOp.encodeDynamic"})
}
func (p *structOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return nil, errors.WithStack(&NotImplementedError{"structOp.decodeDynamic"})
}
