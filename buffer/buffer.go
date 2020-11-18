// Package buffer provides a easy way to manipulate byte slices.
//
// It's always annoying to calculate offset, ensure buffer is large enough ...etc,
// this package wraps annoying operations and lets programmer happy.
package buffer

import (
	"encoding/binary"
	"unsafe"
)

// Buffer is a variable-sized buffer of bytes, provides methods for generic buffer operations
type Buffer struct {
	buf    []byte
	offset int64
	cap    int64
	pbuf   unsafe.Pointer
}

// Create a new Buffer with buffer size
func Create(size int64) *Buffer {
	var cap int64 = size
	if cap < 64 {
		cap = 64
	}
	buf := Buffer{
		buf:    make([]byte, size, cap),
		offset: 0x0,
	}
	buf.update()
	return &buf
}

// From returns a buffer instance with the byte slice provided.
func From(data []byte) *Buffer {
	buf := Buffer{
		buf:    data,
		offset: 0x0,
	}

	buf.update()
	return &buf
}

// Bytes returns internal byte slice of the buffer
func (b *Buffer) Bytes() []byte {
	return b.buf
}

// Capacity returns the capacity of the buffer
func (b *Buffer) Capacity() int64 {
	return b.cap
}

// Offset returns current offset of the buffer
func (b *Buffer) Offset() int64 {
	return b.offset
}

// Grow grows the buffer's capacity, if necessary, to guarantee space for another n bytes.
func (b *Buffer) Grow(n int64) {
	if n < 0 {
		panic(BufferInvalidByteCountError)

	}
	bufCap := int64(cap(b.buf))
	if n <= bufCap-b.cap {
		b.buf = b.buf[:b.cap+n]
		b.update()
		return
	}
	// grow capacity twice larger than length to decrese freq. of copy
	newBuf := make([]byte, b.cap+n, (bufCap+n)*2)
	copy(newBuf, b.buf)
	b.buf = newBuf
	b.update()

}

// Seal returns the slice of bytes which seals to current offset
func (b *Buffer) Seal() []byte {
	return b.buf[:b.offset]
}

// EnsureCap ensures there are enough capacity to access, to guarantee space for another n bytes.
func (b *Buffer) EnsureCap(n int64) {
	if b.cap-b.offset >= n {
		return
	}
	b.Grow(b.offset + n - b.cap)
}

// Seek seeks to the offset of the buffer relative to current position or seeks to absolute position.
func (b *Buffer) Seek(offset int64, relative bool) {
	var pos int64
	if relative {
		pos = b.offset + offset
	} else {
		pos = offset
	}
	if pos < 0 || pos >= b.cap {
		panic(BufferSeekError)
	}
	b.offset = pos
}

// SeekUnsafe seeks to the offset of the buffer relative to current position or seeks to absolute position.
// without boundary checking.
func (b *Buffer) SeekUnsafe(offset int64, relative bool) {
	if relative {
		b.offset += offset
	} else {
		b.offset = offset
	}
}

// WriteByte writes a byte to the buffer at current offset
// and moves the offset forward 1 byte.
//
// This method will ensure buffer capacity is enough to write data.
//
// This method returns current Buffer instance for chainable operation.
func (b *Buffer) WriteByte(data byte) *Buffer {
	b.EnsureCap(1)
	b.buf[b.offset] = data
	b.offset += 1
	return b
}

// WriteBytes writes slice of byte to the buffer at current offset
// and moves the offset forward the amount of bytes written.
//
// This method will ensure buffer capacity is enough to write data.
//
// This method returns current Buffer instance for chainable operation.
func (b *Buffer) WriteBytes(data []byte) *Buffer {
	var writeLen int64 = int64(len(data))
	b.EnsureCap(writeLen)

	var ptr = unsafe.Pointer(uintptr(b.pbuf) + uintptr(b.offset))
	var i = int64(0)
	for ; i < writeLen; i++ {
		*(*byte)(unsafe.Pointer(uintptr(ptr) + uintptr(i))) = data[i]
	}
	b.offset += writeLen
	return b
}

// ReadByte reads a byte from the buffer at current offset
// and moves the offset forward 1 byte.
//
// If the offset is invalid or not safe to read data,
// it will panic with BufferOverreadError or BufferUnderreadError.
func (b *Buffer) ReadByte() byte {
	if b.offset >= b.cap {
		panic(BufferOverreadError)
	}
	if b.offset < 0 {
		panic(BufferUnderreadError)
	}
	b.offset += 1
	return b.buf[b.offset-1]
}

