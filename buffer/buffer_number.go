package buffer

import "unsafe"

// WriteInt8 writes int8 to the buffer at current offset
// and moves the offset forward 1 byte.
//
// This method will ensure buffer capacity is enough to write data.
//
// This method returns current Buffer instance for chainable operation.
func (b *Buffer) WriteInt8(data int8) *Buffer {
	b.EnsureCap(1)
	b.WriteByte(byte(data))
	b.offset += 1
	return b
}

// WriteInt8Unsafe writes int8 to the buffer at current offset
// and moves the offset forward 1 byte
//
// This method doesn't check buffer capacity.
//
// This method returns current Buffer instance for chainable operation.
func (b *Buffer) WriteInt8Unsafe(data int8) *Buffer {
	b.WriteByte(byte(data))
	b.offset += 1
	return b
}

// WriteInt16LE writes int16 to the buffer at current offset in little-endian
// and moves the offset forward 2 bytes.
//
// This method will ensure buffer capacity is enough to write data.
//
// This method returns current Buffer instance for chainable operation.
func (b *Buffer) WriteInt16LE(data int16) *Buffer {
	b.EnsureCap(2)
	b.buf[b.offset] = byte(data)
	b.buf[b.offset+1] = byte(data >> 8)
	b.offset += 2
	return b
}

// WriteInt16LEUnsafe writes int16 to the buffer at current offset in little-endian
// and moves the offset forward 2 bytes.
//
// This method doesn't check buffer capacity.
//
// This method returns current Buffer instance for chainable operation.
func (b *Buffer) WriteInt16LEUnsafe(data int16) *Buffer {
	b.buf[b.offset] = byte(data)
	b.buf[b.offset+1] = byte(data >> 8)
	b.offset += 2
	return b
}

// WriteInt16BE writes int16 to the buffer at current offset in big-endian
// and moves the offset forward 2 bytes.
//
// This method will ensure buffer capacity is enough to write data.
//
// This method returns current Buffer instance for chainable operation.
func (b *Buffer) WriteInt16BE(data int16) *Buffer {
	b.EnsureCap(2)
	b.buf[b.offset] = byte(data >> 8)
	b.buf[b.offset+1] = byte(data)
	b.offset += 2
	return b
}

// WriteInt16BEUnsafe writes int16 to the buffer at current offset in big-endian
// and moves the offset forward 2 bytes.
//
// This method doesn't check buffer capacity.
//
// This method returns current Buffer instance for chainable operation.
func (b *Buffer) WriteInt16BEUnsafe(data int16) *Buffer {
	b.buf[b.offset] = byte(data >> 8)
	b.buf[b.offset+1] = byte(data)
	b.offset += 2
	return b
}

// WriteInt32LE writes int32 to the buffer at current offset in little-endian
// and moves the offset forward 4 bytes.
//
// This method will ensure buffer capacity is enough to write data.
//
// This method returns current Buffer instance for chainable operation.
func (b *Buffer) WriteInt32LE(data int32) *Buffer {
	b.EnsureCap(4)
	b.buf[b.offset] = byte(data)
	b.buf[b.offset+1] = byte(data >> 8)
	b.buf[b.offset+2] = byte(data >> 16)
	b.buf[b.offset+3] = byte(data >> 24)
	b.offset += 4
	return b
}

// WriteInt32LEUnsafe writes int32 to the buffer at current offset in little-endian
// and moves the offset forward 4 bytes.
//
// This method doesn't check buffer capacity.
//
// This method returns current Buffer instance for chainable operation.
func (b *Buffer) WriteInt32LEUnsafe(data int32) *Buffer {
	b.buf[b.offset] = byte(data)
	b.buf[b.offset+1] = byte(data >> 8)
	b.buf[b.offset+2] = byte(data >> 16)
	b.buf[b.offset+3] = byte(data >> 24)
	b.offset += 4
	return b
}

