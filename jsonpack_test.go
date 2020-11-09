package jsonpack

import (
	"bytes"
	"os"
	"reflect"
	"testing"
)

type s1 struct {
	ID  string `json:"id"`
	Msg string `json:"msg"`
}

type s2 struct {
	ID  string `json:"id"`
	Msg string `json:"msg"`
}

func TestAddSchema(t *testing.T) {
	var jsonPacker *JSONPack = NewJSONPack()

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
	_, err := jsonPacker.AddSchema("info", schDef)
	if err != nil {
		t.Errorf("AddSchema info err: %v\n", err)
		os.Exit(1)
	}

	schDefMap := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{"type": "string"},
			"area": map[string]interface{}{"type": "uint32le"},
		},
		"order": []string{"name", "area"},
	}
	_, err = jsonPacker.AddSchema("info", schDefMap)
	if err != nil {
		t.Errorf("AddSchema info err: %v\n", err)
		os.Exit(1)
	}

	schDefSt := SchemaDef{
		Type: "object",
		Properties: map[string]*SchemaDef{
			"name": {Type: "string"},
			"area": {Type: "uint32le"},
		},
		Order: []string{"name", "area"},
	}
	_, err = jsonPacker.AddSchema("schdef_info", schDefSt)
	if err != nil {
		t.Errorf("AddSchema info err: %v\n", err)
		os.Exit(1)
	}
}

func TestGetAndRemoveSchema(t *testing.T) {
	var jsonPacker *JSONPack = NewJSONPack()
	var err error

	_, err = jsonPacker.AddSchema("s1", s1{})
	if err != nil {
		t.Errorf("AddSchema(s1) fail, err: %v", err)
	}

	s1Node := jsonPacker.GetSchema("s1")
	if s1Node == nil {
		t.Error("GetSchema(s1) fail")
	}

	def, _ := s1Node.GetSchemaDef()
	if _, ok := def.Properties["msg"]; !ok {
		t.Error("GetSchemaDef(s1) fail, should contains msg property")
	}
	if _, ok := def.Properties["id"]; !ok {
		t.Error("GetSchemaDef(s2) fail, should contains id property")
	}

	err = jsonPacker.RemoveSchema("s1")
	if err != nil {
		t.Error("RemoveSchema(s1) fail")
	}

	err = jsonPacker.RemoveSchema("nonexist")
	if err == nil {
		t.Error("RemoveSchema(nonexist) should return error")
	}

	jsonPacker.Reset()
}

func TestGetSchemaDefAndText(t *testing.T) {
	var jsonPacker *JSONPack = NewJSONPack()
	var err error

	s1Sch, err := jsonPacker.AddSchema("s1", s1{})
	if err != nil {
		t.Errorf("AddSchema(s1) fail, err: %v", err)
	}
	s1SchDefText := s1Sch.GetSchemaDefText()
	s1SchDef, err := s1Sch.GetSchemaDef()
	if err != nil {
		t.Errorf("GetSchemaDefText(s1) fail, err: %v", err)
	}

	mS1SchDefText, err := jsonPacker.GetSchemaDefText("s1")
	if err != nil || !bytes.Equal(mS1SchDefText, s1SchDefText) {
		t.Errorf("jsonpack.GetSchemaDefText(s1) fail, err: %v", err)
	}

	mS1SchDef, err := jsonPacker.GetSchemaDef("s1")
	if err != nil || !reflect.DeepEqual(mS1SchDef, s1SchDef) {
		t.Errorf("jsonpack.GetSchemaDef(s1) fail, err: %v", err)
	}

}

func TestGetAll(t *testing.T) {
	var jsonPacker *JSONPack = NewJSONPack()
	var err error

	s1Sch, err := jsonPacker.AddSchema("s1", s1{})
	if err != nil {
		t.Errorf("AddSchema(s1) fail, err: %v", err)
	}
	s1SchDefText := s1Sch.GetSchemaDefText()
	s1SchDef, err := s1Sch.GetSchemaDef()
	if err != nil {
		t.Errorf("GetSchemaDefText(s1) fail, err: %v", err)
	}

	s2Sch, err := jsonPacker.AddSchema("s2", s2{})
	if err != nil {
		t.Errorf("AddSchema(s2) fail, err: %v", err)
	}
	s2SchDefText := s2Sch.GetSchemaDefText()
	s2SchDef, err := s2Sch.GetSchemaDef()
	if err != nil {
		t.Errorf("GetSchemaDefText(s2) fail, err: %v", err)
	}

	schemas := jsonPacker.GetAllSchemas()

	if mS1Sch, ok := schemas["s1"]; !ok || mS1Sch != s1Sch {
		t.Error("GetAllSchemas test fail")
	}

	if mS2Sch, ok := schemas["s2"]; !ok || mS2Sch != s2Sch {
		t.Error("GetAllSchemas test fail")
	}

	schDefs := jsonPacker.GetAllSchemaDefs()
	if mS1SchDef, ok := schDefs["s1"]; !ok || !reflect.DeepEqual(mS1SchDef, s1SchDef) {
		t.Error("GetAllSchemaDefs test fail(s1)")
	}

	if mS2SchDef, ok := schDefs["s2"]; !ok || !reflect.DeepEqual(mS2SchDef, s2SchDef) {
		t.Error("GetAllSchemaDefs test fail(s2)")
	}

	schDefTexts := jsonPacker.GetAllSchemaDefTexts()
	if mS1SchDefText, ok := schDefTexts["s1"]; !ok || !bytes.Equal(mS1SchDefText, s1SchDefText) {
		t.Error("GetAllSchemaDefTexts test fail(s1)")
	}

	if mS2SchDefText, ok := schDefTexts["s2"]; !ok || !bytes.Equal(mS2SchDefText, s2SchDefText) {
		t.Error("GetAllSchemaDefTexts test fail(s2)")
	}

}
