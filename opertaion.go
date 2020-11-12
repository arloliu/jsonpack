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
