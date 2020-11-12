package jsonpack

import (
	"unsafe"

	ibuf "github.com/arloliu/jsonpack/buffer"
	"github.com/modern-go/reflect2"
)

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