// WriteInt32BE writes int32 to the buffer at current offset in big-endian
// and moves the offset forward 4 bytes.
//
// This method will ensure buffer capacity is enough to write data.
//
// This method returns current Buffer instance for chainable operation.
func (b *Buffer) WriteInt32BE(data int32) *Buffer {
	b.EnsureCap(4)
	b.buf[b.offset] = byte(data >> 24)
	b.buf[b.offset+1] = byte(data >> 16)
	b.buf[b.offset+2] = byte(data >> 8)
	b.buf[b.offset+3] = byte(data)
	b.offset += 4
	return b
}

// WriteInt32BEUnsafe writes int32 to the buffer at current offset in big-endian
// and moves the offset forward 4 bytes.
//
// This method doesn't check buffer capacity.
//
// This method returns current Buffer instance for chainable operation.
func (b *Buffer) WriteInt32BEUnsafe(data int32) *Buffer {
	b.buf[b.offset] = byte(data >> 24)
	b.buf[b.offset+1] = byte(data >> 16)
	b.buf[b.offset+2] = byte(data >> 8)
	b.buf[b.offset+3] = byte(data)
	b.offset += 4
	return b
}

// WriteInt64LE writes int64 to the buffer at current offset in little-endian
// and moves the offset forward 8 bytes.
//
// This method will ensure buffer capacity is enough to write data.
//
// This method returns current Buffer instance for chainable operation.
func (b *Buffer) WriteInt64LE(data int64) *Buffer {
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
	return b
}

// WriteInt64LEUnsafe writes int64 to the buffer at current offset in little-endian
// and moves the offset forward 8 bytes.
//
// This method doesn't check buffer capacity.
//
// This method returns current Buffer instance for chainable operation.
func (b *Buffer) WriteInt64LEUnsafe(data int64) *Buffer {
	b.buf[b.offset] = byte(data)
	b.buf[b.offset+1] = byte(data >> 8)
	b.buf[b.offset+2] = byte(data >> 16)
	b.buf[b.offset+3] = byte(data >> 24)
	b.buf[b.offset+4] = byte(data >> 32)
	b.buf[b.offset+5] = byte(data >> 40)
	b.buf[b.offset+6] = byte(data >> 48)
	b.buf[b.offset+7] = byte(data >> 56)
	b.offset += 8
	return b
}

// WriteInt64BE writes int64 to the buffer at current offset in big-endian
// and moves the offset forward 8 bytes.
//
// This method will ensure buffer capacity is enough to write data.
//
// This method returns current Buffer instance for chainable operation.
func (b *Buffer) WriteInt64BE(data int64) *Buffer {
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
	return b
}

// WriteInt64BEUnsafe writes int64 to the buffer at current offset in big-endian
// and moves the offset forward 8 bytes.
//
// This method doesn't check buffer capacity.
//
// This method returns current Buffer instance for chainable operation.
func (b *Buffer) WriteInt64BEUnsafe(data int64) *Buffer {
	b.buf[b.offset] = byte(data >> 56)
	b.buf[b.offset+1] = byte(data >> 48)
	b.buf[b.offset+2] = byte(data >> 40)
	b.buf[b.offset+3] = byte(data >> 32)
	b.buf[b.offset+4] = byte(data >> 24)
	b.buf[b.offset+5] = byte(data >> 16)
	b.buf[b.offset+6] = byte(data >> 8)
	b.buf[b.offset+7] = byte(data)
	b.offset += 8
	return b
}

// WriteUint8 writes uint8 to the buffer at current offset
// and moves the offset forward 1 byte.
//
// This method will ensure buffer capacity is enough to write data.
//
// This method returns current Buffer instance for chainable operation.
func (b *Buffer) WriteUint8(data uint8) *Buffer {
	b.EnsureCap(1)
	b.WriteByte(byte(data))
	b.offset += 1
	return b
}

// WriteUint8Unsafe writes uint8 to the buffer at current offset
// and moves the offset forward 1 byte
//
// This method doesn't check buffer capacity.
//
// This method returns current Buffer instance for chainable operation.
func (b *Buffer) WriteUint8Unsafe(data uint8) *Buffer {
	b.WriteByte(byte(data))
	b.offset += 1
	return b
}

