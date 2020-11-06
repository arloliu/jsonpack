package testdata

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"runtime"
)

var ComplexSchDef []byte
var ComplexData map[string]interface{}
var ComplexRawData []byte
var ComplexStructData Complex
var ComplexExpData []byte
var ComplexPbExpData []byte

var StructSchDef []byte
var StructData TestStruct
var StructExpData []byte
var StructPbExpData []byte

var ArraySchDef []byte
var ArrayRawData []byte
var ArrayData []interface{}
var ArrayStructData []TestArrayStruct
var ArrayExpData []byte

var dataDir string

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dataDir = path.Dir(filename)

	var err error
	ComplexSchDef, _, err = loadJsonTextData("complex.def.json")
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

	ArraySchDef, _, err = loadJsonTextData("array.def.json")
	if err != nil {
		os.Exit(1)
	}
	ArrayData, ArrayRawData, err = loadTestArrayJSON("array.data.json")
	if err != nil {
		os.Exit(1)
	}

	err = json.Unmarshal(ArrayRawData, &ArrayStructData)
	if err != nil {
		os.Exit(1)
	}

	ArrayExpData, err = loadRawTestData("array.bin")
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
