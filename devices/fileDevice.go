package devices

import (
	"errors"
	"os"
	"fmt"
)

// FileDevice ...
type FileDevice struct {
	file *os.File
}

// NewFileDevice ..
func NewFileDevice(device byte) *FileDevice {
	if file, err := os.OpenFile(fmt.Sprintf("%#02x.dev", device),
		os.O_RDWR|os.O_CREATE|os.O_APPEND,
		os.ModePerm); err != nil {
		panic(err)
	} else {
		return &FileDevice{file}
	}
}

// Read a single byte from the file device
func (fd *FileDevice) Read() (byte, error) {
	// Try and read one byte from the file
	bytesRead := make([]byte, 1)
	if bytesReadCount, err := fd.file.Read(bytesRead[:1]); err != nil {
		return 0, err
	} else if bytesReadCount <= 0 {
		return 0, errors.New("No bytes read from the device")
	}

	return bytesRead[0], nil
}

// Write a single byte to a file device
func (fd *FileDevice) Write(value byte) error {
	if fd.file == nil {
		return errors.New("File is nil")
	}

	if bytesWritten, err := fd.file.Write([]byte{value}); err != nil {
		return err
	} else if bytesWritten <= 0 {
		return errors.New("No bytes written to the device")
	}

	return nil
}

// Test ...
func (fd *FileDevice) Test() bool {
	return true
}
