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
	return &StderrDevice{
		file: os.Stderr,
	}
}

// Read ...
func (se StderrDevice) Read() (byte, error) {
	// Try and read one byte from the file
	var bytesRead []byte
	bytesReadCount, err := se.file.Read(bytesRead[:1])

	if err != nil {
		return 0, err
	}

	if bytesReadCount <= 0 {
		return 0, errors.New("No bytes read from the device")
	}
	return bytesRead[0], nil
}

// Write ...
func (se StderrDevice) Write(value byte) error {
	if se.file == nil {
		return errors.New("File is nil")
	}

	bytesWritten, err := se.file.Write([]byte{value})

	if err != nil {
		return err
	}

	if bytesWritten <= 0 {
		return errors.New("No bytes written")
	}

	return nil
}

// Test ...
func (se StderrDevice) Test() bool {
	return false
}

func (se StderrDevice) close() error {
	return se.file.Close()
}
