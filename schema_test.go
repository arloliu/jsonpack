package jsonpack

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/pkg/errors"

	"github.com/arloliu/jsonpack/testdata"
)

var jsonPack *JSONPack

func init() {
	var err error
	jsonPack = NewJSONPack()

	_, err = jsonPack.AddSchema("types", testdata.TypesSchDef)
	if err != nil {
		fmt.Printf("AddSchema types err: %+v\n", err)
		os.Exit(1)
	}

	_, err = jsonPack.AddSchema("complex", testdata.ComplexSchDef)
	if err != nil {
		fmt.Printf("AddSchema complex err: %+v\n", err)
		os.Exit(1)
	}

	_, err = jsonPack.AddSchema("sliceObject", testdata.SliceSchDef)
	if err != nil {
		fmt.Printf("AddSchema sliceObject err: %+v\n", err)
		os.Exit(1)
	}
	_, err = jsonPack.AddSchema("testStruct", testdata.StructSchDef)
	if err != nil {
		fmt.Printf("AddSchema testStruct err: %+v\n", err)
		os.Exit(1)
	}
}

func TestTypesEncode(t *testing.T) {
	var err error
	var encData []byte

	// test encode map[string]interface{}
	encData, err = jsonPack.Encode("types", testdata.TypesMapData)
	if err != nil {
		t.Errorf("Encode fail, error: %+v", err)
	}
	if !compareBytes(t, encData, testdata.TypesExpData) {
		t.Errorf("Encoded data mismatch")
	}

	// test encode interface{}
	encData, err = jsonPack.Encode("types", testdata.TypesAnyData)
	if err != nil {
		t.Errorf("Encode fail, error: %+v", err)
	}
	if !compareBytes(t, encData, testdata.TypesExpData) {
		t.Errorf("Encoded data mismatch")
	}
}

func TestTypesDecode(t *testing.T) {
	var err error
	decodeData := testdata.Types{}

	err = jsonPack.Decode("types", testdata.TypesExpData, &decodeData)
	if err != nil {
		t.Errorf("Decode types fail, err: %+v", err)
	}
	// t.Logf("decodeData: %+v", decodeData)
	// t.Logf("TypesStructData: %+v", testdata.TypesStructData)
	if !reflect.DeepEqual(&testdata.TypesStructData, &decodeData) {
		t.Errorf("Decode types data mismatch")
	}
}

func TestComplexAddSchemaFromStruct(t *testing.T) {
	_, err := jsonPack.AddSchema("complex_from_struct", testdata.Complex{})
	if err != nil {
		t.Errorf("AddSchema fail, error: %+v", err)
	}
}

func TestStructEncode(t *testing.T) {
	encData, err := jsonPack.Encode("testStruct", &testdata.StructData)
	if err != nil {
		t.Errorf("Encode fail, error: %+v", err)
	}

	if !compareBytes(t, encData, testdata.StructExpData) {
		t.Errorf("Encoded data mismatch")
	}

}

func TestComplexEncode(t *testing.T) {
	var err error
	var encData []byte

	// test encode interface{}
	encData, err = jsonPack.Encode("complex", testdata.ComplexAnyData)
	if err != nil {
		t.Errorf("Encode fail, error: %+v", err)
	}
	if !compareBytes(t, encData, testdata.ComplexExpData) {
		t.Errorf("Encoded data mismatch")
	}

	// test encode *interface{}
	encData, err = jsonPack.Encode("complex", &testdata.ComplexAnyData)
	if err != nil {
		t.Errorf("Encode fail, error: %+v", err)
	}
	if !compareBytes(t, encData, testdata.ComplexExpData) {
		t.Errorf("Encoded data mismatch")
	}

	// test encode map[string]interface{}
	encData, err = jsonPack.Encode("complex", testdata.ComplexData)
	if err != nil {
		t.Errorf("Encode fail, error: %+v", err)
	}
	if !compareBytes(t, encData, testdata.ComplexExpData) {
		t.Errorf("Encoded data mismatch")
	}

	// test encode *map[string]interface{}
	encData, err = jsonPack.Encode("complex", &testdata.ComplexData)
	if err != nil {
		t.Errorf("Encode fail, error: %+v", err)
	}
	if !compareBytes(t, encData, testdata.ComplexExpData) {
		t.Errorf("Encoded data mismatch")
	}

	var encData2 []byte
	// test encode Complex
	encData2, err = jsonPack.Encode("complex", testdata.ComplexStructData)
	if err != nil {
		t.Errorf("Encode struct fail, error: %+v", err)
	}
	if !compareBytes(t, encData2, testdata.ComplexExpData) {
		t.Errorf("Encoded complex data mismatch")
	}

	// test encode *Complex
	encData2, err = jsonPack.Encode("complex", &testdata.ComplexStructData)
	if err != nil {
		t.Errorf("Encode struct fail, error: %+v", err)
	}
	if !compareBytes(t, encData2, testdata.ComplexExpData) {
		t.Errorf("Encoded complex data mismatch")
	}

	_, err = jsonPack.AddSchema("complex_from_struct", testdata.Complex{})
	if err != nil {
		t.Errorf("AddSchema fail, error: %+v", err)
	}

	var encData3 []byte
	encData3, err = jsonPack.Encode("complex_from_struct", &testdata.ComplexStructData)
	if err != nil {
		t.Errorf("Encode struct fail, error: %+v", err)
	}
	if !compareBytes(t, encData3, testdata.ComplexExpData) {
		t.Errorf("Encoded complex data mismatch")
	}

	var encData4 []byte = make([]byte, 128)
	err = jsonPack.EncodeTo("complex", &testdata.ComplexStructData, &encData4)
	if err != nil {
		t.Errorf("Encode struct fail, error: %+v", err)
	}
	if !compareBytes(t, encData4, testdata.ComplexExpData) {
		t.Errorf("Encoded complex data mismatch")
	}

}

