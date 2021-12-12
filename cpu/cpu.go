package cpu

import (
	"fmt"
	"log"

	instructions "github.com/ayushsherpa111/gameboyEMU/interfaces"
	"github.com/ayushsherpa111/gameboyEMU/memory"
)

type CPU struct {
	registers [8]uint8
	PC        uint16
	SP        uint16
	memory    memory.Mem
	store     []instructions.Instruction
	ime       bool
}

func NewCPU(mem memory.Mem) *CPU {
	return &CPU{
		registers: initRegisters(),
		PC:        0x000,
		SP:        0xFFFE,
		memory:    mem,
		ime:       false,
	}
}

func (c *CPU) SET_CARRY(set bool) {
	if set {
		c.registers[F] |= CARRY
	} else {
		c.registers[F] &= ^CARRY
	}
	// fmt.Printf("CARRY Flag Changed. Final Value: 0b%08b\n", c.registers[F]&CARRY)
}

func (c *CPU) SET_ZERO(set bool) {
	if set {
		c.registers[F] |= ZERO
	} else {
		c.registers[F] &= ^ZERO
	}
	// fmt.Printf("Zero Flag Changed. Final Value: 0b%08b\n", c.registers[F]&ZERO)
}

func (c *CPU) SET_HALF_CARRY(set bool) {
	if set {
		c.registers[F] |= HALFCARRY
	} else {
		c.registers[F] &= ^HALFCARRY
	}
	// fmt.Printf("HalfCarry Flag Changed. Final Value: 0b%08b\n", c.registers[F]&HALFCARRY)
}

func (c *CPU) SET_NEG(set bool) {
	if set {
		c.registers[F] |= NEG
	} else {
		c.registers[F] &= ^NEG
	}
	// fmt.Printf("NEG Flag Changed. Final Value: 0b%08b\n", c.registers[F]&NEG)
}

func (c *CPU) SetRegister(reg uint8, val uint8) {
	if reg < A || reg > L {
		return
	}
	c.registers[reg] = val
}

func (c *CPU) GetRegister(reg uint8) *uint8 {
	return &c.registers[reg]
}

func (c *CPU) Fetch() uint8 {
	b := c.memory.MemRead(c.PC)
	c.PC += 1
	return *b
}

func (c *CPU) Fetch16() uint16 {
	return uint16(c.Fetch()) | uint16(c.Fetch())<<8
}

func (c *CPU) SetMem(addr uint16, val uint8) {
	if e := c.memory.MemWrite(addr, val); e != nil {
		log.Fatalf("Error setting byte in memory.")
	}
}

func (c *CPU) GetMem(addr uint16) *uint8 {
	return c.memory.MemRead(addr)
}

func (c *CPU) combine(reg1, reg2 int) uint16 {
	return uint16(c.registers[reg1])<<8 | uint16(c.registers[reg2])
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

func (c *CPU) FetchDecodeExec(store [0xFF]instructions.Instruction) {
	// FETCH instruction
	inst := c.Fetch()
	fmt.Printf("PC: 0x%02x OP:0x%02x\n", c.PC, inst)
	store[inst].Exec(inst)
}

func (c *CPU) ZeroFlag() bool {
	return (c.registers[F] & ZERO) != 0
}

func (c *CPU) CarryFlag() bool {
	return (c.registers[F] & CARRY) != 0
}

func (c *CPU) NegativeFlag() bool {
	return (c.registers[F] & NEG) != 0
}

func (c *CPU) HalfCarryFlag() bool {
	return (c.registers[F] & HALFCARRY) != 0
}

func (c *CPU) CarryVal() uint8 {
	if c.CarryFlag() {
		return 0x01
	}
	return 0x00
}

func (c *CPU) PushSP(val uint16) {
	c.SP -= 1
	c.SetMem(c.SP, uint8(val>>8))
	c.SP -= 1
	c.SetMem(c.SP, uint8(val))
}

func (c *CPU) FetchSP() uint16 {
	var u16 uint16 = 0
	u16 |= uint16(*c.GetMem(c.SP))
	c.SP++
	u16 |= uint16(*c.GetMem(c.SP)) << 8
	return u16
}
