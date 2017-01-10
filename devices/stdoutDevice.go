package devices

import (
	"errors"
	"os"
)

// StdoutDevice ...
type StdoutDevice struct {
	file *os.File
}

// NewStdoutDevice ..
func NewStdoutDevice() *StdoutDevice {
	return &StdoutDevice{os.Stdin}
}

// Read ...
func (od StdoutDevice) Read() (byte, error) {
	if od.file == nil {
		return 0, errors.New("File is nil")
	}
	// Try and read one byte from the file
	bytesRead := make([]byte, 1)
	if bytesReadCount, err := od.file.Read(bytesRead[:1]); err != nil {
		return 0, err
	} else if bytesReadCount <= 0 {
		return 0, errors.New("No bytes read from the device")
	}
	return bytesRead[0], nil
}

// Write ...
func (od StdoutDevice) Write(value byte) error {
	if od.file == nil {
		return errors.New("File is nil")
	}

	if bytesWritten, err := od.file.Write([]byte{value}); err != nil {
		return err
	} else if bytesWritten <= 0 {
		return errors.New("No bytes written")
	}
	return nil
}

// Test ...
func (od StdoutDevice) Test() bool {
	return true
}
