package jsonpack

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"sync"
)

type ByteOrder int

const (
	LittleEndian ByteOrder = 0
	BigEndian    ByteOrder = 1
)

// Schema of jsonpack, each Schema instance presents a schema added by AddSchema function
type Schema struct {
	// schema name
	Name string
	// rawData stores schema definition from user, can be map, struct, string or slice of byte
	rawData interface{}
	// textData stores text json format of schema defintion
	textData      []byte
	manager       *schemaManager
	rootOp        *operation
	structOpCache *sync.Map
	encodeBufSize int64
	decodeBufSize int64
	byteOrder     ByteOrder
}

// The schema defintion presents the definition that define the structure of JSON data.
type SchemaDef struct {
	Type       string                `json:"type"`
	Properties map[string]*SchemaDef `json:"properties,omitempty"`
	Items      *SchemaDef            `json:"items,omitempty"`
	Order      []string              `json:"order,omitempty"`
}

// GetSchemaDef returns schema definition instance,
// returns nil and error if error occurs.
func (s *Schema) GetSchemaDef() (*SchemaDef, error) {
	schDef := SchemaDef{}
	err := json.Unmarshal(s.textData, &schDef)
	if err != nil {
		return nil, err
	}
	return &schDef, nil
}

// GetSchemaDefText returns JSON format document of schema definition.
func (s *Schema) GetSchemaDefText(schemaName string) []byte {
	return s.textData
}

// create a schema instance by raw schema definition
func newSchema(manager *schemaManager, name string, v ...interface{}) *Schema {
	var rawData interface{}
	var byteOrder ByteOrder = LittleEndian
	if len(v) == 1 {
		rawData = v[0]
	} else if len(v) >= 2 {
		rawData = v[0]
		byteOrder = v[1].(ByteOrder)
	}

	instance := Schema{
		Name:          name,
		rawData:       rawData,
		textData:      nil,
		manager:       manager,
		rootOp:        newOperation("", &nullOp{}, nullOpType),
		structOpCache: &sync.Map{},
		encodeBufSize: 512,
		decodeBufSize: 512,
		byteOrder:     byteOrder,
	}
	return &instance
}

func (s *Schema) build() error {
	var ok bool
	rawDataType := reflect.TypeOf(s.rawData)

	// check if raw data is SchemaDef or pointer of SchemaDef
	switch schDef := s.rawData.(type) {
	case SchemaDef:
		textData, err := json.Marshal(schDef)
		if err != nil {
			return err
		}
		return s.buildFromTextDef(textData)

	case *SchemaDef:
		textData, err := json.Marshal(*schDef)
		if err != nil {
			return err
		}
		return s.buildFromTextDef(textData)
	}

	switch rawDataType.Kind() {
	case reflect.Map:
		var schemaMap map[string]interface{}
		schemaMap, ok = s.rawData.(map[string]interface{})
		if !ok {
			return errors.New("the map type of schema definition needs to be map[string]interface{}")
		}
		return s.buildFromMap(schemaMap)

	case reflect.Struct:
		return s.buildFromStruct(rawDataType)

	case reflect.String:
		textDef := []byte(s.rawData.(string))
		return s.buildFromTextDef(textDef)

	case reflect.Slice, reflect.Array:
		elemType := rawDataType.Elem()
		switch elemType.Kind() {
		case reflect.Uint8: // []byte
			return s.buildFromTextDef(s.rawData.([]byte))

		default:
			return errors.New("the slice type of schema definition needs to be []byte")
		}

	default:
		return errors.New("schema definition invalid, need to be a struct, map, []bytes or string type")
	}
}

func (s *Schema) buildFromStruct(sType reflect.Type) error {
	var err error
	st := SchemaDef{}
	err = s._buildFromStruct(&st, sType)
	if err != nil {
		return err
	}

	var textData []byte
	textData, err = json.Marshal(st)
	if err != nil {
		return err
	}

	return s.buildFromTextDef(textData)
}

