package benchmark

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/arloliu/jsonpack"

	"github.com/arloliu/jsonpack/testdata"

	jsoniter "github.com/json-iterator/go"
	"google.golang.org/protobuf/proto"
)

var jsonPacker *jsonpack.JSONPack

func init() {
	var err error
	jsonPacker = jsonpack.NewJSONPack()

	_, err = jsonPacker.AddSchema("complex", testdata.ComplexSchDef)
	if err != nil {
		os.Exit(1)
	}

	_, err = jsonPacker.AddSchema("testStruct", testdata.StructSchDef)
	if err != nil {
		os.Exit(1)
	}

}

func BenchmarkStruct_Protobuf_Encode(b *testing.B) {
	st := testdata.TestStructPb{}
	err := jsonPacker.Decode("testStruct", testdata.StructExpData, &st)
	if err != nil {
		b.Fatalf("decode testStruct into pb fail, err:%v\n", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		proto.Marshal(&st)
	}
}

func BenchmarkStruct_JSONPACK_Encode(b *testing.B) {
	sch := jsonPacker.GetSchema("testStruct")
	buf := make([]byte, 1024)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sch.EncodeTo(&testdata.StructData, &buf)
	}
}

func BenchmarkStruct_JSON_Encode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		json.Marshal(&testdata.StructData)
	}
}

func BenchmarkStruct_Jsoniter_Encode(b *testing.B) {
	j := jsoniter.ConfigDefault
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		j.Marshal(&testdata.StructData)
	}
}

func BenchmarkComplex_JSONPACK_Encode(b *testing.B) {
	var err error
	sch := jsonPacker.GetSchema("complex")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = sch.Encode(&testdata.ComplexData)
		if err != nil {
			b.FailNow()
		}
	}
}

func BenchmarkComplex_JSON_Encode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		json.Marshal(&testdata.ComplexData)
	}
}

func BenchmarkComplex_Jsoniter_Encode(b *testing.B) {
	j := jsoniter.ConfigDefault
	for i := 0; i < b.N; i++ {
		j.Marshal(&testdata.ComplexData)
	}
}

func BenchmarkComplex_Protobuf_Encode(b *testing.B) {
	var err error
	st := testdata.ComplextPb{}
	err = jsonPacker.Decode("complex", testdata.ComplexExpData, &st)
	if err != nil {
		b.Fatalf("decode complex into pb fail, err:%v\n", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = proto.Marshal(&st)
		if err != nil {
			b.FailNow()
		}
	}
}

func BenchmarkComplex_JSONPACK_Struct_Encode(b *testing.B) {
	sch := jsonPacker.GetSchema("complex")
	// buf := make([]byte, 4096)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = sch.Encode(&testdata.ComplexStructData)
		// _ = sch.EncodeTo(&testdata.ComplexStructData, &buf)
	}
}

func BenchmarkComplex_JSON_Struct_Encode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		json.Marshal(&testdata.ComplexStructData)
	}
}

func BenchmarkComplex_Jsoniter_Struct_Encode(b *testing.B) {
	j := jsoniter.ConfigDefault
	for i := 0; i < b.N; i++ {
		j.Marshal(&testdata.ComplexStructData)
	}
}

func BenchmarkComplex_JSONPACK_Decode(b *testing.B) {
	m := make(map[string]interface{})
	sch := jsonPacker.GetSchema("complex")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := sch.Decode(testdata.ComplexExpData, &m)
		if err != nil {
			b.Fatalf("decode fail, err: %v", err)
		}
	}
}

func BenchmarkComplex_JSON_Decode(b *testing.B) {
	m := make(map[string]interface{})
	for i := 0; i < b.N; i++ {
		json.Unmarshal(testdata.ComplexRawData, &m)
	}
}

func BenchmarkComplex_Jsoniter_Decode(b *testing.B) {
	j := jsoniter.ConfigDefault
	m := make(map[string]interface{})
	for i := 0; i < b.N; i++ {
		j.Unmarshal(testdata.ComplexRawData, &m)
	}
}

func BenchmarkComplex_Protobuf_Decode(b *testing.B) {
	s := testdata.ComplextPb{}
	for i := 0; i < b.N; i++ {
		err := proto.Unmarshal(testdata.ComplexPbExpData, &s)
		if err != nil {
			b.Fatalf("decode fail, err: %v", err)
		}
	}
}
func BenchmarkComplex_JSONPACK_Struct_Decode(b *testing.B) {
	s := testdata.Complex{}
	sch := jsonPacker.GetSchema("complex")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := sch.Decode(testdata.ComplexExpData, &s)
		if err != nil {
			b.Fatalf("decode fail, err: %v", err)
		}
	}
}

func BenchmarkComplex_JSON_Struct_Decode(b *testing.B) {
	s := testdata.Complex{}
	for i := 0; i < b.N; i++ {
		err := json.Unmarshal(testdata.ComplexRawData, &s)
		if err != nil {
			b.Fatalf("decode fail, err: %v", err)
		}
	}
}

func BenchmarkComplex_Jsoniter_Struct_Decode(b *testing.B) {
	j := jsoniter.ConfigDefault
	s := testdata.Complex{}
	for i := 0; i < b.N; i++ {
		err := j.Unmarshal(testdata.ComplexRawData, &s)
		if err != nil {
			b.Fatalf("decode fail, err: %v", err)
		}
	}
}

//nolint:deadcode,unused
func reportEncodeRatio(t *testing.T, name string, enc []byte, ori []byte) {
	var encodeRatio float64 = float64(len(enc)) * 100.0 / float64(len(ori))
	t.Logf("%s encode ratio %f%%", name, encodeRatio)
}