// WriteUint16LE writes uint16 to the buffer at current offset in little-endian
// and moves the offset forward 2 bytes.
//
// This method will ensure buffer capacity is enough to write data.
//
// This method returns current Buffer instance for chainable operation.
func (b *Buffer) WriteUint16LE(data uint16) *Buffer {
	b.EnsureCap(2)
	b.buf[b.offset] = byte(data)
	b.buf[b.offset+1] = byte(data >> 8)
	b.offset += 2
	return b
}

// WriteUint16LEUnsafe writes uint16 to the buffer at current offset in little-endian
// and moves the offset forward 2 bytes.
//
// This method doesn't check buffer capacity.
//
// This method returns current Buffer instance for chainable operation.
func (b *Buffer) WriteUint16LEUnsafe(data uint16) *Buffer {
	b.buf[b.offset] = byte(data)
	b.buf[b.offset+1] = byte(data >> 8)
	b.offset += 2
	return b
}

// WriteUint16BE writes uint16 to the buffer at current offset in big-endian
// and moves the offset forward 2 bytes.
//
// This method will ensure buffer capacity is enough to write data.
//
// This method returns current Buffer instance for chainable operation.
func (b *Buffer) WriteUint16BE(data uint16) *Buffer {
	b.EnsureCap(2)
	b.buf[b.offset] = byte(data >> 8)
	b.buf[b.offset+1] = byte(data)
	b.offset += 2
	return b
}

// WriteUint16BEUnsafe writes uint16 to the buffer at current offset in big-endian
// and moves the offset forward 2 bytes.
//
// This method doesn't check buffer capacity.
//
// This method returns current Buffer instance for chainable operation.
func (b *Buffer) WriteUint16BEUnsafe(data uint16) *Buffer {
	b.buf[b.offset] = byte(data >> 8)
	b.buf[b.offset+1] = byte(data)
	b.offset += 2
	return b
}

// WriteUint32LE writes uint32 to the buffer at current offset in little-endian
// and moves the offset forward 4 bytes.
//
// This method will ensure buffer capacity is enough to write data.
//
// This method returns current Buffer instance for chainable operation.
func (b *Buffer) WriteUint32LE(data uint32) *Buffer {
	b.EnsureCap(4)
	b.buf[b.offset] = byte(data)
	b.buf[b.offset+1] = byte(data >> 8)
	b.buf[b.offset+2] = byte(data >> 16)
	b.buf[b.offset+3] = byte(data >> 24)
	b.offset += 4
	return b
}

// WriteUint32LEUnsafe writes uint32 to the buffer at current offset in little-endian
// and moves the offset forward 4 bytes.
//
// This method doesn't check buffer capacity.
//
// This method returns current Buffer instance for chainable operation.
func (b *Buffer) WriteUint32LEUnsafe(data uint32) *Buffer {
	b.buf[b.offset] = byte(data)
	b.buf[b.offset+1] = byte(data >> 8)
	b.buf[b.offset+2] = byte(data >> 16)
	b.buf[b.offset+3] = byte(data >> 24)
	b.offset += 4
	return b
}

// WriteUint32BE writes uint32 to the buffer at current offset in big-endian
// and moves the offset forward 4 bytes.
//
// This method will ensure buffer capacity is enough to write data.
//
// This method returns current Buffer instance for chainable operation.
func (b *Buffer) WriteUint32BE(data uint32) *Buffer {
	b.EnsureCap(4)
	b.buf[b.offset] = byte(data >> 24)
	b.buf[b.offset+1] = byte(data >> 16)
	b.buf[b.offset+2] = byte(data >> 8)
	b.buf[b.offset+3] = byte(data)
	b.offset += 4
	return b
}

// WriteUint32BEUnsafe writes uint32 to the buffer at current offset in big-endian
// and moves the offset forward 4 bytes.
//
// This method doesn't check buffer capacity.
//
// This method returns current Buffer instance for chainable operation.
func (b *Buffer) WriteUint32BEUnsafe(data uint32) *Buffer {
	b.buf[b.offset] = byte(data >> 24)
	b.buf[b.offset+1] = byte(data >> 16)
	b.buf[b.offset+2] = byte(data >> 8)
	b.buf[b.offset+3] = byte(data)
	b.offset += 4
	return b
}

