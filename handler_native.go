package jsonpack

import (
	"unsafe"

	ibuf "github.com/arloliu/jsonpack/buffer"
	"github.com/pkg/errors"
)

type stringOp struct{}

func (p *stringOp) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	str := buf.ReadString()
	*((*string)(ptr)) = str
	// *((*string)(ptr)) = buf.ReadString()
	return nil
}
func (p *stringOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	d := *((*string)(ptr))
	buf.WriteString(&d)

	return nil
}
func (p *stringOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	d, ok := data.(string)
	if !ok {
		return errors.WithStack(&TypeAssertionError{data, "string"})
	}
	buf.WriteString(&d)

	return nil
}
func (p *stringOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return buf.ReadString(), nil
}

type booleanOp struct{}

func (p *booleanOp) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	d := buf.ReadByte()
	*((*bool)(ptr)) = (d != 0)
	return nil
}
func (p *booleanOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	d := *((*bool)(ptr))
	if d {
		buf.WriteByte(byte(1))
	} else {
		buf.WriteByte(byte(0))
	}
	return nil
}
func (p *booleanOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	d := data.(bool)
	if d {
		buf.WriteByte(byte(1))
	} else {
		buf.WriteByte(byte(0))
	}
	return nil
}
func (p *booleanOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	d := buf.ReadByte()
	return d != 0, nil
}

/*
 * Signed interger handlers
 */
type int8Op struct{}

func (p *int8Op) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	*((*int8)(ptr)) = buf.ReadInt8()
	return nil
}

func (p *int8Op) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	val := *((*byte)(ptr))
	buf.WriteByte(val)
	return nil
}
func (p *int8Op) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	var val byte
	switch data := data.(type) {
	case float64:
		val = byte(data)
	case int8:
		val = byte(data)
	default:
		return errors.WithStack(&TypeAssertionError{data, "int8"})
	}

	buf.WriteByte(byte(val))
	return nil
}
func (p *int8Op) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return buf.ReadInt8(), nil
}

type int16LEOp struct{}

func (p *int16LEOp) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	*((*int16)(ptr)) = buf.ReadInt16LE()
	return nil
}

func (p *int16LEOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	val := *((*int16)(ptr))
	buf.WriteInt16LE(val)
	return nil

}
func (p *int16LEOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	var val int16
	switch data := data.(type) {
	case int16:
		val = data
	case float64:
		val = int16(data)
	default:
		return errors.WithStack(&TypeAssertionError{data, "int16"})
	}
	buf.WriteInt16LE(val)
	return nil
}
func (p *int16LEOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return buf.ReadInt16LE(), nil
}

type int16BEOp struct{}

func (p *int16BEOp) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	*((*int16)(ptr)) = buf.ReadInt16BE()
	return nil
}

func (p *int16BEOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	val := *((*int16)(ptr))
	buf.WriteInt16BE(val)
	return nil
}
func (p *int16BEOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	var val int16
	switch data := data.(type) {
	case int16:
		val = data
	case float64:
		val = int16(data)
	default:
		return errors.WithStack(&TypeAssertionError{data, "int16"})
	}
	buf.WriteInt16BE(val)
	return nil
}
func (p *int16BEOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return buf.ReadInt16BE(), nil
}

type int32LEOp struct{}

func (p *int32LEOp) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	*((*int32)(ptr)) = buf.ReadInt32LE()
	return nil
}

func (p *int32LEOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	val := *((*int32)(ptr))
	buf.WriteInt32LE(val)
	return nil
}
func (p *int32LEOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	var val int32
	switch data := data.(type) {
	case int32:
		val = data
	case float64:
		val = int32(data)
	default:
		return errors.WithStack(&TypeAssertionError{data, "int32"})
	}
	buf.WriteInt32LE(val)
	return nil
}
func (p *int32LEOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return buf.ReadInt32LE(), nil
}

type int32BEOp struct{}

func (p *int32BEOp) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	*((*int32)(ptr)) = buf.ReadInt32BE()
	return nil
}

func (p *int32BEOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	val := *((*int32)(ptr))
	buf.WriteInt32BE(val)
	return nil
}
func (p *int32BEOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	var val int32
	switch data := data.(type) {
	case int32:
		val = data
	case float64:
		val = int32(data)
	default:
		return errors.WithStack(&TypeAssertionError{data, "int32"})
	}
	buf.WriteInt32BE(val)
	return nil
}
func (p *int32BEOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return buf.ReadInt32BE(), nil
}

