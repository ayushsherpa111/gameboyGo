package cpu

import (
	"fmt"

	"github.com/ayushsherpa111/gameboyEMU/memory"
)

type CPU struct {
	Registers [8]register
	PC        uint16
	memory    memory.Memory
}

var opcodeCollection [0x100]memory.Opcode = [0x100]memory.Opcode{}

func NewCPU() *CPU {
	return &CPU{
		Registers: initRegisters(),
		PC:        0x0,
	}
}

func (c *CPU) combine(reg1, reg2 int) uint16 {
	return uint16(c.Registers[reg1])<<8 | uint16(c.Registers[reg2])
}

func (c *CPU) AF() uint16 {
	return c.combine(A, F)
}

func (c *CPU) BC() uint16 {
	return c.combine(B, C)
}

func (c *CPU) DE() uint16 {
	return c.combine(D, E)
}

func (c *CPU) HL() uint16 {
	return c.combine(H, L)
}

func (c *CPU) LD(opcode byte) *memory.Opcode {
	op := memory.NewOpcode(opcode)
	switch opcode {
	case 0x01:
		op.Label = "LD BC, u16"
		op.Length = 3
		op.Steps = append(op.Steps, func() {
			c.Registers[C] = register(c.memory.GetByte(c.PC))
			c.PC++
		})
		op.Steps = append(op.Steps, func() {
			c.Registers[B] = register(c.memory.GetByte(c.PC))
			c.PC++
		})
	case 0x02:
		op.Label = "LD (BC), A"
		op.Length = 1
		op.Steps = append(op.Steps, func() {
			c.memory.SetByte(c.BC(), byte(c.Registers[A]))
		})
	case 0x06:
		op.Label = "LD B, u8"
		op.Length = 2
		op.Steps = append(op.Steps, func() {

		})
	}
	return op
}

func (c *CPU) Step() {
	// FETCH instruction
	inst := c.memory.GetByte(c.PC) // 4 clock cycles
	c.PC++
	opcodeCollection[inst].Execute()
	// Decode the opcode
	fmt.Println(inst)
}

func (c *CPU) ZeroFlag() uint8 {
	return uint8(c.Registers[F]) & ZERO
}