// WriteUint64LE writes uint64 to the buffer at current offset in little-endian
// and moves the offset forward 8 bytes.
//
// This method will ensure buffer capacity is enough to write data.
//
// This method returns current Buffer instance for chainable operation.
func (b *Buffer) WriteUint64LE(data uint64) *Buffer {
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
	return b
}

// WriteUint64LEUnsafe writes uint64 to the buffer at current offset in little-endian
// and moves the offset forward 8 bytes.
//
// This method doesn't check buffer capacity.
//
// This method returns current Buffer instance for chainable operation.
func (b *Buffer) WriteUint64LEUnsafe(data uint64) *Buffer {
	b.buf[b.offset] = byte(data)
	b.buf[b.offset+1] = byte(data >> 8)
	b.buf[b.offset+2] = byte(data >> 16)
	b.buf[b.offset+3] = byte(data >> 24)
	b.buf[b.offset+4] = byte(data >> 32)
	b.buf[b.offset+5] = byte(data >> 40)
	b.buf[b.offset+6] = byte(data >> 48)
	b.buf[b.offset+7] = byte(data >> 56)
	b.offset += 8
	return b
}

// WriteUint64BE writes uint64 to the buffer at current offset in big-endian
// and moves the offset forward 8 bytes.
//
// This method will ensure buffer capacity is enough to write data.
//
// This method returns current Buffer instance for chainable operation.
func (b *Buffer) WriteUint64BE(data uint64) *Buffer {
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
	return b
}

// WriteUint64BEUnsafe writes uint64 to the buffer at current offset in big-endian
// and moves the offset forward 8 bytes.
//
// This method doesn't check buffer capacity.
//
// This method returns current Buffer instance for chainable operation.
func (b *Buffer) WriteUint64BEUnsafe(data uint64) *Buffer {
	b.buf[b.offset] = byte(data >> 56)
	b.buf[b.offset+1] = byte(data >> 48)
	b.buf[b.offset+2] = byte(data >> 40)
	b.buf[b.offset+3] = byte(data >> 32)
	b.buf[b.offset+4] = byte(data >> 24)
	b.buf[b.offset+5] = byte(data >> 16)
	b.buf[b.offset+6] = byte(data >> 8)
	b.buf[b.offset+7] = byte(data)
	b.offset += 8
	return b
}

// WriteFloat32LE writes float32 to the buffer at current offset in little-endian
// and moves the offset forward 4 bytes.
//
// This method will ensure buffer capacity is enough to write data.
//
// This method returns current Buffer instance for chainable operation.
func (b *Buffer) WriteFloat32LE(data float32) *Buffer {
	b.EnsureCap(4)
	var val uint64 = *(*uint64)(unsafe.Pointer(&data))
	b.buf[b.offset] = byte(val)
	b.buf[b.offset+1] = byte(val >> 8)
	b.buf[b.offset+2] = byte(val >> 16)
	b.buf[b.offset+3] = byte(val >> 24)
	b.offset += 4
	return b
}

// WriteFloat32LEUnsafe writes float32 to the buffer at current offset in little-endian
// and moves the offset forward 4 bytes.
//
// This method doesn't check buffer capacity.
//
// This method returns current Buffer instance for chainable operation.
func (b *Buffer) WriteFloat32LEUnsafe(data float32) *Buffer {
	var val uint64 = *(*uint64)(unsafe.Pointer(&data))
	b.buf[b.offset] = byte(val)
	b.buf[b.offset+1] = byte(val >> 8)
	b.buf[b.offset+2] = byte(val >> 16)
	b.buf[b.offset+3] = byte(val >> 24)
	b.offset += 4
	return b
}

