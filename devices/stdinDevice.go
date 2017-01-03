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
	return &StdinDevice{
		file: os.Stdin,
	}
}

// Read ...
func (id StdinDevice) Read() (byte, error) {
	// Try and read one byte from the file
	var bytesRead []byte
	bytesReadCount, err := id.file.Read(bytesRead[:1])

	if err != nil {
		return 0, err
	}

	if bytesReadCount <= 0 {
		return 0, errors.New("No bytes read from the device")
	}
	return bytesRead[0], nil
}

// Write ...
func (id StdinDevice) Write(value byte) error {
	if id.file == nil {
		return errors.New("File is nil")
	}

	bytesWritten, err := id.file.Write([]byte{value})

	if err != nil {
		return err
	}

	if bytesWritten <= 0 {
		return errors.New("No bytes written")
	}

	return nil
}

// Test ...
func (id StdinDevice) Test() bool {
	return false
}

func (id StdinDevice) close() error {
	return id.file.Close()
}