type int64LEOp struct{}

func (p *int64LEOp) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	*((*int64)(ptr)) = buf.ReadInt64LE()
	return nil
}
func (p *int64LEOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	val := *((*int64)(ptr))
	buf.WriteInt64LE(val)
	return nil
}
func (p *int64LEOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	var val int64
	switch data := data.(type) {
	case int64:
		val = data
	case float64:
		val = int64(data)
	default:
		return errors.WithStack(&TypeAssertionError{data, "int64"})
	}
	buf.WriteInt64LE(val)
	return nil
}
func (p *int64LEOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return buf.ReadInt64LE(), nil
}

type int64BEOp struct{}

func (p *int64BEOp) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	*((*int64)(ptr)) = buf.ReadInt64BE()
	return nil
}
func (p *int64BEOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	val := *((*int64)(ptr))
	buf.WriteInt64BE(val)
	return nil
}
func (p *int64BEOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	var val int64
	switch data := data.(type) {
	case int64:
		val = data
	case float64:
		val = int64(data)
	default:
		return errors.WithStack(&TypeAssertionError{data, "int64"})
	}
	buf.WriteInt64BE(val)
	return nil
}
func (p *int64BEOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return buf.ReadInt64BE(), nil
}

/*
 * Unsigned interger handlers
 */
type uint8Op struct{}

func (p *uint8Op) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	*((*uint8)(ptr)) = buf.ReadUint8()
	return nil
}
func (p *uint8Op) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	val := *((*byte)(ptr))
	buf.WriteByte(val)
	return nil
}
func (p *uint8Op) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	var val byte
	switch data := data.(type) {
	case uint8:
		val = byte(data)
	case float64:
		val = byte(data)
	default:
		return errors.WithStack(&TypeAssertionError{data, "uint8"})
	}
	buf.WriteByte(val)
	return nil
}
func (p *uint8Op) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return buf.ReadUint8(), nil
}

type uint16LEOp struct{}

func (p *uint16LEOp) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	*((*uint16)(ptr)) = buf.ReadUint16LE()
	return nil
}
func (p *uint16LEOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	val := *((*uint16)(ptr))
	buf.WriteUint16LE(val)
	return nil
}
func (p *uint16LEOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	var val uint16
	switch data := data.(type) {
	case uint16:
		val = uint16(data)
	case float64:
		val = uint16(data)
	default:
		return errors.WithStack(&TypeAssertionError{data, "uint16"})
	}
	buf.WriteUint16LE(val)
	return nil
}
func (p *uint16LEOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return buf.ReadUint16LE(), nil
}

type uint16BEOp struct{}

func (p *uint16BEOp) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	*((*uint16)(ptr)) = buf.ReadUint16BE()
	return nil
}
func (p *uint16BEOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	val := *((*uint16)(ptr))
	buf.WriteUint16BE(val)
	return nil
}
func (p *uint16BEOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	var val uint16
	switch data := data.(type) {
	case uint16:
		val = uint16(data)
	case float64:
		val = uint16(data)
	default:
		return errors.WithStack(&TypeAssertionError{data, "uint16"})
	}
	buf.WriteUint16BE(val)
	return nil
}
func (p *uint16BEOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return buf.ReadUint16BE(), nil
}

type uint32LEOp struct{}

func (p *uint32LEOp) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	*((*uint32)(ptr)) = buf.ReadUint32LE()
	return nil
}
func (p *uint32LEOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	val := *((*uint32)(ptr))
	buf.WriteUint32LE(val)
	return nil
}
func (p *uint32LEOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	var val uint32
	switch data := data.(type) {
	case uint32:
		val = uint32(data)
	case float64:
		val = uint32(data)
	default:
		return errors.WithStack(&TypeAssertionError{data, "uint32"})
	}
	buf.WriteUint32LE(val)
	return nil
}
func (p *uint32LEOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return buf.ReadUint32LE(), nil
}

type uint32BEOp struct{}

func (p *uint32BEOp) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	*((*uint32)(ptr)) = buf.ReadUint32BE()
	return nil
}

func (p *uint32BEOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	val := *((*uint32)(ptr))
	buf.WriteUint32BE(val)
	return nil
}
func (p *uint32BEOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	var val uint32
	switch data := data.(type) {
	case uint32:
		val = uint32(data)
	case float64:
		val = uint32(data)
	default:
		return errors.WithStack(&TypeAssertionError{data, "uint32"})
	}
	buf.WriteUint32BE(val)
	return nil
}
func (p *uint32BEOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return buf.ReadUint32BE(), nil
}

