package jsonpack

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/arloliu/jsonpack/testdata"
)

var jsonPack *JSONPack

func init() {
	var err error
	jsonPack = NewJSONPack()

	_, err = jsonPack.AddSchema("complex", testdata.ComplexSchDef)
	if err != nil {
		fmt.Printf("AddSchema complex err: %v\n", err)
		os.Exit(1)
	}

	_, err = jsonPack.AddSchema("arrayObject", testdata.ArraySchDef)
	if err != nil {
		fmt.Printf("AddSchema arrayObject err: %v\n", err)
		os.Exit(1)
	}
	_, err = jsonPack.AddSchema("testStruct", testdata.StructSchDef)
	if err != nil {
		fmt.Printf("AddSchema testStruct err: %v\n", err)
		os.Exit(1)
	}
}

func TestComplexAddSchemaFromStruct(t *testing.T) {
	_, err := jsonPack.AddSchema("complex_from_struct", testdata.Complex{})
	if err != nil {
		t.Errorf("AddSchema fail, error: %v", err)
	}
}

func TestAddSchema(t *testing.T) {
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
	_, err := jsonPack.AddSchema("info", schDef)
	if err != nil {
		fmt.Printf("AddSchema info err: %v\n", err)
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
	_, err = jsonPack.AddSchema("info", schDefMap)
	if err != nil {
		fmt.Printf("AddSchema info err: %v\n", err)
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
	_, err = jsonPack.AddSchema("schdef_info", schDefSt)
	if err != nil {
		fmt.Printf("AddSchema info err: %v\n", err)
		os.Exit(1)
	}
}

func TestStructEncode(t *testing.T) {
	encData, err := jsonPack.Encode("testStruct", &testdata.StructData)
	if err != nil {
		t.Errorf("Encode fail, error: %v", err)
	}

	if !compareBytes(t, encData, testdata.StructExpData) {
		t.Errorf("Encoded data mismatch")
	}

}

func TestComplexEncode(t *testing.T) {
	var err error
	var encData []byte
	encData, err = jsonPack.Encode("complex", testdata.ComplexData)
	if err != nil {
		t.Errorf("Encode fail, error: %v", err)
	}

	if !compareBytes(t, encData, testdata.ComplexExpData) {
		t.Errorf("Encoded data mismatch")
	}

	var encData2 []byte
	encData2, err = jsonPack.Encode("complex", &testdata.ComplexStructData)
	if err != nil {
		t.Errorf("Encode struct fail, error: %v", err)
	}
	if !compareBytes(t, encData2, testdata.ComplexExpData) {
		t.Errorf("Encoded complex data mismatch")
	}

	_, err = jsonPack.AddSchema("complex_from_struct", testdata.Complex{})
	if err != nil {
		t.Errorf("AddSchema fail, error: %v", err)
	}

	var encData3 []byte
	encData3, err = jsonPack.Encode("complex_from_struct", &testdata.ComplexStructData)
	if err != nil {
		t.Errorf("Encode struct fail, error: %v", err)
	}
	if !compareBytes(t, encData3, testdata.ComplexExpData) {
		t.Errorf("Encoded complex data mismatch")
	}

	var encData4 []byte = make([]byte, 128)
	err = jsonPack.EncodeTo("complex", &testdata.ComplexStructData, &encData4)
	if err != nil {
		t.Errorf("Encode struct fail, error: %v", err)
	}
	if !compareBytes(t, encData4, testdata.ComplexExpData) {
		t.Errorf("Encoded complex data mismatch")
	}

}

func TestArrayEncode(t *testing.T) {
	var err error

	// test encode slice of maps
	encData, err := jsonPack.Encode("arrayObject", testdata.ArrayData)
	if err != nil {
		t.Errorf("Encode fail, error: %v", err)
	}

	if !compareBytes(t, encData, testdata.ArrayExpData) {
		t.Errorf("Encoded data mismatch")
	}

	// test encode slice of structs
	if err != nil {
		t.Errorf("Unmarshal array struct fail, error: %v", err)
	}

	encData2, err := jsonPack.Encode("arrayObject", &testdata.ArrayStructData)

	if err != nil {
		t.Errorf("Encode array struct fail, error: %v", err)
	}

	if !compareBytes(t, encData2, testdata.ArrayExpData) {
		t.Errorf("Encoded data mismatch")
	}

}

func TestArrayDecode(t *testing.T) {
	var err error

	// decodeData := make([]interface{}, 10)
	var decodeData []interface{}
	err = jsonPack.Decode("arrayObject", testdata.ArrayExpData, &decodeData)
	if err != nil {
		t.Errorf("Decode fail, err: %v", err)
	}
	if !compareMap(decodeData, testdata.ArrayData) {
		t.Errorf("Decode fail, compareMap fail")

	}
	m1, _ := json.Marshal(decodeData)
	m2, _ := json.Marshal(testdata.ArrayData)
	if !reflect.DeepEqual(m1, m2) {
		t.Errorf("Decode fail, marshaled data mismatch")
	}

	// decodeStructData := make([]testdata.TestArrayStruct, 0)
	var decodeStructData []testdata.TestArrayStruct
	err = jsonPack.Decode("arrayObject", testdata.ArrayExpData, &decodeStructData)
	if err != nil {
		t.Errorf("Decode array struct fail, err: %v", err)
	}

	if !reflect.DeepEqual(&testdata.ArrayStructData, &decodeStructData) {
		t.Errorf("Decode array struct data mismatch")
	}

}

func TestStructDecode(t *testing.T) {
	var err error
	decodeData := testdata.TestStruct{}

	err = jsonPack.Decode("testStruct", testdata.StructExpData, &decodeData)
	if err != nil {
		t.Errorf("Decode struct fail, err: %v", err)
	}

	if !reflect.DeepEqual(&testdata.StructData, &decodeData) {
		t.Errorf("Decode struct data mismatch")
	}

}

func toFloat64(a reflect.Value) float64 {
	switch a.Type().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(a.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(a.Uint())
	case reflect.Float32, reflect.Float64:
		return float64(a.Float())
	default:
		return 0
	}
}

func compareMap(a interface{}, b interface{}) bool {
	aVal := reflect.ValueOf(a)
	bVal := reflect.ValueOf(b)
	aType := aVal.Type()
	bType := bVal.Type()
	aKind := aType.Kind()
	bKind := bType.Kind()

	if aKind == reflect.Map || aKind == reflect.Slice || aKind == reflect.Array || aKind == reflect.String {
		if aType != bType {
			return false
		}

		if aKind != bKind {
			return false
		}
	}

	switch aKind {
	case reflect.Map:
		if aVal.Len() != bVal.Len() {
			return false
		}

		iter := aVal.MapRange()
		for iter.Next() {
			k := iter.Key()
			v1 := iter.Value()
			v2 := bVal.MapIndex(k)
			if !v2.IsValid() {
				return false
			}
			if !compareMap(v1.Interface(), v2.Interface()) {
				return false
			}
		}

	case reflect.Slice, reflect.Array:
		if aVal.Len() != bVal.Len() {
			return false
		}
		for i := 0; i < aVal.Len(); i++ {
			if !compareMap(aVal.Index(i).Interface(), bVal.Index(i).Interface()) {
				return false
			}
		}
	case reflect.String:
		return aVal.String() == bVal.String()

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return toFloat64(aVal) == toFloat64(bVal)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return toFloat64(aVal) == toFloat64(bVal)

	case reflect.Float32, reflect.Float64:
		return toFloat64(aVal) == toFloat64(bVal)

	case reflect.Complex64, reflect.Complex128:
		return aVal.Complex() == bVal.Complex()

	default:
		return false
	}
	return true
}

func TestComplexDecode(t *testing.T) {
	var err error

	// test decode into map
	decodeData := make(map[string]interface{})
	err = jsonPack.Decode("complex", testdata.ComplexExpData, &decodeData)
	if err != nil {
		t.Errorf("Decode fail, err: %v", err)
	}

	if !compareMap(decodeData, testdata.ComplexData) {
		t.Errorf("Decode fail, compareMap fail")
	}

	m1, _ := json.Marshal(decodeData)
	m2, _ := json.Marshal(testdata.ComplexData)
	if !reflect.DeepEqual(m1, m2) {
		t.Errorf("Decode fail, marshled data mismatch")
	}

	// test decode into struct
	decodeStructData := testdata.Complex{}
	err = jsonPack.Decode("complex", testdata.ComplexExpData, &decodeStructData)
	if err != nil {
		t.Errorf("Decode struct fail, err: %v", err)
	}
	// t.Logf("decodeStructData: %+v", decodeStructData)

	if !reflect.DeepEqual(&decodeStructData, &testdata.ComplexStructData) {
		t.Errorf("Decode struct data fail")
	}
}

func compareBytes(t *testing.T, a []byte, b []byte) bool {
	if len(a) != len(b) {
		t.Errorf("byte length mismatch, %v : %v", len(a), len(b))
		return false
	}
	for i, v := range a {
		if v != b[i] {
			t.Errorf("byte value mismatch at index %d", i)
			return false
		}
	}
	return true
}
