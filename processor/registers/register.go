package registers

type Register interface {
	Get() int32
	Set(int32)
	Add(int32)
	Sub(int32)
	Clear()
	Multiply(int32)
	Divide(int32)
	ShiftLeft(uint32)
	ShiftRight(uint32)
	And(int32)
	Or(int32)
}
