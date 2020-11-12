package jsonpack

import (
	"math"
	"reflect"
	"testing"

	ibuf "github.com/arloliu/jsonpack/buffer"

	"github.com/modern-go/reflect2"
)

func testNumber(t *testing.T, buf *ibuf.Buffer, data interface{}) {
	var o1 opHandler
	var o2 opHandler
	kind := reflect.TypeOf(data).Kind()
	switch data.(type) {
	case int8:
		o1 = &int8Op{}
	case int16:
		o1 = &int16LEOp{}
		o2 = &int16BEOp{}
	case int32:
		o1 = &int32LEOp{}
		o2 = &int32BEOp{}
	case int64:
		o1 = &int64LEOp{}
		o2 = &int64BEOp{}
	case uint8:
		o1 = &uint8Op{}
	case uint16:
		o1 = &uint16LEOp{}
		o2 = &uint16BEOp{}
	case uint32:
		o1 = &uint32LEOp{}
		o2 = &uint32BEOp{}
	case uint64:
		o1 = &uint64LEOp{}
		o2 = &uint64BEOp{}
	case float32:
		o1 = &float32LEOp{}
		o2 = &float32BEOp{}
	case float64:
		o1 = &float64LEOp{}
		o2 = &float64BEOp{}
	default:
		t.Errorf("Unsupported type: %s", kind)
	}

	buf.Seek(0, false)
	o1.encodeDynamic(buf, &operation{}, data)
	buf.Seek(0, false)
	v1, _ := o1.decodeDynamic(buf, &operation{}, data)
	if v1 != data {
		if o2 != nil {
			t.Errorf("%s operation(LE) write fail, expect %v, got %v", kind, data, v1)
		} else {
			t.Errorf("%s operation write fail, expect %v, got %v", kind, data, v1)
		}
	}

	buf.Seek(0, false)
	o1.encodeStruct(buf, &structOperation{}, reflect2.PtrOf(data))
	buf.Seek(0, false)
	v2, _ := o1.decodeDynamic(buf, &operation{}, data)
	if v2 != data {
		if o2 != nil {
			t.Errorf("%s operation(LE) encodeStruct fail, expect %v, got %v", kind, data, v1)
		} else {
			t.Errorf("%s operation encodeStruct fail, expect %v, got %v", kind, data, v1)
		}
	}

	if o2 != nil {
		buf.Seek(0, false)
		o2.encodeDynamic(buf, &operation{}, data)
		buf.Seek(0, false)
		if v3, _ := o2.decodeDynamic(buf, &operation{}, data); v3 != data {
			t.Errorf("%s operation(BE) write fail", kind)
		}
		buf.Seek(0, false)
		o2.encodeStruct(buf, &structOperation{}, reflect2.PtrOf(data))
		buf.Seek(0, false)
		if v4, _ := o2.decodeDynamic(buf, &operation{}, data); v4 != data {
			t.Errorf("%s operation(BE) encodeStruct fail", kind)
		}
	}
}
func TestOpHandler(t *testing.T) {
	buf := ibuf.Create(1)

	var o opHandler
	o = &stringOp{}

	testStr := "TestOpHandler test string"
	o.encodeDynamic(buf, &operation{}, testStr)
	buf.Seek(0, false)
	s, _ := o.decodeDynamic(buf, &operation{}, nil)
	if s.(string) != testStr {
		t.Errorf("string operation fail")
	}

	buf.Seek(0, false)
	o = &booleanOp{}
	o.encodeDynamic(buf, &operation{}, true)
	buf.Seek(0, false)
	v1, _ := o.decodeDynamic(buf, &operation{}, nil)
	if v1.(bool) != true {
		t.Errorf("boolean operation fail")
	}

	testNumber(t, buf, int8(127))
	testNumber(t, buf, int8(-128))
	testNumber(t, buf, int16(math.MaxInt16))
	testNumber(t, buf, int16(math.MinInt16))
	testNumber(t, buf, int32(math.MaxInt32))
	testNumber(t, buf, int32(math.MinInt32))
	testNumber(t, buf, int64(math.MaxInt64))
	testNumber(t, buf, int64(math.MinInt64))

	testNumber(t, buf, uint8(255))
	testNumber(t, buf, uint8(0))
	testNumber(t, buf, uint16(math.MaxUint16))
	testNumber(t, buf, uint16(0))
	testNumber(t, buf, uint32(math.MaxUint32))
	testNumber(t, buf, uint32(0))
	testNumber(t, buf, uint64(math.MaxUint64))
	testNumber(t, buf, uint64(0))

	testNumber(t, buf, float32(math.MaxFloat32))
	testNumber(t, buf, float32(-math.MaxFloat32))
	testNumber(t, buf, float64(math.MaxFloat64))
	testNumber(t, buf, float64(-math.MaxFloat64))
}

func BenchmarkOpHandlerTypeAssertion(b *testing.B) {
	op := newOperation("", _float32BEOp, float32BEOpType)
	var val int = 0
	for i := 0; i < b.N; i++ {
		switch op.handler.(type) {
		case *objectOp:
			val++
		case *arrayOp:
			val++
		case *float32BEOp:
			val++
		default:
		}
	}
}

func BenchmarkOpHandlerCompareType(b *testing.B) {
	op := newOperation("", _float32BEOp, float32BEOpType)
	var val int = 0
	for i := 0; i < b.N; i++ {
		switch op.handlerType {
		case objectOpType:
			val++
		case arrayOpType:
			val++
		case float32BEOpType:
			val++
		default:
		}
	}
}