func (s *Schema) _buildFromStruct(st *SchemaDef, sType reflect.Type) error {
	var err error
	sKind := sType.Kind()
	switch sKind {
	case reflect.Slice, reflect.Array:
		st.Type = "array"
		st.Items = &SchemaDef{}
		err = s._buildFromStruct(st.Items, sType.Elem())
		if err != nil {
			return err
		}

	case reflect.Struct:
		st.Type = "object"
		st.Properties = make(map[string]*SchemaDef)
		st.Order = make([]string, sType.NumField())
		for i := 0; i < sType.NumField(); i++ {
			field := sType.Field(i)

			var fieldName string

			// lookup field name
			jsonTag, ok := field.Tag.Lookup("json")
			if !ok {
				fieldName = field.Name
			} else {
				// this field is omitted, skip
				if jsonTag == "-" {
					continue
				}
				tagParts := strings.Split(jsonTag, ",")
				if tagParts[0] != "" {
					fieldName = tagParts[0]
				} else {
					fieldName = field.Name
				}
			}

			st.Properties[fieldName] = &SchemaDef{}
			err = s._buildFromStruct(st.Properties[fieldName], field.Type)
			if err != nil {
				return err
			}

			st.Order[i] = fieldName
		}

	case reflect.Ptr:
		return s._buildFromStruct(st, sType.Elem())

	default:
		st.Type = s.getTypeFromKind(sKind)
		if st.Type == "" {
			return &UnknownTypeError{sType.String()}
		}
	}

	return nil
}

func (s *Schema) buildFromTextDef(textDef []byte) error {
	var err error

	schema := make(map[string]interface{})
	err = json.Unmarshal(textDef, &schema)
	if err != nil {
		return err
	}

	return s.buildFromMap(schema)
}

func (s *Schema) buildFromMap(schema map[string]interface{}) error {
	var err error

	schType, ok := schema["type"].(string)
	if !ok {
		return errors.New("Need 'type' property in top-level schema definition")
	}

	schType = strings.ToLower(schType)
	switch schType {
	case "object":
		err = s.compileSchemaObject(schema, s.rootOp)
		if err != nil {
			return err
		}
		s.rootOp.handler = &objectOp{}
		s.rootOp.handlerType = objectOpType
	case "array":
		err = s.compileSchemaArray(schema, s.rootOp)
		if err != nil {
			return err
		}
		s.rootOp.handler = _arrayOp
		s.rootOp.handlerType = arrayOpType
	default:
		return errors.New("type property needs to be 'object' or 'array' in top-level schema definition")
	}
	// parse schema map object into JSON encoded text data
	s.textData, err = json.Marshal(schema)
	if err != nil {
		return err
	}

	return nil
}

func _getProperties(prop interface{}) map[string]interface{} {
	data := make(map[string]interface{})
	switch prop := prop.(type) {
	case map[string]interface{}:
		return prop
	case map[string]string:
		for k, v := range prop {
			data[k] = v
		}
		return data
	}
	return nil
}

func _getOrder(order interface{}) []string {
	switch order := order.(type) {
	case []interface{}:
		data := make([]string, len(order))
		for i, v := range order {
			data[i] = v.(string)
		}
		return data

	case []string:
		return order
	}
	return nil
}

func (s *Schema) compileSchemaObject(schema map[string]interface{}, curOp *operation) error {
	var err error
	err = checkObjectProperties(schema)
	if err != nil {
		return err
	}

	var newOp *operation
	properties := _getProperties(schema["properties"])
	if properties == nil {
		return errPropertiesProp
	}

	order := _getOrder(schema["order"])
	if order == nil {
		return errOrderProp
	}

	for _, fieldName := range order {
		prop := properties[fieldName].(map[string]interface{})
		propType, ok := prop["type"].(string)
		if !ok {
			return errors.New("Object type of schema definition requires valid 'type' field")
		}
		propType = strings.ToLower(propType)

		if propType == "object" {
			// create and append object operation to child operation list
			newOp = newOperation(fieldName, &objectOp{}, objectOpType)
			curOp.children = append(curOp.children, newOp)

			err = s.compileSchemaObject(prop, newOp)

		} else if propType == "array" {
			// create and append object operation to child operation list
			newOp = newOperation(fieldName, _arrayOp, arrayOpType)
			curOp.children = append(curOp.children, newOp)
			err = s.compileSchemaArray(prop, newOp)
		} else if isBuiltinType(&propType) {
			handler := builtinTypes[propType]
			newOp = newOperation(fieldName, handler, builtinOpHandlerTypes[propType])
			curOp.children = append(curOp.children, newOp)
		} else {
			return &UnknownTypeError{propType}
		}
		if err != nil {
			return err
		}
	}
	return err
}

