package registers

// FloatRegister ...
type FloatRegister struct {
	IntRegister
}

func (f *FloatRegister) GetFloat() float32 {
	// Format the float properly
	return float32(f.value)

}

func (f *FloatRegister) SetFloat(v float32) {
	// Format the float properly
	f.value = int32(v)
}
