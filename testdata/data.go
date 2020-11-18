package testdata

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"runtime"
)

var ComplexSchDef []byte
var ComplexAnyData interface{}
var ComplexData map[string]interface{}
var ComplexRawData []byte
var ComplexStructData Complex
var ComplexExpData []byte
var ComplexPbExpData []byte

var StructSchDef []byte
var StructData TestStruct
var StructExpData []byte
var StructPbExpData []byte

var TypesSchDef []byte
var TypesRawData []byte
var TypesAnyData interface{}
var TypesMapData map[string]interface{}
var TypesStructData Types
var TypesExpData []byte

var SliceSchDef []byte
var SliceRawData []byte
var SliceData []interface{}
var SliceMapData []map[string]interface{}
var SliceStructData []TestArrayStruct
var SliceExpData []byte

var ArrayData [3]interface{}
var ArrayMapData [3]map[string]interface{}
var ArrayStructData [3]TestArrayStruct

var dataDir string

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dataDir = path.Dir(filename)

	var err error

	TypesSchDef, _, err = loadJsonTextData("types.def.json")
	if err != nil {
		os.Exit(1)
	}
	TypesRawData, TypesMapData, err = loadJsonTextData("types.data.json")
	if err != nil {
		os.Exit(1)
	}

	_, TypesAnyData, err = loadJsonTextAnyData("types.data.json")
	if err != nil {
		os.Exit(1)
	}

	err = json.Unmarshal(TypesRawData, &TypesStructData)
	if err != nil {
		os.Exit(1)
	}

	TypesExpData, err = loadRawTestData("types.bin")
	if err != nil {
		os.Exit(1)
	}

	ComplexSchDef, _, err = loadJsonTextData("complex.def.json")
	if err != nil {
		os.Exit(1)
	}

	_, ComplexAnyData, err = loadJsonTextAnyData("complex.data.json")
	if err != nil {
		os.Exit(1)
	}

	ComplexRawData, ComplexData, err = loadJsonTextData("complex.data.json")
	if err != nil {
		os.Exit(1)
	}

	err = json.Unmarshal(ComplexRawData, &ComplexStructData)
	if err != nil {
		os.Exit(1)
	}

	ComplexExpData, err = loadRawTestData("complex.bin")
	if err != nil {
		os.Exit(1)
	}

	ComplexPbExpData, err = loadRawTestData("complex.pb.bin")
	if err != nil {
		os.Exit(1)
	}

	SliceSchDef, _, err = loadJsonTextData("array.def.json")
	if err != nil {
		os.Exit(1)
	}

	SliceData, SliceRawData, err = loadTestArrayJSON("array.data.json")
	if err != nil {
		os.Exit(1)
	}

	SliceMapData, _, err = loadTestArrayMapJSON("array.data.json")
	if err != nil {
		os.Exit(1)
	}

	err = json.Unmarshal(SliceRawData, &SliceStructData)
	if err != nil {
		os.Exit(1)
	}

	// clone slice to array
	for i, v := range SliceData {
		ArrayData[i] = v
	}
	for i, v := range SliceMapData {
		ArrayMapData[i] = v
	}
	for i, v := range SliceStructData {
		ArrayStructData[i] = v
	}

	SliceExpData, err = loadRawTestData("array.bin")
	if err != nil {
		os.Exit(1)
	}

	StructSchDef, _, err = loadJsonTextData("struct.def.json")
	if err != nil {
		os.Exit(1)
	}

	structRawData, _, err := loadJsonTextData("struct.data.json")
	if err != nil {
		os.Exit(1)
	}
	err = json.Unmarshal(structRawData, &StructData)
	if err != nil {
		os.Exit(1)
	}

	StructExpData, err = loadRawTestData("struct.bin")
	if err != nil {
		os.Exit(1)
	}

	StructPbExpData, err = loadRawTestData("struct.pb.bin")
	if err != nil {
		os.Exit(1)
	}
}

func loadRawTestData(filename string) ([]byte, error) {
	return ioutil.ReadFile(path.Join(dataDir, filename))
}

func loadJsonTextData(filename string) ([]byte, map[string]interface{}, error) {
	data, err := loadRawTestData(filename)
	if err != nil {
		return nil, nil, err
	}
	jsonData := make(map[string]interface{})
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		return nil, nil, err
	}

	rawData, err := json.Marshal(jsonData)
	if err != nil {
		return nil, nil, err
	}
	return rawData, jsonData, nil
}

func loadJsonTextAnyData(filename string) ([]byte, interface{}, error) {
	data, err := loadRawTestData(filename)
	if err != nil {
		return nil, nil, err
	}
	var jsonData interface{}
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		return nil, nil, err
	}

	rawData, err := json.Marshal(jsonData)
	if err != nil {
		return nil, nil, err
	}
	return rawData, jsonData, nil
}

func loadTestArrayJSON(filename string) ([]interface{}, []byte, error) {
	data, err := loadRawTestData(filename)
	if err != nil {
		return nil, nil, err
	}
	jsonData := make([]interface{}, 0)
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		return nil, data, err
	}
	return jsonData, data, nil
}

func loadTestArrayMapJSON(filename string) ([]map[string]interface{}, []byte, error) {
	data, err := loadRawTestData(filename)
	if err != nil {
		return nil, nil, err
	}
	jsonData := make([]map[string]interface{}, 0)
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		return nil, data, err
	}
	return jsonData, data, nil
}
