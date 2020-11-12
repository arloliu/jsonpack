//nolint:gosec
package buffer

import "unsafe"

func (b *Buffer) WriteInt8(data int8) {
	b.WriteByte(byte(data))
	b.offset += 1
}

func (b *Buffer) WriteInt16LE(data int16) {
	b.EnsureCap(2)
	b.buf[b.offset] = byte(data)
	b.buf[b.offset+1] = byte(data >> 8)
	b.offset += 2
}

func (b *Buffer) WriteInt16BE(data int16) {
	b.EnsureCap(2)
	b.buf[b.offset] = byte(data >> 8)
	b.buf[b.offset+1] = byte(data)
	b.offset += 2
}

func (b *Buffer) WriteInt32LE(data int32) {
	b.EnsureCap(4)
	b.buf[b.offset] = byte(data)
	b.buf[b.offset+1] = byte(data >> 8)
	b.buf[b.offset+2] = byte(data >> 16)
	b.buf[b.offset+3] = byte(data >> 24)
	b.offset += 4
}

func (b *Buffer) WriteInt32BE(data int32) {
	b.EnsureCap(4)
	b.buf[b.offset] = byte(data >> 24)
	b.buf[b.offset+1] = byte(data >> 16)
	b.buf[b.offset+2] = byte(data >> 8)
	b.buf[b.offset+3] = byte(data)
	b.offset += 4
}

func (b *Buffer) WriteInt64LE(data int64) {
	b.EnsureCap(8)
	b.buf[b.offset] = byte(data)
	b.buf[b.offset+1] = byte(data >> 8)
	b.buf[b.offset+2] = byte(data >> 16)
	b.buf[b.offset+3] = byte(data >> 24)
	b.buf[b.offset+4] = byte(data >> 32)
	b.buf[b.offset+5] = byte(data >> 40)
	b.buf[b.offset+6] = byte(data >> 48)
	b.buf[b.offset+7] = byte(data >> 56)
	b.offset += 8
}

func (b *Buffer) WriteInt64BE(data int64) {
	b.EnsureCap(8)
	b.buf[b.offset] = byte(data >> 56)
	b.buf[b.offset+1] = byte(data >> 48)
	b.buf[b.offset+2] = byte(data >> 40)
	b.buf[b.offset+3] = byte(data >> 32)
	b.buf[b.offset+4] = byte(data >> 24)
	b.buf[b.offset+5] = byte(data >> 16)
	b.buf[b.offset+6] = byte(data >> 8)
	b.buf[b.offset+7] = byte(data)
	b.offset += 8
}

func (b *Buffer) WriteUint8(data uint8) {
	b.WriteByte(byte(data))
	b.offset += 1
}

func (b *Buffer) WriteUint16LE(data uint16) {
	b.EnsureCap(2)
	b.buf[b.offset] = byte(data)
	b.buf[b.offset+1] = byte(data >> 8)
	b.offset += 2
}

func (b *Buffer) WriteUint16BE(data uint16) {
	b.EnsureCap(2)
	b.buf[b.offset] = byte(data >> 8)
	b.buf[b.offset+1] = byte(data)
	b.offset += 2
}

func (b *Buffer) WriteUint32LE(data uint32) {
	b.EnsureCap(4)
	b.buf[b.offset] = byte(data)
	b.buf[b.offset+1] = byte(data >> 8)
	b.buf[b.offset+2] = byte(data >> 16)
	b.buf[b.offset+3] = byte(data >> 24)
	b.offset += 4
}

func (b *Buffer) WriteUint32BE(data uint32) {
	b.EnsureCap(4)
	b.buf[b.offset] = byte(data >> 24)
	b.buf[b.offset+1] = byte(data >> 16)
	b.buf[b.offset+2] = byte(data >> 8)
	b.buf[b.offset+3] = byte(data)
	b.offset += 4
}

