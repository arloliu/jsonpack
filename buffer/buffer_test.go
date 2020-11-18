package buffer

import (
	"math"
	"reflect"
	"testing"
)

func TestFrom(t *testing.T) {
	buf := Create(1)
	str1 := "test utf-8 中文 string"
	str2 := "test2 utf-8 中文 string"
	buf.WriteString(&str1)
	buf.WriteString(&str2)

	buf2 := From(buf.Bytes())
	buf2.Seek(0, false)

	if buf2.ReadString() != str1 {
		t.Errorf("fail")
	}
	if buf2.ReadString() != str2 {
		t.Errorf("fail")
	}
}
func TestVarInt(t *testing.T) {
	buf := Create(10)
	nums := [...]int64{0, 127, 128, 16384, math.MinInt64, math.MaxInt64}

	for _, num := range nums {
		n := ByteLenVarInt(num)
		buf.WriteVarInt(num)
		buf.Seek(-n, true)
		v, n2 := buf.ReadVarInt()
		if v != num || int64(n2) != n {
			t.Errorf("Write/Read varint fail, expected %d, got %d, len: %d", num, v, n)
		}
	}
}
func TestVarUint(t *testing.T) {
	buf := Create(10)
	nums := [...]uint64{0, 127, 128, 16384, math.MaxUint64}

	for _, num := range nums {
		n := ByteLenVarUint(num)
		buf.WriteVarUint(num)
		buf.Seek(-n, true)
		v, n2 := buf.ReadVarUint()
		if v != num || n != int64(n2) {
			t.Errorf("Write/Read varuint fail, expect %d, got %d, len: %d", num, v, n)
		}
	}
}

func TestString(t *testing.T) {
	buf := Create(128)
	str := "test utf-8 中文 string"
	off := buf.Offset()
	buf.WriteString(&str)
	buf.Seek(off, false)
	result := buf.ReadString()

	if result != str {
		t.Errorf("Write/Read string fail, expect %s, got %s", str, result)
	}
}

func TestStringPtr(t *testing.T) {
	buf := Create(128)
	str := "test utf-8 中文 string"
	off := buf.Offset()
	buf.WriteString(&str)
	buf.Seek(off, false)
	result := buf.ReadStringPtr()

	if *result != str {
		t.Errorf("Write/Read string pointer fail, expect %s, got %s", str, *result)
	}
}

func TestFloatTypes(t *testing.T) {
	buf := Create(1)
	inFloat32 := []float32{1.3, math.MaxFloat32}
	inFloat64 := []float64{1.3, 123141.32, math.MaxFloat64}

	var outFloat32 []float32 = make([]float32, len(inFloat32))
	var outFloat64 []float64 = make([]float64, len(inFloat64))

	// little-endian
	for _, v := range inFloat32 {
		buf.WriteFloat32LE(v)
	}
	buf.Seek(int64(-4*len(inFloat32)), true)
	for i := 0; i < len(inFloat32); i++ {
		outFloat32[i] = buf.ReadFloat32LE()
	}
	if reflect.DeepEqual(inFloat32, outFloat32) == false {
		t.Errorf("Write/Read float32le fail, expect %v, got %v", inFloat32, outFloat32)
	}

	for _, v := range inFloat64 {
		buf.WriteFloat64LE(v)
	}
	buf.Seek(int64(-8*len(inFloat64)), true)
	for i := 0; i < len(inFloat64); i++ {
		outFloat64[i] = buf.ReadFloat64LE()
	}
	if reflect.DeepEqual(inFloat64, outFloat64) == false {
		t.Errorf("Write/Read float64le fail, expect %v, got %v", inFloat64, outFloat64)
	}

	// big-endian
	for _, v := range inFloat32 {
		buf.WriteFloat32BE(v)
	}
	buf.Seek(int64(-4*len(inFloat32)), true)
	for i := 0; i < len(inFloat32); i++ {
		outFloat32[i] = buf.ReadFloat32BE()
	}
	if reflect.DeepEqual(inFloat32, outFloat32) == false {
		t.Errorf("Write/Read float32BE fail, expect %v, got %v", inFloat32, outFloat32)
	}

	for _, v := range inFloat64 {
		buf.WriteFloat64BE(v)
	}
	buf.Seek(int64(-8*len(inFloat64)), true)
	for i := 0; i < len(inFloat64); i++ {
		outFloat64[i] = buf.ReadFloat64BE()
	}
	if reflect.DeepEqual(inFloat64, outFloat64) == false {
		t.Errorf("Write/Read float64BE fail, expect %v, got %v", inFloat64, outFloat64)
	}

}
