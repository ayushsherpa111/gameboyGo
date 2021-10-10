package cpu

import (
	"fmt"

	"github.com/ayushsherpa111/gameboyEMU/instructions"
	"github.com/ayushsherpa111/gameboyEMU/memory"
)

type CPU struct {
	registers [8]Register
	pc        uint16
	SP        uint16
	memory    *memory.Memory
}

func NewCPU() *CPU {
	return &CPU{
		registers: initRegisters(),
		pc:        0x0,
		memory:    memory.InitMem(),
	}
}

func (c *CPU) SetRegister(reg uint8, val Register) {
	if reg < A || reg > L {
		return
	}
	c.registers[reg] = val
}

func (c *CPU) GetRegister(reg uint8) Register {
	return c.registers[reg]
}

func (c *CPU) Fetch() uint8 {
	b := c.memory.GetByte(c.pc)
	c.pc++
	return b
}

func (c *CPU) Fetch16() uint16 {
	return uint16(c.Fetch()) | uint16(c.Fetch())<<8
}

func (c *CPU) SetMem(addr uint16, val uint8) {
	c.memory.SetByte(addr, val)
}

func (c *CPU) combine(reg1, reg2 int) uint16 {
	return uint16(c.registers[reg1])<<8 | uint16(c.registers[reg2])
}

func (c *CPU) af() uint16 {
	return c.combine(A, F)
}

func (c *CPU) BC() uint16 {
	return c.combine(B, C)
}

func (c *CPU) de() uint16 {
	return c.combine(D, E)
}

func (c *CPU) hl() uint16 {
	return c.combine(H, L)
}

func (c *CPU) Decode(store []instructions.Instruction) {
	// FETCH instruction
	inst := c.Fetch()
	fmt.Println(inst)
	store[inst].Exec()
}

func (c *CPU) ZeroFlag() uint8 {
	return uint8(c.registers[F]) & ZERO
}
