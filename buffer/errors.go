package buffer

import (
	"fmt"

	"github.com/pkg/errors"
)

// Error implements custom error in buffer
type BufferError struct {
	error string
}

// Error formats BufferError as string
func (e BufferError) Error() string {
	return fmt.Sprintf("buf: %v", e.error)
}

var (
	// BufferOverreadError represents an instance in which a read
	// attempted to read past the buffer itself
	BufferSeekError = errors.WithStack(&BufferError{
		error: "seek offset is invalid",
	})
	// BufferOverreadError represents an instance in which a read
	// attempted to read past the buffer itself
	BufferOverreadError = errors.WithStack(&BufferError{
		error: "read exceeds buffer capacity",
	})

	// BufferUnderreadError represents an instance in which a read
	// attempted to read before the buffer itself
	BufferUnderreadError = errors.WithStack(&BufferError{
		error: "read offset is less than zero",
	})

	// BufferOverwriteError represents an instance in which a write
	// attempted to write past the buffer itself
	BufferOverwriteError = errors.WithStack(&BufferError{
		error: "write offset exceeds buffer capacity",
	})

	// BufferUnderwriteError represents an instance in which a write
	// attempted to write before the buffer itself
	BufferUnderwriteError = errors.WithStack(&BufferError{
		error: "write offset is less than zero",
	})

	// BufferInvalidByteCountError represents an instance in which an
	// invalid byte count was passed to one of the buffer's methods
	BufferInvalidByteCountError = errors.WithStack(&BufferError{
		error: "invalid byte count requested",
	})

	// BytesBufNegativeReadError represents an instance in which a
	// reader returned a negative count from its Read method
	BytesBufNegativeReadError = errors.WithStack(&BufferError{
		error: "reader returned negative count from Read",
	})
)