// WriteFloat32BE writes float32 to the buffer at current offset in big-endian
// and moves the offset forward 4 bytes.
//
// This method will ensure buffer capacity is enough to write data.
//
// This method returns current Buffer instance for chainable operation.
func (b *Buffer) WriteFloat32BE(data float32) *Buffer {
	b.EnsureCap(4)
	var val uint64 = *(*uint64)(unsafe.Pointer(&data))
	b.buf[b.offset] = byte(val >> 24)
	b.buf[b.offset+1] = byte(val >> 16)
	b.buf[b.offset+2] = byte(val >> 8)
	b.buf[b.offset+3] = byte(val)
	b.offset += 4
	return b
}

// WriteFloat32BEUnsafe writes float32 to the buffer at current offset in big-endian
// and moves the offset forward 4 bytes.
//
// This method doesn't check buffer capacity.
//
// This method returns current Buffer instance for chainable operation.
func (b *Buffer) WriteFloat32BEUnsafe(data float32) *Buffer {
	var val uint64 = *(*uint64)(unsafe.Pointer(&data))
	b.buf[b.offset] = byte(val >> 24)
	b.buf[b.offset+1] = byte(val >> 16)
	b.buf[b.offset+2] = byte(val >> 8)
	b.buf[b.offset+3] = byte(val)
	b.offset += 4
	return b
}

// WriteFloat64LE writes float64 to the buffer at current offset in little-endian
// and moves the offset forward 8 bytes.
//
// This method will ensure buffer capacity is enough to write data.
//
// This method returns current Buffer instance for chainable operation.
func (b *Buffer) WriteFloat64LE(data float64) *Buffer {
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
	return b
}

// WriteFloat64LEUnsafe writes float64 to the buffer at current offset in little-endian
// and moves the offset forward 8 bytes.
//
// This method doesn't check buffer capacity.
//
// This method returns current Buffer instance for chainable operation.
func (b *Buffer) WriteFloat64LEUnsafe(data float64) *Buffer {
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
	return b
}

// WriteFloat64BE writes float64 to the buffer at current offset in big-endian
// and moves the offset forward 8 bytes.
//
// This method will ensure buffer capacity is enough to write data.
//
// This method returns current Buffer instance for chainable operation.
func (b *Buffer) WriteFloat64BE(data float64) *Buffer {
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
	return b
}

// WriteFloat64BEUnsafe writes float64 to the buffer at current offset in big-endian
// and moves the offset forward 8 bytes.
//
// This method doesn't check buffer capacity.
//
// This method returns current Buffer instance for chainable operation.
func (b *Buffer) WriteFloat64BEUnsafe(data float64) *Buffer {
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
	return b
}

// ReadInt8 reads int8 from the buffer at current offset
// and moves the offset forward 1 byte.
//
// If the offset is invalid or not safe to read data,
// it will panic with BufferOverreadError or BufferUnderreadError.
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

// ReadInt8Unsafe reads int8 from the buffer at current offset
// and moves the offset forward 1 byte
//
// This method doesn't check the read safety.
func (b *Buffer) ReadInt8Unsafe() int8 {
	b.offset += 1
	return int8(b.buf[b.offset-1])
}

// ReadInt16LE reads int16 from the buffer at current offset in little-endian
// and moves the offset forward 2 bytes.
//
// If the offset is invalid or not safe to read data,
// it will panic with BufferOverreadError or BufferUnderreadError.
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

// ReadInt16LEUnsafe reads int16 from the buffer at current offset in little-endian
// and moves the offset forward 2 bytes.
//
// This method doesn't check the read safety.
func (b *Buffer) ReadInt16LEUnsafe() int16 {
	var result int16 = int16(b.buf[b.offset]) | int16(b.buf[b.offset+1])<<8
	b.offset += 2
	return result
}

// ReadInt16BE reads int16 from the buffer at current offset in big-endian
// and moves the offset forward 2 bytes.
//
// If the offset is invalid or not safe to read data,
// it will panic with BufferOverreadError or BufferUnderreadError.
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

// ReadInt16BEUnsafe reads int16 from the buffer at current offset in big-endian
// and moves the offset forward 2 bytes.
//
// This method doesn't check the read safety.
func (b *Buffer) ReadInt16BEUnsafe() int16 {
	var result int16 = int16(b.buf[b.offset])<<8 | int16(b.buf[b.offset+1])
	b.offset += 2
	return result
}

