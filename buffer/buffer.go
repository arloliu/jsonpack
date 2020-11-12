//nolint:gosec
package buffer

import (
	"encoding/binary"
	"unsafe"
)

// Buffer a buffer instance
type Buffer struct {
	// crunch.Buffer
	buf    []byte
	offset int64
	cap    int64
	pbuf   unsafe.Pointer
}

// Create a new Buffer with buffer size
func Create(size int64) *Buffer {
	buf := Buffer{
		buf:    make([]byte, size),
		offset: 0x0,
	}
	buf.Update()
	return &buf
}

func From(data []byte) *Buffer {
	buf := Buffer{
		buf:    data,
		offset: 0x0,
	}

	buf.Update()
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

// Grow makes the buffer's capacity bigger by n bytes
func (b *Buffer) Grow(n int64) {
	if n < 0 {
		panic(BufferInvalidByteCountError)

	}
	bufCap := int64(cap(b.buf))
	if n <= bufCap-b.cap {
		b.buf = b.buf[:b.cap+n]
		b.Update()
		return
	}
	// grow capacity twice larger than length to decrese freq. of copy
	newBuf := make([]byte, b.cap+n, (bufCap+n)*2)
	copy(newBuf, b.buf)
	b.buf = newBuf
	b.Update()

}

// Seal return slice which seals to current offset
func (b *Buffer) Seal() []byte {
	return b.buf[:b.offset]
}

// EnsureCap ensure there are enough capacity to access
func (b *Buffer) EnsureCap(n int64) {
	if b.cap-b.offset >= n {
		return
	}
	b.Grow(b.offset + n - b.cap)
}

// Seek seeks to position off of the buffer
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

// UnsafeSeek seeks to position off of the buffer without boundary checking
func (b *Buffer) UnsafeSeek(offset int64, relative bool) {
	if relative {
		b.offset += offset
	} else {
		b.offset = offset
	}
}

// WriteByte writes byte to the buffer at the current offset and moves
// forward the amount of bytes written
//nolint
func (b *Buffer) WriteByte(data byte) {
	b.EnsureCap(1)
	b.buf[b.offset] = data
	b.UnsafeSeek(1, true)
}

// WriteBytes writes bytes to the buffer at the current offset and moves
// forward the amount of bytes written
func (b *Buffer) WriteBytes(data []byte) {
	var writeLen int64 = int64(len(data))
	b.EnsureCap(writeLen)

	var ptr = unsafe.Pointer(uintptr(b.pbuf) + uintptr(b.offset))
	var i = int64(0)
	for ; i < writeLen; i++ {
		*(*byte)(unsafe.Pointer(uintptr(ptr) + uintptr(i))) = data[i]
	}
	b.offset += writeLen
}

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

/******************
 * Varint Routine *
 ******************/
// ByteLenVarUint get size of varuint number
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

// ByteLenVarInt get size of varint number
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
// Returns number of bytes written
func (b *Buffer) WriteVarUint(x uint64) int64 {
	b.EnsureCap(10)
	// n := binary.PutUvarint(b.buf[b.offset:], x)
	n := writeVarUint(b.buf[b.offset:], x)
	b.UnsafeSeek(int64(n), true)

	return int64(n)
}

// WriteVarInt write varint to the buffer at the current offset
// and moves the offset forward the amount of bytes written.
// Returns number of bytes written
func (b *Buffer) WriteVarInt(x int64) int64 {
	b.EnsureCap(10)
	// n := binary.PutVarint(b.buf[b.offset:], x)
	n := writeVarInt(b.buf[b.offset:], x)
	b.offset += int64(n)
	// b.UnsafeSeek(int64(n), true)

	return int64(n)
}

// ReadVarUint reads a varuint from the buffer at the current offset
// and moves the offset forward the amount of bytes read
func (b *Buffer) ReadVarUint() (uint64, int) {
	result, n := binary.Uvarint(b.buf[b.offset:])
	if n <= 0 {
		return 0, n
	}
	b.UnsafeSeek(int64(n), true)
	return result, n
}

// ReadVarInt reads a varuint from the buffer at the current offset
// and moves the offset forward the amount of bytes read
func (b *Buffer) ReadVarInt() (int64, int) {
	result, n := binary.Varint(b.buf[b.offset:])
	if n <= 0 {
		return 0, n
	}
	b.UnsafeSeek(int64(n), true)
	return result, n
}

func (b *Buffer) WriteString(str *string) {
	writeLen := int64(len(*str))
	b.WriteVarUint(uint64(writeLen))
	// b.WriteBytes([]byte(*str))
	b.WriteBytes(*(*[]byte)(unsafe.Pointer(str)))
}

// ReadString reads packed string from the buffer at the current offset
// and moves the offset forward the amount of bytes read.
// Retrurns string instance
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

func (b *Buffer) Update() {
	b.cap = int64(len(b.buf))
	if len(b.buf) > 0 {
		b.pbuf = unsafe.Pointer(&b.buf[0])
	}
}

// func max(x, y int64) int64 {
// 	if x > y {
// 		return x
// 	}
// 	return y
// }
