package memory

import (
	//"fmt"
	"github.com/uroshercog/sic-machine/obj"
	"fmt"
)

// MaxAddress ...
const MaxAddress = 0xF000 // bytes

const (
	errInvalidMemoryAddress = "Invalid memory address"
)

// RAM ...
type RAM struct {
	cells []byte
}

// GetByte ...
func (ram *RAM) GetByte(addr int32) (byte) {
	ram.ValidAddress(addr)
	return ram.cells[addr]
}

// SetByte ...
func (ram *RAM) SetByte(addr int32, value byte) (err error) {
	ram.ValidAddress(addr)
	ram.cells[addr] = value
	return
}

// GetWord ...
func (ram *RAM) GetWord(addr int32) (ret int32) {
	ram.ValidAddress(addr + 2)

	val := ram.GetByte(addr)
	ret = int32(val)
	val = ram.GetByte(addr + 1)
	ret = (ret << 8) + int32(val)
	val = ram.GetByte(addr + 2)
	ret = (ret << 8) + int32(val)
	return
}

// SetWord ...
func (ram *RAM) SetWord(addr int32, value int32) {
	ram.ValidAddress(addr + 2)
	ram.SetByte(addr, byte((value & 0xFF0000) >> 16))
	ram.SetByte(addr+1, byte((value & 0xFF00) >> 8))
	ram.SetByte(addr+2, byte(value & 0xFF))
}

// Load ...
func (ram *RAM) Load(objCode *obj.ObjectCode) {
	for _, body := range objCode.Code {
		for i, code := range body.Code {
			addr := body.StartAddr + int32(i)
			ram.ValidAddress(addr)
			ram.cells[addr] = code
		}
	}
}

func (ram *RAM) ValidAddress(addr int32) {
	if addr < 0 || addr >= MaxAddress {
		panic(fmt.Errorf(fmt.Sprintf("%s %s", errInvalidMemoryAddress, "%#x"), addr))
	}
}

func (ram *RAM) GetRaw() []byte {
	return ram.cells
}

// New ...
func New() *RAM {
	return &RAM{make([]byte, MaxAddress)}
}
