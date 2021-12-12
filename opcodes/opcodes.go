package opcodes

import (
	"github.com/ayushsherpa111/gameboyEMU/cpu"
	instructions "github.com/ayushsherpa111/gameboyEMU/interfaces"
)

func NewOpcodeStore(cpu *cpu.CPU) [0xFF]instructions.Instruction {
	NOP := NewNOP()
	STOP := NewStop()
	LD := NewLD(cpu)
	INC := NewInc(cpu)
	DEC := NewDEC(cpu)
	JR := NewJR(cpu)
	ADD := NewADD(cpu)
	ADC := NewADC(cpu)
	OR := NewOR(cpu)
	AND := NewAND(cpu)
	SUB := NewSub(cpu)
	CP := NewCP(cpu)
	SBC := NewSBC(cpu)
	XOR := NewXOR(cpu)
	JP := NewJP(cpu)
	RR := NewRR(cpu)
	RL := NewRl(cpu)
	DAA := NewDAA(cpu)
	CPL := NewCPL(cpu)
	SCF := NewSCF(cpu)
	RET := NewRet(cpu)
	POP := NewPOP(cpu)
	CB := NewCB(cpu)
	PUSH := NewPush(cpu)
	CALL := NewCall(cpu)

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

	// 0x40 to 0x7F
	for i := 0x40; i < 0x80; i++ {
		if i == 0x76 {
			continue
		}
		opStore[i] = LD
	}

	// 0x80 to 0xBF
	for i := 0x80; i < 0xC0; i++ {
		mask := i & 0x0F
		if mask < 0x08 {
			// ADD
			opStore[i|0x80] = ADD
			opStore[i|0x90] = SUB
			opStore[i|0xA0] = AND
			opStore[i|0xB0] = OR
		}
		if mask > 0x7 {
			// ADC
			opStore[i|0x80] = ADC
			opStore[i|0x90] = SBC
			opStore[i|0xA0] = XOR
			opStore[i|0xB0] = CP
		}
	}
	opStore[0xC0] = RET

	for i := 0xC1; i < 0x101; i += 0x10 {
		opStore[i] = POP
	}

	opStore[0xC2] = JP
	opStore[0xC3] = JP
	opStore[0xC4] = CALL

	// 0xC5-0xF5
	for i := 0xC5; i < 0x101; i += 0x10 {
		opStore[i] = PUSH
	}
	opStore[0xC6] = ADD
	// TODO: Implement RST.go
	// opStore[0xC7] = RST
	opStore[0xC8] = RET
	opStore[0xC9] = RET
	opStore[0xCA] = JP
	opStore[0xCB] = CB

	opStore[0xCC] = CALL
	opStore[0xCD] = CALL
	opStore[0xCE] = ADC
	// opStore[0xCF] = RST

	opStore[0xD0] = RET
	opStore[0xD2] = JP
	opStore[0xD4] = CALL
	opStore[0xD5] = PUSH
	opStore[0xD6] = SUB
	// opStore[0xD7] = RST
	opStore[0xD8] = RET
	// opStore[0xD9] = RETI
	opStore[0xDA] = JP
	opStore[0xDC] = CALL
	opStore[0xDE] = SBC
	// opStore[0xDF] = RST

	opStore[0xE0] = LD
	opStore[0xE2] = LD
	opStore[0xE5] = PUSH
	opStore[0xE6] = AND
	// opStore[0xE7] = RST
	opStore[0xE8] = ADD
	opStore[0xE9] = JP
	opStore[0xEA] = LD
	opStore[0xEE] = XOR
	// opStore[0xEF] = RST

	opStore[0xF0] = LD
	opStore[0xF2] = LD
	// TODO: Find out what the IME flag is and how to use it.
	// opStore[0xF3] = DI
	opStore[0xF5] = PUSH
	opStore[0xF6] = OR
	// opStore[0xF7] = RST
	opStore[0xF8] = LD
	opStore[0xF9] = LD
	opStore[0xFA] = LD
	// opStore[0xFB] = EI
	opStore[0xFE] = CP
	// opStore[0xFF] = RST

	return opStore
}
