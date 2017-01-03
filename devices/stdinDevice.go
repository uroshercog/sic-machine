package devices

import (
	"errors"
	"os"
)

// StdinDevice ...
type StdinDevice struct {
	file *os.File
}

// NewStdinDevice ..
func NewStdinDevice() *StdinDevice {
	return &StdinDevice{os.Stdin}
}

// Read ...
func (id *StdinDevice) Read() (byte, error) {
	// Try and read one byte from the file
	if id.file == nil {
		return 0, errors.New("File is nil")
	}

	var bytesRead []byte
	if bytesReadCount, err := id.file.Read(bytesRead[:1]); err != nil {
		return 0, err
	} else if bytesReadCount <= 0 {
		return 0, errors.New("No bytes read from the device")
	}
	return bytesRead[0], nil
}

// Write ...
func (id *StdinDevice) Write(value byte) error {
	if id.file == nil {
		return errors.New("File is nil")
	}

	if bytesWritten, err := id.file.Write([]byte{value}); err != nil {
		return err
	} else if bytesWritten <= 0 {
		return errors.New("No bytes written")
	}
	return nil
}

// Test ...
func (id *StdinDevice) Test() bool {
	return true
}
