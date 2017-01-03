package registers

type SwRegister struct {
	IntRegister
}

func (sw *SwRegister) IsEqual() bool {
	return sw.value&0x40 > 0
}

func (sw *SwRegister) IsLess() bool {
	return sw.value&0x20 > 0
}

func (sw *SwRegister) IsGreater() bool {
	return sw.value&0x80 > 0
}

func (sw *SwRegister) Compare(a, b int32) {
	sw.Clear()
	if a < b {
		sw.value = 0x20
	} else if a > b {
		sw.value = 0x80
	} else {
		sw.value = 0x40
	}
}