func (s *Schema) compileSchemaArray(schema map[string]interface{}, curOp *operation) error {
	var err error
	err = checkArrayProperties(schema)
	if err != nil {
		return err
	}

	items := schema["items"].(map[string]interface{})
	itemType := items["type"].(string)
	itemType = strings.ToLower(itemType)

	var newOp *operation
	if itemType == "object" {
		// create and append object operation to child operation list
		newOp = newOperation("", &objectOp{}, objectOpType)
		curOp.children = append(curOp.children, newOp)
		err = s.compileSchemaObject(items, newOp)
	} else if itemType == "array" {
		newOp = newOperation("", _arrayOp, arrayOpType)
		curOp.children = append(curOp.children, newOp)
		err = s.compileSchemaArray(items, newOp)
	} else if isBuiltinType(&itemType) {
		handler := builtinTypes[itemType]
		newOp = newOperation("", handler, builtinOpHandlerTypes[itemType])
		curOp.children = append(curOp.children, newOp)
	} else {
		return &UnknownTypeError{itemType}
	}
	return err
}

const uintSize = 32 << (^uint(0) >> 32 & 1)

func (s *Schema) getTypeFromKind(kind reflect.Kind) string {
	switch kind {
	case reflect.Bool:
		return "boolean"
	case reflect.String:
		return "string"
	case reflect.Int:
		if uintSize == 32 {
			return s.getTypeEndian("int32")
		} else {
			return s.getTypeEndian("int64")
		}
	case reflect.Int8:
		return "uint8"
	case reflect.Int16:
		return s.getTypeEndian("int16")
	case reflect.Int32:
		return s.getTypeEndian("int32")
	case reflect.Int64:
		return s.getTypeEndian("int64")
	case reflect.Uint:
		if uintSize == 32 {
			return s.getTypeEndian("uint32")
		} else {
			return s.getTypeEndian("uint64")
		}
	case reflect.Uint8:
		return "uint8"
	case reflect.Uint16:
		return s.getTypeEndian("uint16")
	case reflect.Uint32:
		return s.getTypeEndian("uint32")
	case reflect.Uint64:
		return s.getTypeEndian("uint64")
	case reflect.Float32:
		return s.getTypeEndian("float32")
	case reflect.Float64:
		return s.getTypeEndian("float64")
	default:
		return ""
	}
}

func (s *Schema) getTypeEndian(typ string) string {
	if s.byteOrder == BigEndian {
		return typ + "BE"
	}
	return typ + "LE"
}

// SetEncodeBufSize set default encode buffer allocation bytes
func (s *Schema) SetEncodeBufSize(size int64) {
	s.encodeBufSize = size
}

// SetDecodeBufSize set default decode buffer allocation bytes
func (s *Schema) SetDecodeBufSize(size int64) {
	s.decodeBufSize = size
}

func maxInt64(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

func checkObjectProperties(schema map[string]interface{}) error {
	var okIface, okStr bool

	_, okIface = schema["properties"].(map[string]interface{})
	_, okStr = schema["properties"].(map[string]string)
	if !okIface && !okStr {
		return errPropertiesProp
	}

	_, okIface = schema["order"].([]interface{})
	_, okStr = schema["order"].([]string)
	if !okIface && !okStr {
		return errOrderProp
	}
	return nil
}

func checkArrayProperties(schema map[string]interface{}) error {
	items, ok := schema["items"].(map[string]interface{})
	if !ok {
		return errItemsProp
	}

	_, ok = items["type"].(string)
	if !ok {
		return errTypeProp
	}
	return nil
}