// ReadInt32LE reads int32 from the buffer at current offset in little-endian
// and moves the offset forward 4 bytes.
//
// If the offset is invalid or not safe to read data,
// it will panic with BufferOverreadError or BufferUnderreadError.
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

// ReadInt32LEUnsafe reads int32 from the buffer at current offset in little-endian
// and moves the offset forward 4 bytes.
//
// This method doesn't check the read safety.
func (b *Buffer) ReadInt32LEUnsafe() int32 {
	var result int32 = int32(b.buf[b.offset]) | int32(b.buf[b.offset+1])<<8 | int32(b.buf[b.offset+2])<<16 | int32(b.buf[b.offset+3])<<24
	b.offset += 4
	return result
}

// ReadInt32BE reads int32 from the buffer at current offset in big-endian
// and moves the offset forward 4 bytes.
//
// If the offset is invalid or not safe to read data,
// it will panic with BufferOverreadError or BufferUnderreadError.
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

// ReadInt32BEUnsafe reads int32 from the buffer at current offset in big-endian
// and moves the offset forward 4 bytes.
//
// This method doesn't check the read safety.
func (b *Buffer) ReadInt32BEUnsafe() int32 {
	var result int32 = int32(b.buf[b.offset])<<24 | int32(b.buf[b.offset+1])<<16 | int32(b.buf[b.offset+2])<<8 | int32(b.buf[b.offset+3])
	b.offset += 4
	return result
}

// ReadInt64LE reads int64 from the buffer at current offset in little-endian
// and moves the offset forward 8 bytes.
//
// If the offset is invalid or not safe to read data,
// it will panic with BufferOverreadError or BufferUnderreadError.
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

// ReadInt64LEUnsafe reads int64 from the buffer at current offset in little-endian
// and moves the offset forward 8 bytes.
//
// This method doesn't check the read safety.
func (b *Buffer) ReadInt64LEUnsafe() int64 {
	var result int64 = int64(b.buf[b.offset]) | int64(b.buf[b.offset+1])<<8 | int64(b.buf[b.offset+2])<<16 | int64(b.buf[b.offset+3])<<24 | int64(b.buf[b.offset+4])<<32 | int64(b.buf[b.offset+5])<<40 | int64(b.buf[b.offset+6])<<48 | int64(b.buf[b.offset+7])<<56
	b.offset += 8
	return result
}

// ReadInt64BE reads int64 from the buffer at current offset in big-endian
// and moves the offset forward 8 bytes.
//
// If the offset is invalid or not safe to read data,
// it will panic with BufferOverreadError or BufferUnderreadError.
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

// ReadInt64BEUnsafe reads int64 from the buffer at current offset in big-endian
// and moves the offset forward 8 bytes.
//
// This method doesn't check the read safety.
func (b *Buffer) ReadInt64BEUnsafe() int64 {
	var result int64 = int64(b.buf[b.offset])<<56 | int64(b.buf[b.offset+1])<<48 | int64(b.buf[b.offset+2])<<40 | int64(b.buf[b.offset+3])<<32 | int64(b.buf[b.offset+4])<<24 | int64(b.buf[b.offset+5])<<16 | int64(b.buf[b.offset+6])<<8 | int64(b.buf[b.offset+7])
	b.offset += 8
	return result
}

// ReadUint8 reads uint8 from the buffer at current offset
// and moves the offset forward 1 byte.
//
// If the offset is invalid or not safe to read data,
// it will panic with BufferOverreadError or BufferUnderreadError.
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

// ReadUint8Unsafe reads uint8 from the buffer at current offset
// and moves the offset forward 1 byte
//
// This method doesn't check the read safety.
func (b *Buffer) ReadUint8Unsafe() uint8 {
	b.offset += 1
	return uint8(b.buf[b.offset-1])
}

// ReadUint16LE reads uint16 from the buffer at current offset in little-endian
// and moves the offset forward 2 bytes.
//
// If the offset is invalid or not safe to read data,
// it will panic with BufferOverreadError or BufferUnderreadError.
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

