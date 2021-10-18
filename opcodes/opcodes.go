package opcodes

import (
	"github.com/ayushsherpa111/gameboyEMU/cpu"
	"github.com/ayushsherpa111/gameboyEMU/instructions"
)

func NewOpcodeStore(cpu *cpu.CPU) [0xFF]instructions.Instruction {
	NOP := NewNOP("NOP", 0x0)
	LD := NewLD(cpu)
	INC := NewInc(cpu)
	DEC := NewDEC(cpu)
	JR := NewJR(cpu)
	opStore := [0xFF]instructions.Instruction{
		NOP, // 0x00
		LD,  // 0x01
		LD,  // 0x02
		INC, // 0x03
		INC, // 0x04
		DEC, // 0x05
		LD,  // 0x06
		JR,  // FIX
	}
	return opStore
}