type uint64LEOp struct{}

func (p *uint64LEOp) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	*((*uint64)(ptr)) = buf.ReadUint64LE()
	return nil
}

func (p *uint64LEOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	val := *((*uint64)(ptr))
	buf.WriteUint64LE(val)
	return nil
}
func (p *uint64LEOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	var val uint64
	switch data := data.(type) {
	case uint64:
		val = uint64(data)
	case float64:
		val = uint64(data)
	default:
		return errors.WithStack(&TypeAssertionError{data, "uint64"})
	}
	buf.WriteUint64LE(val)
	return nil
}
func (p *uint64LEOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return buf.ReadUint64LE(), nil
}

type uint64BEOp struct{}

func (p *uint64BEOp) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	*((*uint64)(ptr)) = buf.ReadUint64BE()
	return nil
}
func (p *uint64BEOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	val := *((*uint64)(ptr))
	buf.WriteUint64BE(val)
	return nil
}
func (p *uint64BEOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	var val uint64
	switch data := data.(type) {
	case uint64:
		val = uint64(data)
	case float64:
		val = uint64(data)
	default:
		return errors.WithStack(&TypeAssertionError{data, "uint64"})
	}
	buf.WriteUint64BE(val)
	return nil
}
func (p *uint64BEOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return buf.ReadUint64BE(), nil
}

/*
 * floating number handlers
 */
type float32LEOp struct{}

func (p *float32LEOp) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	*((*float32)(ptr)) = buf.ReadFloat32LE()
	return nil
}
func (p *float32LEOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	val := *((*float32)(ptr))
	buf.WriteFloat32LE(val)
	return nil
}
func (p *float32LEOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	var val float32
	switch data := data.(type) {
	case float32:
		val = float32(data)
	case float64:
		val = float32(data)
	default:
		return errors.WithStack(&TypeAssertionError{data, "float32"})
	}
	buf.WriteFloat32LE(val)
	return nil
}
func (p *float32LEOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return buf.ReadFloat32LE(), nil
}

type float32BEOp struct{}

func (p *float32BEOp) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	*((*float32)(ptr)) = buf.ReadFloat32BE()
	return nil
}

func (p *float32BEOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	val := *((*float32)(ptr))
	buf.WriteFloat32BE(val)
	return nil
}
func (p *float32BEOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	var val float32
	switch data := data.(type) {
	case float32:
		val = float32(data)
	case float64:
		val = float32(data)
	default:
		return errors.WithStack(&TypeAssertionError{data, "float32"})
	}
	buf.WriteFloat32BE(val)
	return nil
}
func (p *float32BEOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return buf.ReadFloat32BE(), nil
}

type float64LEOp struct{}

func (p *float64LEOp) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	*((*float64)(ptr)) = buf.ReadFloat64LE()
	return nil
}

func (p *float64LEOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	val := *((*float64)(ptr))
	buf.WriteFloat64LE(val)
	return nil
}
func (p *float64LEOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	var val float64
	switch data := data.(type) {
	case float64:
		val = data
	default:
		return errors.WithStack(&TypeAssertionError{data, "float64"})
	}
	buf.WriteFloat64LE(val)
	return nil
}
func (p *float64LEOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return buf.ReadFloat64LE(), nil
}

type float64BEOp struct{}

func (p *float64BEOp) decodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	*((*float64)(ptr)) = buf.ReadFloat64BE()
	return nil
}

func (p *float64BEOp) encodeStruct(buf *ibuf.Buffer, opNode *structOperation, ptr unsafe.Pointer) error {
	val := *((*float64)(ptr))
	buf.WriteFloat64BE(val)
	return nil
}
func (p *float64BEOp) encodeDynamic(buf *ibuf.Buffer, opNode *operation, data interface{}) error {
	var val float64
	switch data := data.(type) {
	case float64:
		val = data
	default:
		return errors.WithStack(&TypeAssertionError{data, "float64"})
	}
	buf.WriteFloat64BE(val)
	return nil
}
func (p *float64BEOp) decodeDynamic(buf *ibuf.Buffer, opNode *operation, v interface{}) (interface{}, error) {
	return buf.ReadFloat64BE(), nil
}
