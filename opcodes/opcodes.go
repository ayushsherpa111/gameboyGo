package opcodes

import (
	"github.com/ayushsherpa111/gameboyEMU/cpu"
	"github.com/ayushsherpa111/gameboyEMU/instructions"
)

func NewOpcodeStore(cpu *cpu.CPU) [0xFF]instructions.Instruction {
	NOP := NewNOP()
	STOP := NewStop()
	LD := NewLD(cpu)
	INC := NewInc(cpu)
	DEC := NewDEC(cpu)
	JR := NewJR(cpu)
	ADD := NewADD(cpu)
	JP := NewJP(cpu)
	RR := NewRR(cpu)
	RL := NewRl(cpu)
	DAA := NewDAA(cpu)
	CPL := NewCPL(cpu)
	SCF := NewSCF(cpu)

	opStore := [0xFF]instructions.Instruction{
		NOP,  // 0x00
		LD,   // 0x01
		LD,   // 0x02
		INC,  // 0x03
		INC,  // 0x04
		DEC,  // 0x05
		LD,   // 0x06
		RL,   // 0x07
		LD,   // 0x08
		ADD,  // 0x09
		LD,   // 0x0a
		DEC,  // 0x0b
		INC,  // 0x0c
		DEC,  // 0x0d
		LD,   // 0x0e
		RR,   // 0x0f
		STOP, // 0x10
		LD,   // 0x11
		LD,   // 0x12
		INC,  // 0x13
		INC,  // 0x14
		DEC,  // 0x15
		LD,   // 0x16
		RL,   // 0x17
		JR,   // 0x18
		ADD,  // 0x19
		LD,   // 0x1A
		DEC,  // 0x1B
		INC,  // 0x1C
		DEC,  // 0x1D
		LD,   // 0x1E
		RR,   // 0x1F
		JR,   // 0x20
		LD,   // 0x21
		LD,   // 0x22
		INC,  // 0x23
		INC,  // 0x24
		DEC,  // 0x25
		LD,   // 0x26
		DAA,  // 0x27
		JR,   // 0x28
		ADD,  // 0x29
		LD,   // 0x2A
		DEC,  // 0x2B
		INC,  // 0x2C
		DEC,  // 0x2D
		LD,   // 0x2E
		CPL,  // 0x2F
		JR,   // 0x30
		LD,   // 0x31
		LD,   // 0x32
		INC,  // 0x33
		INC,  // 0x34
		DEC,  // 0x35
		LD,   // 0x36
		SCF,  // 0x37
		JR,   // 0x38
		ADD,  // 0x39
		LD,   // 0x3a
		DEC,  // 0x3b
		INC,  // 0x3c
		DEC,  // 0x3d
		LD,   // 0x3e
		nil,  // 0x3f
	}
	opStore[0xC3] = JP
	return opStore
}
