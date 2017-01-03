package devices

import (
	"errors"
	"os"
)

// FileDevice ...
type FileDevice struct {
	file *os.File
}

// NewFileDevice ..
func NewFileDevice(filename string) *FileDevice {
	file, err := os.OpenFile(filename, os.O_RDWR, os.ModeExclusive|os.ModeDevice)

	if err != nil {
		panic(err)
	}

	return &FileDevice{
		file: file,
	}
}

// Read ...
func (fd FileDevice) Read() (byte, error) {
	// Try and read one byte from the file
	var bytesRead []byte
	bytesReadCount, err := fd.file.Read(bytesRead[:1])

	if err != nil {
		return 0, err
	}

	if bytesReadCount <= 0 {
		return 0, errors.New("No bytes read from the device")
	}
	return bytesRead[0], nil
}

// Write ...
func (fd FileDevice) Write(value byte) error {
	if fd.file == nil {
		return errors.New("File is nil")
	}

	bytesWritten, err := fd.file.Write([]byte{value})

	if err != nil {
		return err
	}

	if bytesWritten <= 0 {
		return errors.New("No bytes written")
	}

	return nil
}

// Test ...
func (fd FileDevice) Test() bool {
	return false
}
