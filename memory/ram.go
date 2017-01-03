package memory

import (
	//"fmt"
	"github.com/uroshercog/sic-machine/obj"
	"fmt"
)

// MaxAddress ...
const MaxAddress = 1000 // bytes

const (
	errInvalidMemoryAddress = "Invalid memory address"
)

// RAM ...
type RAM struct {
	cells []byte
}

// GetByte ...
func (ram *RAM) GetByte(addr int32) (byte) {
	ram.validAddress(addr)
	return ram.cells[addr]
}

// SetByte ...
func (ram *RAM) SetByte(addr int32, value byte) (err error) {
	ram.validAddress(addr)
	ram.cells[addr] = value
	return
}

// GetWord ...
func (ram *RAM) GetWord(addr int32) (ret int32) {
	ram.validAddress(addr)
	//fmt.Printf("* Loading from %#x\n", addr)
	val := ram.GetByte(addr)
	//fmt.Printf("\t* Loaded %#x\n", val)
	ret = int32(val)
	//fmt.Printf("* Loading from %#x\n", addr+1)
	val = ram.GetByte(addr + 1)
	//fmt.Printf("\t* Loaded %#x\n", val)
	ret = (ret << 3) + int32(val)
	//fmt.Printf("* Loading from %#x\n", addr+2)
	val = ram.GetByte(addr + 2)
	//fmt.Printf("\t* Loaded %#x\n", val)
	ret = (ret << 3) + int32(val)
	return
}

// SetWord ...
func (ram *RAM) SetWord(addr int32, value int32) {
	ram.validAddress(addr + 2)
	ram.SetByte(addr, byte((value & 0xFF0000) >> 16))
	ram.SetByte(addr+1, byte((value & 0xFF00) >> 8))
	ram.SetByte(addr+2, byte(value & 0xFF))
}

// Load ...
func (ram *RAM) Load(objCode *obj.ObjectCode) {
	for _, body := range objCode.Code {
		for i, code := range body.Code {
			addr := body.StartAddr + int32(i)
			ram.validAddress(addr)
			ram.cells[addr] = code
		}
	}
}

func (ram *RAM) validAddress(addr int32) {
	if addr < 0 || addr >= MaxAddress {
		panic(fmt.Errorf(fmt.Sprintf("%s %s", errInvalidMemoryAddress, "%#x"), addr))
	}
}

func (ram *RAM) GetRaw() []byte {
	return ram.cells
}

// New ...
func New() *RAM {
	return &RAM{
		cells: make([]byte, MaxAddress),
	}
}