// ReadByteUnsafe reads a byte from the buffer at current offset
// and moves the offset forward 1 byte.
//
// This method doesn't check the read safety.
func (b *Buffer) ReadByteUnsafe() byte {
	b.offset += 1
	return b.buf[b.offset-1]
}

// ReadByte reads n bytes from the buffer at current offset
// and moves the offset forward n byte.
//
// If the offset is invalid or not safe to read data,
// it will panic with BufferOverreadError or BufferUnderreadError.
func (b *Buffer) ReadBytes(n int64) []byte {
	if (b.offset + n) > b.cap {
		panic(BufferOverreadError)
	}
	if b.offset < 0 {
		panic(BufferUnderreadError)
	}
	d := b.buf[b.offset : b.offset+n]
	b.offset += n
	return d
}

// ReadBytesUnsafe reads n bytes from the buffer at current offset
// and moves the offset forward n byte.
//
// This method doesn't check the read safety.
func (b *Buffer) ReadBytesUnsafe(n int64) []byte {
	d := b.buf[b.offset : b.offset+n]
	b.offset += n
	return d
}

/******************
 * Varint Routine *
 ******************/
// ByteLenVarUint returns the byte length of variant integer number
func ByteLenVarUint(x uint64) int64 {
	switch {
	case x < 128: // 2^7
		return 1
	case x < 16384: // 2^14
		return 2
	case x < 2097152: // 2^21
		return 3
	case x < 268435456: // 2^28
		return 4
	case x < 34359738368:
		return 5
	case x < 4398046511104:
		return 6
	case x < 562949953421312:
		return 7
	case x < 72057594037927936:
		return 8
	case x < 9223372036854775808:
		return 9
	}
	return 10
}

// ByteLenVarInt returns the byte length of unsigned variant integer number
func ByteLenVarInt(x int64) int64 {
	var y uint64
	if x >= 0 {
		y = uint64(x) * 2
	} else {
		y = uint64(x*-2) - 1
	}
	return ByteLenVarUint(y)
}

func writeVarUint(b []byte, v uint64) int64 {
	switch {
	case v < 1<<7:
		b[0] = byte(v)
		return 1
	case v < 1<<14:
		b[0] = byte((v>>0)&0x7f | 0x80)
		b[1] = byte(v >> 7)
		return 2
	case v < 1<<21:
		b[0] = byte((v>>0)&0x7f | 0x80)
		b[1] = byte((v>>7)&0x7f | 0x80)
		b[2] = byte(v >> 14)
		return 3
	case v < 1<<28:
		b[0] = byte((v>>0)&0x7f | 0x80)
		b[1] = byte((v>>7)&0x7f | 0x80)
		b[2] = byte((v>>14)&0x7f | 0x80)
		b[3] = byte(v >> 21)
		return 4
	case v < 1<<35:
		b[0] = byte((v>>0)&0x7f | 0x80)
		b[1] = byte((v>>7)&0x7f | 0x80)
		b[2] = byte((v>>14)&0x7f | 0x80)
		b[3] = byte((v>>21)&0x7f | 0x80)
		b[4] = byte(v >> 28)
		return 5
	case v < 1<<42:
		b[0] = byte((v>>0)&0x7f | 0x80)
		b[1] = byte((v>>7)&0x7f | 0x80)
		b[2] = byte((v>>14)&0x7f | 0x80)
		b[3] = byte((v>>21)&0x7f | 0x80)
		b[4] = byte((v>>28)&0x7f | 0x80)
		b[5] = byte(v >> 35)
		return 6
	case v < 1<<49:
		b[0] = byte((v>>0)&0x7f | 0x80)
		b[1] = byte((v>>7)&0x7f | 0x80)
		b[2] = byte((v>>14)&0x7f | 0x80)
		b[3] = byte((v>>21)&0x7f | 0x80)
		b[4] = byte((v>>28)&0x7f | 0x80)
		b[5] = byte((v>>35)&0x7f | 0x80)
		b[6] = byte(v >> 42)
		return 7
	case v < 1<<56:
		b[0] = byte((v>>0)&0x7f | 0x80)
		b[1] = byte((v>>7)&0x7f | 0x80)
		b[2] = byte((v>>14)&0x7f | 0x80)
		b[3] = byte((v>>21)&0x7f | 0x80)
		b[4] = byte((v>>28)&0x7f | 0x80)
		b[5] = byte((v>>35)&0x7f | 0x80)
		b[6] = byte((v>>42)&0x7f | 0x80)
		b[7] = byte(v >> 49)
		return 8
	case v < 1<<63:
		b[0] = byte((v>>0)&0x7f | 0x80)
		b[1] = byte((v>>7)&0x7f | 0x80)
		b[2] = byte((v>>14)&0x7f | 0x80)
		b[3] = byte((v>>21)&0x7f | 0x80)
		b[4] = byte((v>>28)&0x7f | 0x80)
		b[5] = byte((v>>35)&0x7f | 0x80)
		b[6] = byte((v>>42)&0x7f | 0x80)
		b[7] = byte((v>>49)&0x7f | 0x80)
		b[8] = byte(v >> 56)
		return 9
	default:
		b[0] = byte((v>>0)&0x7f | 0x80)
		b[1] = byte((v>>7)&0x7f | 0x80)
		b[2] = byte((v>>14)&0x7f | 0x80)
		b[3] = byte((v>>21)&0x7f | 0x80)
		b[4] = byte((v>>28)&0x7f | 0x80)
		b[5] = byte((v>>35)&0x7f | 0x80)
		b[6] = byte((v>>42)&0x7f | 0x80)
		b[7] = byte((v>>49)&0x7f | 0x80)
		b[8] = byte((v>>56)&0x7f | 0x80)
		b[9] = byte(1)
		return 10
	}
}