// ReadUint16LEUnsafe reads uint16 from the buffer at current offset in little-endian
// and moves the offset forward 2 bytes.
//
// This method doesn't check the read safety.
func (b *Buffer) ReadUint16LEUnsafe() uint16 {
	var result uint16 = uint16(b.buf[b.offset]) | uint16(b.buf[b.offset+1])<<8
	b.offset += 2
	return result
}

// ReadUint16BE reads uint16 from the buffer at current offset in big-endian
// and moves the offset forward 2 bytes.
//
// If the offset is invalid or not safe to read data,
// it will panic with BufferOverreadError or BufferUnderreadError.
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

// ReadUint16BEUnsafe reads uint16 from the buffer at current offset in big-endian
// and moves the offset forward 2 bytes.
//
// This method doesn't check the read safety.
func (b *Buffer) ReadUint16BEUnsafe() uint16 {
	var result uint16 = uint16(b.buf[b.offset])<<8 | uint16(b.buf[b.offset+1])
	b.offset += 2
	return result
}

// ReadUint32LE reads uint32 from the buffer at current offset in little-endian
// and moves the offset forward 4 bytes.
//
// If the offset is invalid or not safe to read data,
// it will panic with BufferOverreadError or BufferUnderreadError.
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

// ReadUint32LEUnsafe reads uint32 from the buffer at current offset in little-endian
// and moves the offset forward 4 bytes.
//
// This method doesn't check the read safety.
func (b *Buffer) ReadUint32LEUnsafe() uint32 {
	var result uint32 = uint32(b.buf[b.offset]) | uint32(b.buf[b.offset+1])<<8 | uint32(b.buf[b.offset+2])<<16 | uint32(b.buf[b.offset+3])<<24
	b.offset += 4
	return result
}

// ReadUint32BE reads uint32 from the buffer at current offset in big-endian
// and moves the offset forward 4 bytes.
//
// If the offset is invalid or not safe to read data,
// it will panic with BufferOverreadError or BufferUnderreadError.
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

// ReadUint32BEUnsafe reads uint32 from the buffer at current offset in big-endian
// and moves the offset forward 4 bytes.
//
// This method doesn't check the read safety.
func (b *Buffer) ReadUint32BEUnsafe() uint32 {
	var result uint32 = uint32(b.buf[b.offset])<<24 | uint32(b.buf[b.offset+1])<<16 | uint32(b.buf[b.offset+2])<<8 | uint32(b.buf[b.offset+3])
	b.offset += 4
	return result
}

// ReadUint64LE reads uint64 from the buffer at current offset in little-endian
// and moves the offset forward 8 bytes.
//
// If the offset is invalid or not safe to read data,
// it will panic with BufferOverreadError or BufferUnderreadError.
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

// ReadUint64LEUnsafe reads uint64 from the buffer at current offset in little-endian
// and moves the offset forward 8 bytes.
//
// This method doesn't check the read safety.
func (b *Buffer) ReadUint64LEUnsafe() uint64 {
	var result uint64 = uint64(b.buf[b.offset]) | uint64(b.buf[b.offset+1])<<8 | uint64(b.buf[b.offset+2])<<16 | uint64(b.buf[b.offset+3])<<24 | uint64(b.buf[b.offset+4])<<32 | uint64(b.buf[b.offset+5])<<40 | uint64(b.buf[b.offset+6])<<48 | uint64(b.buf[b.offset+7])<<56
	b.offset += 8
	return result
}

// ReadUint64BE reads uint64 from the buffer at current offset in big-endian
// and moves the offset forward 8 bytes.
//
// If the offset is invalid or not safe to read data,
// it will panic with BufferOverreadError or BufferUnderreadError.
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

