package processor

import (
	oc "github.com/uroshercog/sic-machine/opcodes"
	"github.com/uroshercog/sic-machine/memory"
	dev "github.com/uroshercog/sic-machine/devices"
	reg "github.com/uroshercog/sic-machine/processor/registers"
	"time"
	"fmt"
)

const (
	regA  = iota
	regX
	regL
	regB
	regS
	regT
	regF
	regPC
	regSW
)

const (
	nanoseconds = int64(1000 * 1000 * 1000) // nanoseconds in a second
)

var (
	cmdMap = []string{"LDA", "LDX", "LDL", "STA", "STX", "STL", "ADD", "SUB",
					  "MUL", "DIV", "COMP", "TIX", "JEQ", "JGT", "JLT", "J",
					  "AND", "OR", "JSUB", "RSUB", "LDCH", "STCH", "ADDF", "SUBF",
					  "MULF", "DIVF", "LDB", "LDS", "LDF", "LDT", "STB", "STS",
					  "STF", "STT", "COMPF", "", "ADDR", "SUBR", "MULR", "DIVR",
					  "COMPR", "SHIFTL", "SHIFTR", "RMO", "SVC", "CLEAR", "TIXR", "",
					  "FLOAT", "FIX", "NORM", "", "LPS", "STI", "RD", "WD",
					  "TD", "", "STSW", "SSK", "SIO", "HIO", "TIO", ""}
)

// CPU ...
type CPU struct {
	registers [9]reg.Register
	ram       *memory.RAM
	devices   *dev.DeviceManager
	clock     *time.Ticker
	speed     int64 // OP/s
	running   bool
	OnStart   []func()
	OnStop    []func()
	OnExec    []func(cmd string)
}

func (cpu *CPU) GetRegisters() []string {
	return []string{
		fmt.Sprintf("[A] %#x", cpu.registers[regA].Get()),
		fmt.Sprintf("[X] %#x", cpu.registers[regX].Get()),
		fmt.Sprintf("[L] %#x", cpu.registers[regL].Get()),
		fmt.Sprintf("[B] %#x", cpu.registers[regB].Get()),
		fmt.Sprintf("[S] %#x", cpu.registers[regS].Get()),
		fmt.Sprintf("[T] %#x", cpu.registers[regT].Get()),
		fmt.Sprintf("[F] %f", float32(cpu.registers[regF].Get())),
		fmt.Sprintf("[PC] %#x", cpu.registers[regPC].Get()),
		fmt.Sprintf("[SW] %#x", cpu.registers[regSW].Get()),
	}
}

func (cpu *CPU) SetStart(start int32) {
	cpu.ram.ValidAddress(start)
	cpu.registers[regPC].Set(start)
}

// Run ...
func (cpu *CPU) run() byte {
	// Sets the PC register and starts executing the commands
	cpu.running = true
	defer func() {
		cpu.running = false
	}()

	pcReg := cpu.registers[regPC]

	var command, operand byte

	// Load the first byte from the memory (from location in PC)
	command = cpu.ram.GetByte(pcReg.Get())
	// Increment the program counter by the number of bytes used in the instruction
	pcReg.Add(0x1)
	// Command is 8 bits
	if executed := cpu.executeF1(command); executed {
		// The command was format 1
		return command
	}

	// Load another byte
	operand = cpu.ram.GetByte(pcReg.Get())
	// Increment the program counter by the number of bytes used in the instruction
	pcReg.Add(0x1)

	// Check if the two bytes represent a command and an operand
	if executed := cpu.executeF2(command, int32(operand)); executed {
		// The command was format 2
		return command
	}

	// Load a third byte
	operandEx := (int32(operand) << 8) | int32(cpu.ram.GetByte(pcReg.Get()))
	// Increment the PC register by the number of bytes used in the instruction
	pcReg.Add(0x1)

	// Parse the bits in the command and the operand
	// The operand in SIC, F3 and F4 is 6 bits, last two bits are N (indirect addressing) and I (immediate addressing)
	bits := map[string]bool{
		// command: _ _ _ _ n i
		"n": (command & 0x2) > 0,
		"i": (command & 0x1) > 0,
		// operand: x b p e _ _ ...
		"x": (operandEx & 0x8000) > 0,
		"b": (operandEx & 0x4000) > 0,
		"p": (operandEx & 0x2000) > 0,
		"e": (operandEx & 0x1000) > 0,
	}

	// If bits n and i are 0 that means SIC format
	// Otherwise its F3 or F4

	if !bits["n"] && !bits["i"] {
		// SIC format
		// 8 bit operand, (bottom) 15 bit operand
		operandEx &= 0x7FFF
	} else if bits["e"] {
		// Extended -> format 4
		// Operand is (bottom) 20 bits

		// Load the 4th byte
		operandEx = (operandEx << 8) | int32(cpu.ram.GetByte(pcReg.Get()))
		// Format 4 operand is bottom 20 bits
		operandEx &= 0xFFFFF
		// Increment the PC counter by one because we loaded an additional byte as part of the operand value
		pcReg.Add(0x1)

		if executed := cpu.execute(command, operandEx, bits); executed {
			return command
		} else {
			panic("Format 4 command was not executed")
		}
	} else {
		// Format 3
		// Check if its PC relative addressing
		if bits["p"] {
			if operandEx >= 2048 {
				operandEx -= 4096
			}
			operandEx += pcReg.Get()
		} else if bits["b"] {
			// Check its base relative addressing
			operandEx += cpu.registers[regB].Get()
		} else if bits["p"] && bits["b"] {
			// It cannot be PC and base relative at the same time
			panic("Invalid addressing")
		}
	}

	if bits["x"] {
		// Check if its indexed addressing
		if bits["n"] && bits["i"] {
			operandEx += cpu.registers[regX].Get()
		} else {
			panic("Invalid addressing")
		}
	}

	// The operand is only the top 6 bits
	command &= 0xFC
	// The operand is only bottom 12 bits (offset/displacement)
	operandEx &= 0xFFF

	if executed := cpu.execute(command, operandEx, bits); !executed {
		panic("Format 3 command was not executed")
	}
	return command
}

