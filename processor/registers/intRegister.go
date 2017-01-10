package registers

// IntRegister ...
type IntRegister struct {
	value int32
}

// Get ...
func (reg *IntRegister) Get() int32 {
	return reg.value
}

// Set ...
func (reg *IntRegister) Set(value int32) {
	reg.value = value
}

// Add ...
func (reg *IntRegister) Add(value int32) {
	reg.value += value
}

// Sub ...
func (reg *IntRegister) Sub(value int32) {
	reg.value -= value
}

// Clear ...
func (reg *IntRegister) Clear() {
	reg.value = 0
}

// Multiply ...
func (reg *IntRegister) Multiply(value int32) {
	reg.value *= value
}

// Divide ...
func (reg *IntRegister) Divide(value int32) {
	reg.value /= value
}

// ShiftLeft ...
func (reg *IntRegister) ShiftLeft(bitCount uint32) {
	reg.value <<= bitCount
}

// ShiftRight ...
func (reg *IntRegister) ShiftRight(bitCount uint32) {
	reg.value >>= bitCount
}

// And ...
func (reg *IntRegister) And(value int32) {
	reg.value &= value
}

// Or ...
func (reg *IntRegister) Or(value int32) {
	reg.value |= value
}
