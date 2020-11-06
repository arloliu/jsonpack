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
schema definition from existing struct without write schema definition.

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
	sch, err := jsonPack.AddSchema("info", Info{}, jsonpack.BigEndian)
*/
func (p *JSONPack) AddSchema(schemaName string, v ...interface{}) (*Schema, error) {
	sch, err := p.schemaManager.add(schemaName, v...)
	if err != nil {
		return nil, &CompileError{schemaName, err}
	}
	return sch, nil
}

// GetSchema returns schema instance with schemaName, returns nil if schema doesn't exist.
func (p *JSONPack) GetSchema(schemaName string) *Schema {
	return p.schemaManager.get(schemaName)
}

// GetSchemaDef is a wrapper of Schema.GetSchemaDef, it gets a Schema instance by schemaName
// and call Schema instance's GetSchemaDef function.
func (p *JSONPack) GetSchemaDef(schemaName string) (*SchemaDef, error) {
	schema := p.schemaManager.get(schemaName)
	if schema == nil {
		return nil, &SchemaNonExistError{schemaName}
	}
	return schema.GetSchemaDef()
}

// GetSchemaDefText is a wrapper of Schema.GetSchemaDefText,
// it gets a Schema instance by schemaName
// and call Schema instance's GetSchemaDefText function,
// It returns error if schema doesn't exist.
func (p *JSONPack) GetSchemaDefText(schemaName string) ([]byte, error) {
	schema := p.schemaManager.get(schemaName)
	if schema == nil {
		return nil, &SchemaNonExistError{schemaName}
	}
	return schema.textData, nil
}

// Encode is a wrapper of Schema.Encode, it gets a Schema instance by schemaName
// and call Schema instance's Encode function.
func (p *JSONPack) Encode(schemaName string, v interface{}) ([]byte, error) {
	schema := p.schemaManager.get(schemaName)
	if schema == nil {
		return nil, &SchemaNonExistError{schemaName}
	}
	return schema.Encode(v)
}

// EncodeTo is a wrapper of Schema.EncodeTo, it gets a Schema instance by schemaName
// and call Schema instance's EncodeTo function.
func (p *JSONPack) EncodeTo(schemaName string, v interface{}, dataPtr *[]byte) error {
	schema := p.schemaManager.get(schemaName)
	if schema == nil {
		return &SchemaNonExistError{schemaName}
	}
	return schema.EncodeTo(v, dataPtr)
}

// Decode is a wrapper of Schema.Decode, it gets a Schema instance by schemaName
// and call Schema instance's Decode function.
func (p *JSONPack) Decode(schemaName string, data []byte, v interface{}) error {
	schema := p.schemaManager.get(schemaName)
	if schema == nil {
		return &SchemaNonExistError{schemaName}
	}
	return schema.decode(data, v, true)
}