func (b *Buffer) WriteUint64LE(data uint64) {
	b.EnsureCap(8)
	b.buf[b.offset] = byte(data)
	b.buf[b.offset+1] = byte(data >> 8)
	b.buf[b.offset+2] = byte(data >> 16)
	b.buf[b.offset+3] = byte(data >> 24)
	b.buf[b.offset+4] = byte(data >> 32)
	b.buf[b.offset+5] = byte(data >> 40)
	b.buf[b.offset+6] = byte(data >> 48)
	b.buf[b.offset+7] = byte(data >> 56)
	b.offset += 8
}

func (b *Buffer) WriteUint64BE(data uint64) {
	b.EnsureCap(8)
	b.buf[b.offset] = byte(data >> 56)
	b.buf[b.offset+1] = byte(data >> 48)
	b.buf[b.offset+2] = byte(data >> 40)
	b.buf[b.offset+3] = byte(data >> 32)
	b.buf[b.offset+4] = byte(data >> 24)
	b.buf[b.offset+5] = byte(data >> 16)
	b.buf[b.offset+6] = byte(data >> 8)
	b.buf[b.offset+7] = byte(data)
	b.offset += 8
}

func (b *Buffer) WriteFloat32LE(data float32) {
	b.EnsureCap(4)
	var val uint64 = *(*uint64)(unsafe.Pointer(&data))
	b.buf[b.offset] = byte(val)
	b.buf[b.offset+1] = byte(val >> 8)
	b.buf[b.offset+2] = byte(val >> 16)
	b.buf[b.offset+3] = byte(val >> 24)
	b.offset += 4
}

func (b *Buffer) WriteFloat32BE(data float32) {
	b.EnsureCap(4)
	var val uint64 = *(*uint64)(unsafe.Pointer(&data))
	b.buf[b.offset] = byte(val >> 24)
	b.buf[b.offset+1] = byte(val >> 16)
	b.buf[b.offset+2] = byte(val >> 8)
	b.buf[b.offset+3] = byte(val)
	b.offset += 4
}

func (b *Buffer) WriteFloat64LE(data float64) {
	b.EnsureCap(8)
	var val uint64 = *(*uint64)(unsafe.Pointer(&data))
	b.buf[b.offset] = byte(val)
	b.buf[b.offset+1] = byte(val >> 8)
	b.buf[b.offset+2] = byte(val >> 16)
	b.buf[b.offset+3] = byte(val >> 24)
	b.buf[b.offset+4] = byte(val >> 32)
	b.buf[b.offset+5] = byte(val >> 40)
	b.buf[b.offset+6] = byte(val >> 48)
	b.buf[b.offset+7] = byte(val >> 56)
	b.offset += 8
}

func (b *Buffer) WriteFloat64BE(data float64) {
	b.EnsureCap(8)
	var val uint64 = *(*uint64)(unsafe.Pointer(&data))
	b.buf[b.offset] = byte(val >> 56)
	b.buf[b.offset+1] = byte(val >> 48)
	b.buf[b.offset+2] = byte(val >> 40)
	b.buf[b.offset+3] = byte(val >> 32)
	b.buf[b.offset+4] = byte(val >> 24)
	b.buf[b.offset+5] = byte(val >> 16)
	b.buf[b.offset+6] = byte(val >> 8)
	b.buf[b.offset+7] = byte(val)
	b.offset += 8
}

func (b *Buffer) ReadInt8() int8 {
	if b.offset >= b.cap {
		panic(BufferOverreadError)
	}
	if b.offset < 0 {
		panic(BufferUnderreadError)
	}
	b.offset += 1
	return int8(b.buf[b.offset-1])
}

func (b *Buffer) ReadInt16LE() int16 {
	if (b.offset + 2) > b.cap {
		panic(BufferOverreadError)
	}
	if b.offset < 0 {
		panic(BufferUnderreadError)
	}
	var result int16 = int16(b.buf[b.offset]) | int16(b.buf[b.offset+1])<<8
	b.offset += 2
	return result
}

func (b *Buffer) ReadInt16BE() int16 {
	if (b.offset + 2) > b.cap {
		panic(BufferOverreadError)
	}
	if b.offset < 0 {
		panic(BufferUnderreadError)
	}
	var result int16 = int16(b.buf[b.offset])<<8 | int16(b.buf[b.offset+1])
	b.offset += 2
	return result
}