// ReadUint64BEUnsafe reads uint64 from the buffer at current offset in big-endian
// and moves the offset forward 8 bytes.
//
// This method doesn't check the read safety.
func (b *Buffer) ReadUint64BEUnsafe() uint64 {
	var result uint64 = uint64(b.buf[b.offset])<<56 | uint64(b.buf[b.offset+1])<<48 | uint64(b.buf[b.offset+2])<<40 | uint64(b.buf[b.offset+3])<<32 | uint64(b.buf[b.offset+4])<<24 | uint64(b.buf[b.offset+5])<<16 | uint64(b.buf[b.offset+6])<<8 | uint64(b.buf[b.offset+7])
	b.offset += 8
	return result
}

// ReadFloat32LE reads float32 from the buffer at current offset in little-endian
// and moves the offset forward 4 bytes.
//
// If the offset is invalid or not safe to read data,
// it will panic with BufferOverreadError or BufferUnderreadError.
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

// ReadFloat32LEUnsafe reads float32 from the buffer at current offset in little-endian
// and moves the offset forward 4 bytes.
//
// This method doesn't check the read safety.
func (b *Buffer) ReadFloat32LEUnsafe() float32 {
	var val uint64 = uint64(b.buf[b.offset]) | uint64(b.buf[b.offset+1])<<8 | uint64(b.buf[b.offset+2])<<16 | uint64(b.buf[b.offset+3])<<24
	var result float32 = *(*float32)(unsafe.Pointer(&val))
	b.offset += 4
	return result
}

// ReadFloat32BE reads float32 from the buffer at current offset in big-endian
// and moves the offset forward 4 bytes.
//
// If the offset is invalid or not safe to read data,
// it will panic with BufferOverreadError or BufferUnderreadError.
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

// ReadFloat32BEUnsafe reads float32 from the buffer at current offset in big-endian
// and moves the offset forward 4 bytes.
//
// This method doesn't check the read safety.
func (b *Buffer) ReadFloat32BEUnsafe() float32 {
	var val uint64 = uint64(b.buf[b.offset])<<24 | uint64(b.buf[b.offset+1])<<16 | uint64(b.buf[b.offset+2])<<8 | uint64(b.buf[b.offset+3])
	var result float32 = *(*float32)(unsafe.Pointer(&val))
	b.offset += 4
	return result
}

// ReadFloat64LE reads float64 from the buffer at current offset in little-endian
// and moves the offset forward 8 bytes.
//
// If the offset is invalid or not safe to read data,
// it will panic with BufferOverreadError or BufferUnderreadError.
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

// ReadFloat64LEUnsafe reads float64 from the buffer at current offset in little-endian
// and moves the offset forward 8 bytes.
//
// This method doesn't check the read safety.
func (b *Buffer) ReadFloat64LEUnsafe() float64 {
	var val uint64 = uint64(b.buf[b.offset]) | uint64(b.buf[b.offset+1])<<8 | uint64(b.buf[b.offset+2])<<16 | uint64(b.buf[b.offset+3])<<24 | uint64(b.buf[b.offset+4])<<32 | uint64(b.buf[b.offset+5])<<40 | uint64(b.buf[b.offset+6])<<48 | uint64(b.buf[b.offset+7])<<56
	var result float64 = *(*float64)(unsafe.Pointer(&val))
	b.offset += 8
	return result
}

// ReadFloat64BE reads float64 from the buffer at current offset in big-endian
// and moves the offset forward 8 bytes.
//
// If the offset is invalid or not safe to read data,
// it will panic with BufferOverreadError or BufferUnderreadError.
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

// ReadFloat64BEUnsafe reads float64 from the buffer at current offset in big-endian
// and moves the offset forward 8 bytes.
//
// This method doesn't check the read safety.
func (b *Buffer) ReadFloat64BEUnsafe() float64 {
	var val uint64 = uint64(b.buf[b.offset])<<56 | uint64(b.buf[b.offset+1])<<48 | uint64(b.buf[b.offset+2])<<40 | uint64(b.buf[b.offset+3])<<32 | uint64(b.buf[b.offset+4])<<24 | uint64(b.buf[b.offset+5])<<16 | uint64(b.buf[b.offset+6])<<8 | uint64(b.buf[b.offset+7])
	var result float64 = *(*float64)(unsafe.Pointer(&val))
	b.offset += 8
	return result
}
