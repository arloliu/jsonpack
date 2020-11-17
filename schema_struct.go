package jsonpack

import (
	"reflect"
	"strings"

	"github.com/modern-go/reflect2"
)

func (s *Schema) getStructOperation(typ reflect2.Type, st interface{}) (*structOperation, error) {
	cacheKey := typ.RType()
	sop, ok := s.structOpCache.Load(cacheKey)
	// sop, ok := s.structOpCache[cacheKey]
	if !ok {
		newOp, err := s.buildStructOperation(s.rootOp, st)
		if err != nil {
			return nil, err
		}
		s.structOpCache.Store(cacheKey, newOp)
		return newOp, nil
	}
	return sop.(*structOperation), nil
}

func (s *Schema) buildStructOperation(op *operation, st interface{}) (*structOperation, error) {
	sop := newStructOperation(_nullOp, nullOpType)
	err := s._buildStructOperation(sop, op, reflect2.TypeOf(st))
	if err != nil {
		return nil, err
	}
	return sop, nil
}

func (s *Schema) _buildStructOperation(sop *structOperation, op *operation, typ reflect2.Type) error {
	var err error
	sop.handler = op.handler
	sop.handlerType = op.handlerType
	sop.opType = typ
	sop.opInstance = op
	sop.isPtrType = (typ.Kind() == reflect.Ptr)

	switch op.handlerType {
	case objectOpType:
		// skip dynamic operation
		if typ.Kind() == reflect.Map {
			sop.opInstance = cloneAnonymousObjectOp(op)
			return nil
		}
		var structType *reflect2.UnsafeStructType
		structType, err = toStructType(typ)
		if err != nil {
			return err
		}

		fieldMap := parseStructFieldName(structType)

		// overwrite object operation to struct operation
		sop.handler = _structOp
		sop.handlerType = structOpType

		for _, opNode := range op.children {
			childOp := newStructOperation(opNode.handler, opNode.handlerType)
			childOp.field = structType.FieldByName(fieldMap[opNode.propName])
			if childOp.field == nil {
				return &StructFieldNonExistError{typ.String(), opNode.propName}
			}
			// childOp.fieldType = childOp.field.Type()
			err = s._buildStructOperation(childOp, opNode, childOp.field.Type())
			if err != nil {
				return err
			}
			sop.appendChild(childOp)
		}

	case sliceOpType, arrayOpType:
		switch typ := typ.(type) {
		case *reflect2.UnsafeArrayType:
			sop.handler = _arrayOp
			sop.handlerType = arrayOpType
			itemOp := op.children[0]
			childOp := newStructOperation(itemOp.handler, itemOp.handlerType)
			err = s._buildStructOperation(childOp, itemOp, typ.Elem())
			if err != nil {
				return err
			}
			sop.appendChild(childOp)

		case *reflect2.UnsafeSliceType:
			sop.handler = _sliceOp
			sop.handlerType = sliceOpType
			itemOp := op.children[0]
			childOp := newStructOperation(itemOp.handler, itemOp.handlerType)
			err = s._buildStructOperation(childOp, itemOp, typ.Elem())
			if err != nil {
				return err
			}
			sop.appendChild(childOp)
		default:
			return &UnknownTypeError{typ.String()}
		}
	default:
		// do nothing with builtin operations
	}
	return nil
}

type jsonTag struct {
	name       string
	omit       bool
	omitEmpty  bool
	stringMode bool
}

func parseJsonTag(field *reflect2.StructField) (*jsonTag, bool) {
	tag, ok := (*field).Tag().Lookup("json")
	if !ok {
		return nil, false
	}
	data := &jsonTag{}
	//  if the field tag is "-", the field is always omitted.
	if tag == "-" {
		data.omit = true
		return data, true
	}

	// data.private = unicode.IsLower(rune((*field).Name()[0]))

	tagParts := strings.Split(tag, ",")
	if tagParts[0] != "" {
		data.name = tagParts[0]
	}
	if len(tagParts) > 1 {
		for _, tagPart := range tagParts[1:] {
			switch tagPart {
			case "omitempty":
				data.omitEmpty = true
			case "string":
				data.stringMode = true
			}
		}
	}

	return data, true
}
func cloneAnonymousObjectOp(op *operation) *operation {
	newOp := operation{
		propName: "",
		handler:  op.handler,
		children: op.children,
	}
	return &newOp
}

func parseStructFieldName(structType *reflect2.UnsafeStructType) map[string]string {
	numField := structType.NumField()
	fieldMap := make(map[string]string, numField)

	for i := 0; i < numField; i++ {
		field := structType.Field(i)
		tag, ok := parseJsonTag(&field)
		var mapKey string
		if !ok || tag.omit || tag.name == "" {
			mapKey = field.Name()
		} else {
			mapKey = tag.name
		}
		fieldMap[mapKey] = field.Name()
	}
	return fieldMap
}

func toPtrElemType(typ reflect2.Type) reflect2.Type {
	ptrType := typ.(*reflect2.UnsafePtrType)
	return ptrType.Elem()
}

func toStructType(dType reflect2.Type) (*reflect2.UnsafeStructType, error) {
	dKind := dType.Kind()
	if dKind == reflect.Ptr {
		dType = toPtrElemType(dType)
		dKind = dType.Kind()
	}
	if dKind != reflect.Struct {
		return nil, &WrongTypeError{dType.String()}
	}

	return dType.(*reflect2.UnsafeStructType), nil
}