func TestSliceEncode(t *testing.T) {
	var err error
	var encData []byte

	// test encode []interface{}
	encData, err = jsonPack.Encode("sliceObject", testdata.SliceData)
	if err != nil {
		t.Errorf("Encode fail, error: %+v", err)
	}
	if !compareBytes(t, encData, testdata.SliceExpData) {
		t.Errorf("Encoded data mismatch")
	}

	// test encode *[]interface{}
	encData, err = jsonPack.Encode("sliceObject", &testdata.SliceData)
	if err != nil {
		t.Errorf("Encode fail, error: %+v", err)
	}
	if !compareBytes(t, encData, testdata.SliceExpData) {
		t.Errorf("Encoded data mismatch")
	}

	// test encode []map[string]interface{}
	encData, err = jsonPack.Encode("sliceObject", testdata.SliceMapData)
	if err != nil {
		t.Errorf("Encode fail, error: %+v", err)
	}
	if !compareBytes(t, encData, testdata.SliceExpData) {
		t.Errorf("Encoded data mismatch")
	}

	// test encode *[]map[string]interface{}
	encData, err = jsonPack.Encode("sliceObject", &testdata.SliceMapData)
	if err != nil {
		t.Errorf("Encode fail, error: %+v", err)
	}
	if !compareBytes(t, encData, testdata.SliceExpData) {
		t.Errorf("Encoded data mismatch")
	}

	// test encode []TestArrayStruct
	encData, err = jsonPack.Encode("sliceObject", testdata.SliceStructData)
	if err != nil {
		t.Errorf("Encode array struct fail, error: %+v", err)
	}
	if !compareBytes(t, encData, testdata.SliceExpData) {
		t.Errorf("Encoded data mismatch")
	}

	// test encode *[]TestArrayStruct
	encData, err = jsonPack.Encode("sliceObject", &testdata.SliceStructData)
	if err != nil {
		t.Errorf("Encode array struct fail, error: %+v", err)
	}
	if !compareBytes(t, encData, testdata.SliceExpData) {
		t.Errorf("Encoded data mismatch")
	}
}

func TestArrayEncode(t *testing.T) {
	var err error
	var encData []byte

	// test encode [3]interface{}
	encData, err = jsonPack.Encode("sliceObject", testdata.ArrayData)
	if err != nil {
		t.Errorf("Encode fail, error: %+v", err)
	}
	if !compareBytes(t, encData, testdata.SliceExpData) {
		t.Errorf("Encoded data mismatch")
	}

	// test encode *[3]interface{}
	encData, err = jsonPack.Encode("sliceObject", &testdata.ArrayData)
	if err != nil {
		t.Errorf("Encode fail, error: %+v", err)
	}
	if !compareBytes(t, encData, testdata.SliceExpData) {
		t.Errorf("Encoded data mismatch")
	}

	// test encode [3]map[string]interface{}
	encData, err = jsonPack.Encode("sliceObject", testdata.ArrayMapData)
	if err != nil {
		t.Errorf("Encode fail, error: %+v", err)
	}
	if !compareBytes(t, encData, testdata.SliceExpData) {
		t.Errorf("Encoded data mismatch")
	}

	// test encode *[3]map[string]interface{}
	encData, err = jsonPack.Encode("sliceObject", &testdata.ArrayMapData)
	if err != nil {
		t.Errorf("Encode fail, error: %+v", err)
	}
	if !compareBytes(t, encData, testdata.SliceExpData) {
		t.Errorf("Encoded data mismatch")
	}
}

