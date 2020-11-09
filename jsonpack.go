package jsonpack

// JSONPack JSON packer structure
type JSONPack struct {
	schemaManager *schemaManager
}

// NewJSONPack Create instance of JSONPacker
func NewJSONPack() *JSONPack {
	instance := JSONPack{}
	instance.schemaManager = newSchemaManager()
	return &instance
}

/*
AddSchema will compile schema definition and stores compiled result in internal schema manager.

It's a variadic function which accept three types of input parameters in the following.

AddSchema(schemaName string, v interface{})

The v is schema definition which want to compile.
The value of v can be a JSON format of text data with []byte or string type, a map presents JSON
format of schema definition or a SchemaDef struct presents schema definition.

Example of add new schema from JSON text string:
	schDef := `
	{
		"type": "object",
		"properties": {
			"name": {"type": "string"},
			"area": {"type": "uint32le"}
		},
		"order": ["name", "area"]
	}
	`

	jsonPack := jsonpack.NewJSONPack()
	sch, err := jsonPack.AddSchema("info", schDef)

Example of adding new schema from map of schema definition:

	schDef := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{"type": "string"},
			"area": map[string]interface{}{"type": "uint32le"},
		},
		"order": []string{"name", "area"},
	}

	jsonPack := jsonpack.NewJSONPack()
	sch, err := jsonPack.AddSchema("info", schDef)

Example of adding new schema from SchemaDef struct:

	schDef := jsonpack.SchemaDef{
		Type: "object",
		Properties: map[string]*jsonpack.SchemaDef{
			"name": {Type: "string"},
			"area": {Type: "uint32le"},
		},
		Order: []string{"name", "area"},
	}
	jsonPack := jsonpack.NewJSONPack()
	sch, err := jsonPack.AddSchema("info", schDef)


AddSchema(schemaName string, v interface{}, byteOrder jsonpack.ByteOrder)

For fast prototyping and fast prototyping, AddSchema method supports generate
schema definition from existing struct without write schema definition by hand.

In this scenario, the value of v is the struct which want to generate,
and byteOrder parameter indicates the byte order, can be either jsonpack.LittleEndian
or jsonpack.BigEndian, it defaults to little-endian if byteOrder not specified.

Example of adding new schema and build schema definition from struct:

	type Info struct {
		Name string `json:"name"`
		Area uint32 `json:"area"`
		// omit this field
		ExcludeField string `-`
	}

	jsonPack := jsonpack.NewJSONPack()
	sch, err := jsonPack.AddSchema("Info", Info{}, jsonpack.BigEndian)
*/
func (p *JSONPack) AddSchema(schemaName string, v ...interface{}) (*Schema, error) {
	sch, err := p.schemaManager.add(schemaName, v...)
	if err != nil {
		return nil, &CompileError{schemaName, err}
	}
	return sch, nil
}

// EncodeTo is a wrapper of Schema.EncodeTo,
// it returns *SchemaNonExistError error if schema doesn't exist.
func (p *JSONPack) EncodeTo(schemaName string, v interface{}, dataPtr *[]byte) error {
	schema := p.schemaManager.get(schemaName)
	if schema == nil {
		return &SchemaNonExistError{schemaName}
	}
	return schema.EncodeTo(v, dataPtr)
}

// Decode is a wrapper of Schema.Decode,
// it returns *SchemaNonExistError error if schema doesn't exist.
func (p *JSONPack) Decode(schemaName string, data []byte, v interface{}) error {
	schema := p.schemaManager.get(schemaName)
	if schema == nil {
		return &SchemaNonExistError{schemaName}
	}
	return schema.decode(data, v, true)
}

// GetSchema returns schema instance by schemaName, returns nil if schema doesn't exist.
func (p *JSONPack) GetSchema(schemaName string) *Schema {
	return p.schemaManager.get(schemaName)
}

// GetSchemaDef is a wrapper of Schema.GetSchemaDef, it gets a Schema instance by schemaName,
// it returns *SchemaNonExistError error if schema doesn't exist.
func (p *JSONPack) GetSchemaDef(schemaName string) (*SchemaDef, error) {
	schema := p.schemaManager.get(schemaName)
	if schema == nil {
		return nil, &SchemaNonExistError{schemaName}
	}
	return schema.GetSchemaDef()
}

// GetSchemaDefText is a wrapper of Schema.GetSchemaDefText,
// it returns *SchemaNonExistError error if schema doesn't exist.
func (p *JSONPack) GetSchemaDefText(schemaName string) ([]byte, error) {
	schema := p.schemaManager.get(schemaName)
	if schema == nil {
		return nil, &SchemaNonExistError{schemaName}
	}
	return schema.textData, nil
}

// GetAllSchemas returns a map which contains all existed schema instances,
// key of map it schema name, and value of map is *Schema.
func (p *JSONPack) GetAllSchemas() map[string]*Schema {
	return p.schemaManager.getAllSchemas()
}

// GetAllSchemaDefs returns a map which contains all existed schema definitions,
// key of map it schema name, and value of map is *SchemaDef.
func (p *JSONPack) GetAllSchemaDefs() map[string]*SchemaDef {
	return p.schemaManager.getAllSchemaDefs()
}

// GetAllSchemaDefTexts returns a map which contains all existed schema text definitions,
// key of map it schema name, and value of map is text format of schema defintion which
// presented as []byte.
func (p *JSONPack) GetAllSchemaDefTexts() map[string][]byte {
	return p.schemaManager.getAllSchemaDefTexts()
}

// Encode is a wrapper of Schema.Encode,
// it returns *SchemaNonExistError error if schema doesn't exist.
func (p *JSONPack) Encode(schemaName string, v interface{}) ([]byte, error) {
	schema := p.schemaManager.get(schemaName)
	if schema == nil {
		return nil, &SchemaNonExistError{schemaName}
	}
	return schema.Encode(v)
}

// RemoveSchema removes schema by schemaName, it returns *SchemaNonExistError error
// if schema doesn't exist.
func (p *JSONPack) RemoveSchema(schemaName string) error {
	return p.schemaManager.remove(schemaName)
}

// Reset removes all schema instances
func (p *JSONPack) Reset() {
	p.schemaManager.reset()
}
