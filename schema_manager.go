package jsonpack

import (
	"sync"
)

// schemaManager manages schema instances
type schemaManager struct {
	schemas sync.Map // provides thread safety map
}

// newSchemaManager returns a new schema manager instance
func newSchemaManager() *schemaManager {
	return &schemaManager{}
}

// getAll returns a cloned map of schema instances
func (s *schemaManager) getAllSchemas() map[string]*Schema {
	schemas := make(map[string]*Schema)
	s.schemas.Range(func(key, value interface{}) bool {
		schemas[key.(string)] = value.(*Schema)
		return true
	})
	return schemas
}

// getAllSchemaDefs returns a map which contains all existed schema definitions,
// key of map it schema name, and value of map is *SchemaDef.
func (s *schemaManager) getAllSchemaDefs() map[string]*SchemaDef {
	schDefs := make(map[string]*SchemaDef)
	s.schemas.Range(func(key, value interface{}) bool {
		schema := value.(*Schema)
		schDef, err := schema.GetSchemaDef()
		if err == nil {
			schDefs[key.(string)] = schDef
		}
		return true
	})
	return schDefs
}

// getAllSchemaDefTexts returns a map which contains all existed schema text definitions,
// key of map it schema name, and value of map is text format of schema defintion, which
// presented as []byte.
func (s *schemaManager) getAllSchemaDefTexts() map[string][]byte {
	schDefTexts := make(map[string][]byte)
	s.schemas.Range(func(key, value interface{}) bool {
		schema := value.(*Schema)
		schDefTexts[key.(string)] = schema.GetSchemaDefText()
		return true
	})
	return schDefTexts
}

// get returns schema instance in schema manager or returns nil if schema not found.
func (s *schemaManager) get(name string) *Schema {
	instance, ok := s.schemas.Load(name)
	if !ok {
		return nil
	}
	return instance.(*Schema)
}

// add a new schema or or replace existing one, compile and store it in schema manager.
// It returns error if can't not build new schema and it will not add to schemna manager.
func (s *schemaManager) add(name string, v ...interface{}) (*Schema, error) {
	schema := newSchema(name, v...)
	err := schema.build()
	if err != nil {
		return nil, err
	}

	s.schemas.Delete(name)
	s.schemas.Store(name, schema)

	return schema, nil
}

// remove schema by name in schema manager, it returns *SchemaNonExistError error
// if schema doesn't exist.
func (s *schemaManager) remove(name string) error {
	_, ok := s.schemas.LoadAndDelete(name)
	if !ok {
		return &SchemaNonExistError{name}
	}
	return nil
}

// reset removes all schema instance in schema manager.
func (s *schemaManager) reset() {
	s.schemas.Range(func(key, value interface{}) bool {
		s.schemas.Delete(key)
		return true
	})
}
