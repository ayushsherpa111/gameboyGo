package opcodes

import (
	"github.com/ayushsherpa111/gameboyEMU/cpu"
	"github.com/ayushsherpa111/gameboyEMU/instructions"
)

func NewOpcodeStore(cpu *cpu.CPU) []instructions.Instruction {
	opStore := make([]instructions.Instruction, 0xFF)
	opStore[0x00] = NewNOP("NOP", 0x00)
	opStore[0x01] = NewLD("LD BC, u16", 0x01, cpu)
	return opStore
}