// Starts the CPU clock
func (cpu *CPU) Start() {
	if !cpu.running {
		// Create a new ticker
		cpu.clock = time.NewTicker(time.Duration(nanoseconds / cpu.speed))

		for _, f := range cpu.OnStart {
			f()
		}

		go func() {
			for range cpu.clock.C {
				cmd := cpu.run()
				for _, f := range cpu.OnExec {
					f(cmdMap[cmd >> 2])
				}
			}
		}()
	}
}

// Pauses the cpu clock, for breakpoints or w/e
func (cpu *CPU) Stop() {
	if cpu.running {
		cpu.clock.Stop()
		for _, f := range cpu.OnStop {
			f()
		}
	}
}

func (cpu *CPU) Step() {
	if !cpu.running {
		cmd := cpu.run()
		for _, f := range cpu.OnExec {
			f(cmdMap[cmd >> 2])
		}
	}
}

func (cpu *CPU) executeF1(command byte) bool {
	switch command {
	case oc.FIX:
		// Move register F to A and convert to integer
		//f := cpu.registers[regF].(reg.FloatRegister).GetFloat()
		//cpu.registers[regA].Set(int32(f))
		panic("Not implemented")
	case oc.FLOAT:
		// Move register A to F and convert to float
		//a := cpu.registers[regA].Get()
		//cpu.registers[regB].(reg.FloatRegister).SetFloat(float32(a))
		panic("Not implemented")
	case oc.HIO:
		// Halt I/O channel no. (A)
		panic("Not implemented")
	case oc.NORM:
		// Normalize F and save to F
		panic("Not implemented")
	case oc.SIO:
		// Start I/O channel number (A). Address of channel if given in S.
		panic("Not implemented")
	case oc.TIO:
		// Test I/O channel number (A)
		panic("Not implemented")
	default:
		return false
	}
	return true
}
func (cpu *CPU) executeF2(command byte, operand int32) bool {
	v1 := (operand & 0xF0) >> 4
	v2 := operand & 0xF

	switch command {
	case oc.ADDR:
		//R2 <- (R2) + (R1)
		// Load one more byte, upper 4 bits are R1 and lower 4 bits are R2
		cpu.registers[v2].Add(cpu.registers[v1].Get())
	case oc.CLEAR:
		cpu.registers[v2].Clear()
	case oc.COMPR:
		// Load one more byte, upper 4 bits are R1 and lower 4 bits are R2
		r1 := cpu.registers[v1].Get()
		r2 := cpu.registers[v2].Get()

		rSW := cpu.registers[regSW]

		if r1 < r2 {
			rSW.Set(0x2)
		} else if r1 > r2 {
			rSW.Set(0x1)
		} else {
			rSW.Set(0x0)
		}
	case oc.DIVR:
		//R2 <- (R2) + (R1)
		// Load one more byte, upper 4 bits are R1 and lower 4 bits are R2
		cpu.registers[v2].Divide(cpu.registers[v1].Get())
	case oc.MULR:
		//R2 <- (R2) * (R1)
		// Load one more byte, upper 4 bits are R1 and lower 4 bits are R2
		cpu.registers[v2].Add(cpu.registers[v1].Get())
	case oc.RMO:
		//R2 <- (R1)
		// Load one more byte, upper 4 bits are R1 and lower 4 bits are R2
		cpu.registers[v2].Set(cpu.registers[v1].Get())
	case oc.SHIFTR:
		cpu.registers[v1].ShiftRight(uint32(v2))
	case oc.SHITFTL:
		cpu.registers[v1].ShiftLeft(uint32(v2))
	case oc.SUBR:
		cpu.registers[v2].Sub(cpu.registers[v1].Get())
	case oc.SVC:
		panic("Not implemented")
	case oc.TIXR:
		panic("Not implemented")
	default:
		return false
	}
	return true
}
func (cpu *CPU) execute(command byte, operand int32, flags map[string]bool) bool {
	switch command {
	case oc.ADD:
		// A <- (A) + (m..m + 2)
		cpu.registers[regA].Add(cpu.resolveWordOperand(operand, flags))
	case oc.ADDF:
		panic("Not implemented")
	case oc.AND:
		// A <- (A) + (m..m + 2)
		cpu.registers[regA].And(cpu.resolveWordOperand(operand, flags))
	case oc.COMP:
		sw := cpu.registers[regSW].(*reg.SwRegister)
		sw.Compare(cpu.registers[regA].Get(), cpu.resolveWordOperand(operand, flags))
	case oc.COMPF:
		panic("Not implemented")
	case oc.DIV:
		cpu.registers[regA].Divide(cpu.resolveWordOperand(operand, flags))
	case oc.DIVF:
		panic("Not implemented")
	case oc.J:
		if flags["n"] && !flags["i"] {
			operand = cpu.ram.GetWord(operand)
		}
		cpu.registers[regPC].Set(operand)
	case oc.JEQ:
		r := cpu.registers[regSW].(*reg.SwRegister)
		if r.IsEqual() {
			if flags["n"] && !flags["i"] {
				operand = cpu.ram.GetWord(operand)
			}
			cpu.registers[regPC].Set(operand)
		}
	case oc.JGT:
		r := cpu.registers[regSW].(*reg.SwRegister)
		if r.IsGreater() {
			if flags["n"] && !flags["i"] {
				operand = cpu.ram.GetWord(operand)
			}
			cpu.registers[regPC].Set(operand)
		}
		r.Clear()
	case oc.JLT:
		r := cpu.registers[regSW].(*reg.SwRegister)
		if r.IsLess() {
			if flags["n"] && !flags["i"] {
				operand = cpu.ram.GetWord(operand)
			}
			cpu.registers[regPC].Set(operand)
		}
		r.Clear()
	case oc.JSUB:
		cpu.registers[regL].Set(cpu.registers[regPC].Get())
		if flags["n"] && !flags["i"] {
			operand = cpu.ram.GetWord(operand)
		}
		cpu.registers[regPC].Set(operand)
	case oc.LDA:
		// 	Load from memory location <operand>
		cpu.registers[regA].Set(cpu.resolveWordOperand(operand, flags))
	case oc.LDB:
		// 	Load from memory location <operand>
		cpu.registers[regB].Set(cpu.resolveWordOperand(operand, flags))
	case oc.LDCH:
		// 	Load from memory location <operand>
		cpu.registers[regA].Set(int32(cpu.resolveByteOperand(operand, flags) & 0xFF))
	case oc.LDF:
		panic("Not implemented")
	case oc.LDL:
		// 	Load from memory location <operand>
		cpu.registers[regL].Set(cpu.resolveWordOperand(operand, flags))
	case oc.LDS:
		// 	Load from memory location <operand>
		cpu.registers[regS].Set(cpu.resolveWordOperand(operand, flags))
	case oc.LDT:
		// 	Load from memory location <operand>
		cpu.registers[regT].Set(cpu.resolveWordOperand(operand, flags))
	case oc.LDX:
		// 	Load from memory location <operand>
		cpu.registers[regX].Set(cpu.resolveWordOperand(operand, flags))
	case oc.LPS:
		panic("Not implemented")
	case oc.MUL:
		cpu.registers[regA].Multiply(cpu.resolveWordOperand(operand, flags))
	case oc.MULF:
		panic("Not implemented")
	case oc.OR:
		cpu.registers[regA].Or(cpu.resolveWordOperand(operand, flags))
	case oc.RD:
		if m, err := cpu.devices.Get(byte(operand)).Read(); err == nil {
			cpu.registers[regA].Set(int32(m))
		} else {
			panic(err)
		}
	case oc.RSUB:
		cpu.registers[regPC].Set(cpu.registers[regL].Get())
	case oc.SSK:
		panic("Not implemented")
	case oc.STA:
		if flags["n"] && !flags["i"] {
			operand = cpu.ram.GetWord(operand)
		}
		cpu.ram.SetWord(operand, cpu.registers[regA].Get())
	case oc.STB:
		if flags["n"] && !flags["i"] {
			operand = cpu.ram.GetWord(operand)
		}
		cpu.ram.SetWord(operand, cpu.registers[regB].Get())
	case oc.STCH:
		if flags["n"] && !flags["i"] {
			operand = cpu.ram.GetWord(operand)
		}
		cpu.ram.SetByte(operand, byte(cpu.registers[regB].Get()))
	case oc.STF:
		panic("Not implemented")
	case oc.STI:
		panic("Not implemented")
	case oc.STL:
		if flags["n"] && !flags["i"] {
			operand = cpu.ram.GetWord(operand)
		}
		cpu.ram.SetWord(operand, cpu.registers[regL].Get())
	case oc.STS:
		if flags["n"] && !flags["i"] {
			operand = cpu.ram.GetWord(operand)
		}
		cpu.ram.SetWord(operand, cpu.registers[regS].Get())
	case oc.STSW:
		if flags["n"] && !flags["i"] {
			operand = cpu.ram.GetWord(operand)
		}
		cpu.ram.SetWord(operand, cpu.registers[regSW].Get())
	case oc.STT:
		if flags["n"] && !flags["i"] {
			operand = cpu.ram.GetWord(operand)
		}
		cpu.ram.SetWord(operand, cpu.registers[regT].Get())
	case oc.STX:
		if flags["n"] && !flags["i"] {
			operand = cpu.ram.GetWord(operand)
		}
		cpu.ram.SetWord(operand, cpu.registers[regX].Get())
	case oc.SUB:
		//R2 <- (R2) - (R1)
		// Load one more byte, upper 4 bits are R1 and lower 4 bits are R2
		cpu.registers[regA].Sub(cpu.resolveWordOperand(operand, flags))
	case oc.SUBF:
		panic("Not implemented")
	case oc.TD:
		panic("Not implemented")
	case oc.TIX:
		panic("Not implemented")
	case oc.WD:
		cpu.devices.Get(cpu.resolveByteOperand(operand, flags)).Write(byte(cpu.registers[regA].Get()))
	default:
		return false
	}
	return true
}

