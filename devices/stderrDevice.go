package devices

import (
	"errors"
	"os"
)

// StderrDevice ...
type StderrDevice struct {
	file *os.File
}

// NewStderrDevice ..
func NewStderrDevice() *StderrDevice {
	return &StderrDevice{os.Stderr}
}

// Read ...
func (se StderrDevice) Read() (byte, error) {
	if se.file == nil {
		return 0, errors.New("File is nil")
	}
	// Try and read one byte from the file
	bytesRead := make([]byte, 1)
	if bytesReadCount, err := se.file.Read(bytesRead[:1]); err != nil {
		return 0, err
	} else if bytesReadCount <= 0 {
		return 0, errors.New("No bytes read from the device")
	}
	return bytesRead[0], nil
}

// Write ...
func (se StderrDevice) Write(value byte) error {
	if se.file == nil {
		return errors.New("File is nil")
	}

	if bytesWritten, err := se.file.Write([]byte{value}); err != nil {
		return err
	} else if bytesWritten <= 0 {
		return errors.New("No bytes written")
	}
	return nil
}

// Test ...
func (se StderrDevice) Test() bool {
	return true
}
