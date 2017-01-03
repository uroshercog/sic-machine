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
	return &StdoutDevice{
		file: os.Stdin,
	}
}

// Read ...
func (od StdoutDevice) Read() (byte, error) {
	// Try and read one byte from the file
	var bytesRead []byte
	bytesReadCount, err := od.file.Read(bytesRead[:1])

	if err != nil {
		return 0, err
	}

	if bytesReadCount <= 0 {
		return 0, errors.New("No bytes read from the device")
	}
	return bytesRead[0], nil
}

// Write ...
func (od StdoutDevice) Write(value byte) error {
	if od.file == nil {
		return errors.New("File is nil")
	}

	bytesWritten, err := od.file.Write([]byte{value})

	if err != nil {
		return err
	}

	if bytesWritten <= 0 {
		return errors.New("No bytes written")
	}

	return nil
}

// Test ...
func (od StdoutDevice) Test() bool {
	return false
}

func (od StdoutDevice) close() error {
	return od.file.Close()
}