func (b *Buffer) ReadInt32LE() int32 {
	if (b.offset + 4) > b.cap {
		panic(BufferOverreadError)
	}
	if b.offset < 0 {
		panic(BufferUnderreadError)
	}
	var result int32 = int32(b.buf[b.offset]) | int32(b.buf[b.offset+1])<<8 | int32(b.buf[b.offset+2])<<16 | int32(b.buf[b.offset+3])<<24
	b.offset += 4
	return result
}

func (b *Buffer) ReadInt32BE() int32 {
	if (b.offset + 4) > b.cap {
		panic(BufferOverreadError)
	}
	if b.offset < 0 {
		panic(BufferUnderreadError)
	}
	var result int32 = int32(b.buf[b.offset])<<24 | int32(b.buf[b.offset+1])<<16 | int32(b.buf[b.offset+2])<<8 | int32(b.buf[b.offset+3])
	b.offset += 4
	return result
}

func (b *Buffer) ReadInt64LE() int64 {
	if (b.offset + 8) > b.cap {
		panic(BufferOverreadError)
	}
	if b.offset < 0 {
		panic(BufferUnderreadError)
	}
	var result int64 = int64(b.buf[b.offset]) | int64(b.buf[b.offset+1])<<8 | int64(b.buf[b.offset+2])<<16 | int64(b.buf[b.offset+3])<<24 | int64(b.buf[b.offset+4])<<32 | int64(b.buf[b.offset+5])<<40 | int64(b.buf[b.offset+6])<<48 | int64(b.buf[b.offset+7])<<56
	b.offset += 8
	return result
}

func (b *Buffer) ReadInt64BE() int64 {
	if (b.offset + 8) > b.cap {
		panic(BufferOverreadError)
	}
	if b.offset < 0 {
		panic(BufferUnderreadError)
	}
	var result int64 = int64(b.buf[b.offset])<<56 | int64(b.buf[b.offset+1])<<48 | int64(b.buf[b.offset+2])<<40 | int64(b.buf[b.offset+3])<<32 | int64(b.buf[b.offset+4])<<24 | int64(b.buf[b.offset+5])<<16 | int64(b.buf[b.offset+6])<<8 | int64(b.buf[b.offset+7])
	b.offset += 8
	return result
}

func (b *Buffer) ReadUint8() uint8 {
	if b.offset >= b.cap {
		panic(BufferOverreadError)
	}
	if b.offset < 0 {
		panic(BufferUnderreadError)
	}
	b.offset += 1
	return uint8(b.buf[b.offset-1])
}

func (b *Buffer) ReadUint16LE() uint16 {
	if (b.offset + 2) > b.cap {
		panic(BufferOverreadError)
	}
	if b.offset < 0 {
		panic(BufferUnderreadError)
	}
	var result uint16 = uint16(b.buf[b.offset]) | uint16(b.buf[b.offset+1])<<8
	b.offset += 2
	return result
}

func (b *Buffer) ReadUint16BE() uint16 {
	if (b.offset + 2) > b.cap {
		panic(BufferOverreadError)
	}
	if b.offset < 0 {
		panic(BufferUnderreadError)
	}
	var result uint16 = uint16(b.buf[b.offset])<<8 | uint16(b.buf[b.offset+1])
	b.offset += 2
	return result
}

func (b *Buffer) ReadUint32LE() uint32 {
	if (b.offset + 4) > b.cap {
		panic(BufferOverreadError)
	}
	if b.offset < 0 {
		panic(BufferUnderreadError)
	}
	var result uint32 = uint32(b.buf[b.offset]) | uint32(b.buf[b.offset+1])<<8 | uint32(b.buf[b.offset+2])<<16 | uint32(b.buf[b.offset+3])<<24
	b.offset += 4
	return result
}

func (b *Buffer) ReadUint32BE() uint32 {
	if (b.offset + 4) > b.cap {
		panic(BufferOverreadError)
	}
	if b.offset < 0 {
		panic(BufferUnderreadError)
	}
	var result uint32 = uint32(b.buf[b.offset])<<24 | uint32(b.buf[b.offset+1])<<16 | uint32(b.buf[b.offset+2])<<8 | uint32(b.buf[b.offset+3])
	b.offset += 4
	return result
}

