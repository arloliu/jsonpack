package jsonpack

import (
	"github.com/modern-go/reflect2"
)

type operation struct {
	propName    string
	handler     opHandler
	handlerType opHandlerType
	children    []*operation
}

type structOperation struct {
	handler     opHandler
	handlerType opHandlerType
	opType      reflect2.Type
	opInstance  *operation
	field       reflect2.StructField
	isPtrType   bool
	children    []*structOperation
}

func newOperation(field string, handler opHandler, handlerType opHandlerType) *operation {
	return &operation{field, handler, handlerType, make([]*operation, 0)}
}

func newStructOperation(handler opHandler, handlerType opHandlerType) *structOperation {
	op := &structOperation{}
	op.handler = handler
	op.handlerType = handlerType
	op.children = make([]*structOperation, 0)
	return op
}

func (s *structOperation) appendChild(child *structOperation) {
	s.children = append(s.children, child)
}

//nolint:deadcode,unused
func handlerName(handlerType opHandlerType) string {
	switch handlerType {
	case nullOpType:
		return "nullOp"
	case objectOpType:
		return "objectOp"
	case structOpType:
		return "structOp"
	case arrayOpType:
		return "arrayOp"
	case sliceOpType:
		return "sliceOp"
	default:
		for k, v := range builtinOpHandlerTypes {
			if handlerType == v {
				return k + "Op"
			}
		}
	}
	return "unknownOp"
}

// //nolint:deadcode,unused
// func (s *structOperation) dump() {
// 	level := 0
// 	_dump(s, level)
// }

// //nolint:deadcode,unused
// func _dump(op *structOperation, level int) {
// 	var space string
// 	for i := 0; i < level; i++ {
// 		space = space + "    "
// 	}

// 	name := handlerName(op.handlerType)
// 	var fieldName string = "nil"
// 	if op.field != nil {
// 		fieldName = op.field.Name()
// 	}

// 	logger.Debugf("%sHandler: %s Type: %v Field: %v propName: %v", space, name, op.opType.String(), fieldName, op.opInstance.propName)

// 	for _, childOp := range op.children {
// 		_dump(childOp, level+1)
// 	}
// }
