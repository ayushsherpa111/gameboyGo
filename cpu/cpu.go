package cpu

import (
	"errors"
	"fmt"
	"log"

	"github.com/ayushsherpa111/gameboyEMU/interfaces"
	"github.com/ayushsherpa111/gameboyEMU/memory"
)

const (
	V_BLANK uint8 = 1 << iota
	LCD_STAT
	TIMER
	SERIAL
	JOYPAD
)

var (
	VB_VEC = 0x40
	ST_VEC = 0x48
	TM_VEC = 0x50
	SE_VEC = 0x58
	JP_VEC = 0x60
)

type CPU struct {
	registers  [8]uint8
	PC         uint16
	SP         uint16
	memory     memory.Mem
	store      []interfaces.Instruction
	ime        bool
	CloseChan  chan struct{}
	CycleCount uint64
	Scheduler  interfaces.Scheduler
}

func NewCPU(mem memory.Mem) *CPU {
	return &CPU{
		registers:  [8]uint8{},
		PC:         0x000,
		SP:         0xFFFE,
		memory:     mem,
		ime:        false,
		CloseChan:  make(chan struct{}),
		CycleCount: 0,
	}
}

func (c *CPU) SetIME(v bool) func() {
	return func() {
		c.ime = v
	}
}

func (c *CPU) SET_CARRY(set bool) {
	if set {
		c.registers[F] |= CARRY
	} else {
		c.registers[F] &= ^CARRY
	}
	// fmt.Printf("CARRY flag Changed. Final Value: 0b%08b\n", c.registers[F]&CARRY)
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
}

func (c *CPU) SET_NEG(set bool) {
	if set {
		c.registers[F] |= NEG
	} else {
		c.registers[F] &= ^NEG
	}
}

func (c *CPU) SetRegister(reg uint8, val uint8) {
	c.registers[reg] = val
}

func (c *CPU) GetRegister(reg uint8) *uint8 {
	// fmt.Println("Registers ", c.registers)
	return &c.registers[reg]
}

func (c *CPU) tick() {
	c.Scheduler.Tick()
	c.CycleCount++

	// INFO: Tick PPU 4 times. [1 T cycle]
	c.memory.TickAllComponents(c.CycleCount)
}

func (c *CPU) Fetch() (uint8, error) {
	if c.PC >= 0x100 {
		c.memory.UnloadBootloader()
	}
	b := c.memory.MemRead(c.PC)
	c.tick()

	if b == nil {
		return 0, errors.New("PC is pointing to an invalid address")
	}

	c.PC += 1
	return *b, nil
}

func (c *CPU) Fetch16() uint16 {
	HB, _ := c.Fetch()
	LB, _ := c.Fetch()
	return uint16(HB) | uint16(LB)<<8
}

func (c *CPU) SetMem(addr uint16, val uint8) {
	if e := c.memory.MemWrite(addr, val); e != nil {
		log.Fatalf("Error setting byte in memory.")
	}
	c.tick()
}

func (c *CPU) GetMem(addr uint16) *uint8 {
	return c.memory.MemRead(addr)
}

func (c *CPU) combine(reg1, reg2 int) uint16 {
	return uint16(c.registers[reg1])<<8 | uint16(c.registers[reg2])
}

func (c *CPU) setMulReg(reg1, reg2 int, val uint16) {
	c.registers[reg1] = uint8(val >> 8)
	c.registers[reg2] = uint8(val)
}

func (c *CPU) AF() uint16 {
	return c.combine(A, F)
}

func (c *CPU) SetAF(val uint16) {
	c.setMulReg(A, F, val)
}

func (c *CPU) BC() uint16 {
	return c.combine(B, C)
}

func (c *CPU) SetBC(val uint16) {
	c.setMulReg(B, C, val)
}

func (c *CPU) DE() uint16 {
	return c.combine(D, E)
}

func (c *CPU) SetDE(val uint16) {
	c.setMulReg(D, E, val)
}

func (c *CPU) HL() uint16 {
	return c.combine(H, L)
}

func (c *CPU) SetHL(val uint16) {
	c.setMulReg(H, L, val)
}

func hexVals(regs []uint8) string {
	var reg_bin string
	for i, v := range regs {
		reg_bin += fmt.Sprintf("%s: 0x%04x ", GetRegName(uint8(i)), v)
	}
	return reg_bin
}

func (c *CPU) FetchDecodeExec(store [0x100]interfaces.Instruction) error {
	// FETCH instruction
	inst, err := c.Fetch()
	if err != nil {
		return err
	}

	// fmt.Printf("A: %02X F: %02X B: %02X C: %02X D: %02X E: %02X H: %02X L: %02X SP: %04X PC: %04X \n",
	// 	c.registers[A], c.registers[F], c.registers[B], c.registers[C], c.registers[D], c.registers[E], c.registers[H], c.registers[L], c.SP, c.PC-1)

	store[inst].Exec(inst)

	// handle interrupt at the end of each cycle
	c.handleInterrupt()

	return nil
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
	if v := c.CarryFlag(); v {
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
	c.SP++
	return u16
}

func isEnabled(inter uint8, bit uint8) bool {
	return (inter & bit) == bit
}

func (c *CPU) handleInterrupt() {
	if !c.ime {
		return
	}
	IE, IF := c.memory.MemRead(memory.INTERRUPT_ENABLE), c.memory.MemRead(memory.INTERRUPT_FLAG)
	interrupt := *IE & *IF

	if interrupt&0x0F >= 1 {
		c.ime = false
		c.PushSP(c.SP)
	}

	switch {
	case isEnabled(interrupt, V_BLANK):
		*IF = clearBit(*IF, V_BLANK)
		c.PC = uint16(VB_VEC)
	case isEnabled(interrupt, LCD_STAT):
		*IF = clearBit(*IF, LCD_STAT)
		c.PC = uint16(VB_VEC)

	case isEnabled(interrupt, TIMER):
		*IF = clearBit(*IF, TIMER)
		c.PC = uint16(VB_VEC)

	case isEnabled(interrupt, SERIAL):
		*IF = clearBit(*IF, SERIAL)
		c.PC = uint16(VB_VEC)

	case isEnabled(interrupt, JOYPAD):
		*IF = clearBit(*IF, JOYPAD)
		c.PC = uint16(VB_VEC)
	}

}

func clearBit(flag, bitNum uint8) uint8 {
	return flag &^ bitNum
}