func writeVarInt(b []byte, x int64) int64 {
	ux := uint64(x) << 1
	if x < 0 {
		ux = ^ux
	}
	return writeVarUint(b, ux)
}

// WriteVarUint write varuint to the buffer at the current offset
// and moves the offset forward the amount of bytes written.
//
// This method returns current Buffer instance for chainable operation.
func (b *Buffer) WriteVarUint(x uint64) *Buffer {
	b.EnsureCap(10)
	n := writeVarUint(b.buf[b.offset:], x)
	b.offset += n
	return b
}

// WriteVarInt write varint to the buffer at the current offset
// and moves the offset forward the amount of bytes written.
//
// This method returns current Buffer instance for chainable operation.
func (b *Buffer) WriteVarInt(x int64) *Buffer {
	b.EnsureCap(10)
	n := writeVarInt(b.buf[b.offset:], x)
	b.offset += int64(n)
	return b
}

// ReadVarUint reads a varuint from the buffer at the current offset
// and moves the offset forward the amount of bytes read
func (b *Buffer) ReadVarUint() (uint64, int) {
	result, n := binary.Uvarint(b.buf[b.offset:])
	if n <= 0 {
		return 0, n
	}
	b.offset += int64(n)
	return result, n
}

// ReadVarInt reads a varuint from the buffer at the current offset
// and moves the offset forward the amount of bytes read
func (b *Buffer) ReadVarInt() (int64, int) {
	result, n := binary.Varint(b.buf[b.offset:])
	if n <= 0 {
		return 0, n
	}
	b.offset += int64(n)
	return result, n
}

// WriteString writes length of string and string data into buffer at
// current offset and moves the offset forward the amount of bytes written.
//
// This method returns current Buffer instance for chainable operation.
func (b *Buffer) WriteString(str *string) *Buffer {
	writeLen := int64(len(*str))
	b.WriteVarUint(uint64(writeLen))
	// b.WriteBytes([]byte(*str))
	b.WriteBytes(*(*[]byte)(unsafe.Pointer(str)))
	return b
}

// ReadString reads packed string from the buffer at the current offset
// and moves the offset forward the amount of bytes read.
func (b *Buffer) ReadString() string {
	size, _ := b.ReadVarUint()
	data := b.ReadBytes(int64(size))
	// return string(data)
	return *(*string)(unsafe.Pointer(&data)) //nolint:gosec
}

// ReadStringPtr reads packed string from the buffer at the current offset
// and moves the offset forward the amount of bytes read
// Retrurns pointer of string instance
func (b *Buffer) ReadStringPtr() *string {
	size, _ := b.ReadVarUint()
	data := b.ReadBytes(int64(size))
	return (*string)(unsafe.Pointer(&data)) //nolint:gosec
}

// update will update capacity of buffer by current byte length of buffer
func (b *Buffer) update() {
	b.cap = int64(len(b.buf))
	if len(b.buf) > 0 {
		b.pbuf = unsafe.Pointer(&b.buf[0])
	}
}