func (b *Buffer) ReadUint64LE() uint64 {
	if (b.offset + 8) > b.cap {
		panic(BufferOverreadError)
	}
	if b.offset < 0 {
		panic(BufferUnderreadError)
	}
	var result uint64 = uint64(b.buf[b.offset]) | uint64(b.buf[b.offset+1])<<8 | uint64(b.buf[b.offset+2])<<16 | uint64(b.buf[b.offset+3])<<24 | uint64(b.buf[b.offset+4])<<32 | uint64(b.buf[b.offset+5])<<40 | uint64(b.buf[b.offset+6])<<48 | uint64(b.buf[b.offset+7])<<56
	b.offset += 8
	return result
}

func (b *Buffer) ReadUint64BE() uint64 {
	if (b.offset + 8) > b.cap {
		panic(BufferOverreadError)
	}
	if b.offset < 0 {
		panic(BufferUnderreadError)
	}
	var result uint64 = uint64(b.buf[b.offset])<<56 | uint64(b.buf[b.offset+1])<<48 | uint64(b.buf[b.offset+2])<<40 | uint64(b.buf[b.offset+3])<<32 | uint64(b.buf[b.offset+4])<<24 | uint64(b.buf[b.offset+5])<<16 | uint64(b.buf[b.offset+6])<<8 | uint64(b.buf[b.offset+7])
	b.offset += 8
	return result
}

func (b *Buffer) ReadFloat32LE() float32 {
	if (b.offset + 4) > b.cap {
		panic(BufferOverreadError)
	}
	if b.offset < 0 {
		panic(BufferUnderreadError)
	}
	var val uint64 = uint64(b.buf[b.offset]) | uint64(b.buf[b.offset+1])<<8 | uint64(b.buf[b.offset+2])<<16 | uint64(b.buf[b.offset+3])<<24
	var result float32 = *(*float32)(unsafe.Pointer(&val))
	b.offset += 4
	return result
}

func (b *Buffer) ReadFloat32BE() float32 {
	if (b.offset + 4) > b.cap {
		panic(BufferOverreadError)
	}
	if b.offset < 0 {
		panic(BufferUnderreadError)
	}
	var val uint64 = uint64(b.buf[b.offset])<<24 | uint64(b.buf[b.offset+1])<<16 | uint64(b.buf[b.offset+2])<<8 | uint64(b.buf[b.offset+3])
	var result float32 = *(*float32)(unsafe.Pointer(&val))
	b.offset += 4
	return result
}

func (b *Buffer) ReadFloat64LE() float64 {
	if (b.offset + 8) > b.cap {
		panic(BufferOverreadError)
	}
	if b.offset < 0 {
		panic(BufferUnderreadError)
	}
	var val uint64 = uint64(b.buf[b.offset]) | uint64(b.buf[b.offset+1])<<8 | uint64(b.buf[b.offset+2])<<16 | uint64(b.buf[b.offset+3])<<24 | uint64(b.buf[b.offset+4])<<32 | uint64(b.buf[b.offset+5])<<40 | uint64(b.buf[b.offset+6])<<48 | uint64(b.buf[b.offset+7])<<56
	var result float64 = *(*float64)(unsafe.Pointer(&val))
	b.offset += 8
	return result
}

func (b *Buffer) ReadFloat64BE() float64 {
	if (b.offset + 8) > b.cap {
		panic(BufferOverreadError)
	}
	if b.offset < 0 {
		panic(BufferUnderreadError)
	}
	var val uint64 = uint64(b.buf[b.offset])<<56 | uint64(b.buf[b.offset+1])<<48 | uint64(b.buf[b.offset+2])<<40 | uint64(b.buf[b.offset+3])<<32 | uint64(b.buf[b.offset+4])<<24 | uint64(b.buf[b.offset+5])<<16 | uint64(b.buf[b.offset+6])<<8 | uint64(b.buf[b.offset+7])
	var result float64 = *(*float64)(unsafe.Pointer(&val))
	b.offset += 8
	return result
}