func TestSliceDecode(t *testing.T) {
	var err error

	var decodeData []interface{}
	err = jsonPack.Decode("sliceObject", testdata.SliceExpData, &decodeData)
	if err != nil {
		t.Errorf("Decode fail, err: %+v", err)
	}
	if !compareMap(decodeData, testdata.SliceData) {
		t.Errorf("Decode fail, compareMap fail")
	}

	var decodeData2 []map[string]interface{}
	err = jsonPack.Decode("sliceObject", testdata.SliceExpData, &decodeData2)
	if err != nil {
		t.Errorf("Decode fail, err: %+v", err)
	}
	if !compareMap(decodeData2, testdata.SliceMapData) {
		t.Errorf("Decode fail, compareMap fail")
	}

	// decodeStructData := make([]testdata.TestArrayStruct, 0)
	var decodeStructData []testdata.TestArrayStruct
	err = jsonPack.Decode("sliceObject", testdata.SliceExpData, &decodeStructData)
	if err != nil {
		t.Errorf("Decode array struct fail, err: %+v", err)
	}

	if !reflect.DeepEqual(&testdata.SliceStructData, &decodeStructData) {
		t.Errorf("Decode array struct data mismatch")
	}
}

func TestArrayDecode(t *testing.T) {
	var err error

	var decodeData [3]interface{}
	err = jsonPack.Decode("sliceObject", testdata.SliceExpData, &decodeData)
	if err != nil {
		t.Errorf("Decode fail, err: %+v", err)
	}
	// if !compareMap(decodeData, testdata.ArrayData) {
	// 	t.Errorf("Decode fail, compareMap fail")
	// }

	// decodeStructData := make([]testdata.TestArrayStruct, 0)
	var decodeStructData []testdata.TestArrayStruct
	err = jsonPack.Decode("sliceObject", testdata.SliceExpData, &decodeStructData)
	if err != nil {
		t.Errorf("Decode array struct fail, err: %+v", err)
	}

	if !reflect.DeepEqual(&testdata.SliceStructData, &decodeStructData) {
		t.Errorf("Decode array struct data mismatch")
	}
}

func TestStructDecode(t *testing.T) {
	var err error
	decodeData := testdata.TestStruct{}

	err = jsonPack.Decode("testStruct", testdata.StructExpData, &decodeData)
	if err != nil {
		t.Errorf("Decode struct fail, err: %+v", err)
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
		t.Errorf("Decode fail, err: %+v", err)
	}
	if !compareMap(decodeData, testdata.ComplexData) {
		t.Errorf("Decode fail, compareMap fail")
	}

	// test decode into interface{}
	var decodeData2 interface{}
	err = jsonPack.Decode("complex", testdata.ComplexExpData, &decodeData2)
	var expectErr *DecodeError
	if !errors.As(err, &expectErr) {
		t.Errorf("Decode fail, err: %+v", err)
	}

	// test decode into struct
	decodeStructData := testdata.Complex{}
	err = jsonPack.Decode("complex", testdata.ComplexExpData, &decodeStructData)
	if err != nil {
		t.Errorf("Decode struct fail, err: %+v", err)
	}
	// t.Logf("decodeStructData: %+v", decodeStructData)

	if !reflect.DeepEqual(&decodeStructData, &testdata.ComplexStructData) {
		t.Errorf("Decode struct data fail")
	}

	err = jsonPack.Decode("complex", testdata.ComplexExpData, &decodeStructData)
	if err != nil {
		t.Errorf("Decode struct fail, err: %+v", err)
	}
	// t.Logf("decodeStructData: %+v", decodeStructData)

	if !reflect.DeepEqual(&decodeStructData, &testdata.ComplexStructData) {
		t.Errorf("Decode struct data fail")
	}
}

func compareBytes(t *testing.T, a []byte, b []byte) bool {
	if len(a) != len(b) {
		t.Errorf("byte length mismatch, %+v : %+v", len(a), len(b))
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
