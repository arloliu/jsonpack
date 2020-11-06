package jsonpack

// schemaManager structure
type schemaManager struct {
	schemas map[string]*Schema
}

// create schema manager instance
func newSchemaManager() *schemaManager {
	return &schemaManager{schemas: make(map[string]*Schema)}
}

// get schema instance
func (s *schemaManager) get(name string) *Schema {
	instance, ok := s.schemas[name]
	if !ok {
		return nil
	}
	return instance
}

// add a new schema or or replace existing one, compile it and store
func (s *schemaManager) add(name string, v ...interface{}) (*Schema, error) {
	schema := newSchema(s, name, v...)
	// add schema instance into map first, then build it
	_, ok := s.schemas[name]
	if ok {
		delete(s.schemas, name)
	}
	s.schemas[name] = schema

	err := schema.build()
	if err != nil {
		return nil, err
	}

	return schema, nil
}