// Takes an operand and determines the actual value of the operand
func (cpu *CPU) resolveWordOperand(operand int32, flags map[string]bool) int32 {
	if flags["i"] && !flags["n"] {
		return operand
	}

	operand = cpu.ram.GetWord(operand)

	if flags["n"] && !flags["i"] {
		operand = cpu.ram.GetWord(operand)
	}

	return operand
}

func (cpu *CPU) resolveByteOperand(operand int32, flags map[string]bool) byte {
	if flags["i"] && !flags["n"] {
		return byte(operand & 0xFF)
	}

	if flags["n"] && !flags["i"] {
		return cpu.ram.GetByte(cpu.ram.GetWord(operand))
	}

	return cpu.ram.GetByte(operand)
}

// New ...
func NewCPU(ram *memory.RAM, devices *dev.DeviceManager) *CPU {
	registers := [...]reg.Register{
		&reg.IntRegister{},
		&reg.IntRegister{},
		&reg.IntRegister{},
		&reg.IntRegister{},
		&reg.IntRegister{},
		&reg.IntRegister{},
		&reg.FloatRegister{},
		&reg.IntRegister{},
		&reg.SwRegister{},
	}

	ops := int64(1000)

	return &CPU{
		registers: registers,
		ram:       ram,
		devices:   devices,
		clock:     time.NewTicker(time.Duration(nanoseconds / ops)),
		speed:     ops,
		OnStart:   []func(){},
		OnStop:    []func(){},
		OnExec:    []func(cmd string){},
	}
}
