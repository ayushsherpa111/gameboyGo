package cpu

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/ayushsherpa111/gameboyEMU/interfaces"
	"github.com/ayushsherpa111/gameboyEMU/memory"
	"github.com/ayushsherpa111/gameboyEMU/types"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	V_BLANK uint8 = 1 << iota
	LCD_STAT
	TIMER
	SERIAL
	JOYPAD
)

var (
	VB_VEC uint16 = 0x40
	ST_VEC uint16 = 0x48
	TM_VEC uint16 = 0x50
	SE_VEC uint16 = 0x58
	JP_VEC uint16 = 0x60
)

var BOOTLOADER_UNLOADED bool = false

type CPU struct {
	registers  [8]uint8
	PC         uint16
	SP         uint16
	memory     interfaces.Mem
	store      []interfaces.Instruction
	ime        bool
	CloseChan  chan struct{}
	CycleCount uint64
	Scheduler  interfaces.Scheduler
	Halted     bool
	joypadChan <-chan types.KeyboardEvent
	logFile    *os.File
	logChan    chan string
	isDebug    bool
}

func NewCPU(mem interfaces.Mem, joyPad <-chan types.KeyboardEvent, debug bool) *CPU {
	logFile, err := os.OpenFile("vram.test", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		log.Fatalf("%s\n", err.Error())
	}
	return &CPU{
		registers:  [8]uint8{},
		PC:         0x000,
		SP:         0xFFFE,
		memory:     mem,
		ime:        false,
		CloseChan:  make(chan struct{}),
		CycleCount: 0,
		Halted:     false,
		joypadChan: joyPad,
		logChan:    make(chan string, 100),
		logFile:    logFile,
	}
}
func (c *CPU) ListenForKeyPress() {
	for {
		inp := <-c.joypadChan
		if inp.Key == sdl.K_q {
			c.logChan <- ""
			break
		}
		// fmt.Println("Keyboard pressed")
		c.memory.HandleInput(inp)
	}
}

func (c *CPU) WriteToFile() {
	for {
		log := <-c.logChan
		if len(log) == 0 {
			break
		}
		fmt.Fprint(c.logFile, log)
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

func (c *CPU) ScheduleEI(cycles uint64) {
	c.Scheduler.ScheduleEvent(c.SetIME(true), cycles, types.EV_EI)
}

func (c *CPU) tick() {
	c.CycleCount += 4

	// INFO: Tick PPU 4 times. [1 T cycle]
	c.memory.TickAllComponents(c.CycleCount)
}

func (c *CPU) Fetch() (uint8, error) {
	c.tick()
	b := c.memory.MemRead(c.PC, c.CycleCount)

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
	c.tick()
	if e := c.memory.MemWrite(addr, val, c.CycleCount); e != nil {
		log.Fatalf("Error setting byte in memory.")
	}
}

func (c *CPU) GetMem(addr uint16) *uint8 {
	c.tick()
	return c.memory.MemRead(addr, c.CycleCount)
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
	if !c.Halted {
		inst, err := c.Fetch()
		if err != nil {
			return err
		}

		// arg1 := c.memory.MemRead(c.PC, c.CycleCount)
		// arg2 := c.memory.MemRead(c.PC+1, c.CycleCount)
		// c.logChan <- fmt.Sprintf("SP:0x%x PC: 0x%02x OP:0x%02x ARG1:0x%02x ARG2:0x%02x  Registers: %s Flag: 0b%04b\n",
		// 	c.SP, c.PC-1, inst, *arg1, *arg2, hexVals(c.registers[:]), c.registers[F]>>4)

		// if c.PC >= 0xFF00 {
		// 	c.logChan <- ""
		// 	log.Fatalln("AT MEMORY 0xFF00")
		// }

		store[inst].Exec(inst)
	} else {
		// c.CycleCount++
		c.tick()
		fmt.Printf("HALTED: %d\n", c.CycleCount)
	}

	c.Scheduler.Tick()

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
	IE, IF := c.memory.MemRead(memory.INTERRUPT_ENABLE, c.CycleCount), c.memory.MemRead(memory.INTERRUPT_FLAG, c.CycleCount)
	interrupt := *IE & *IF

	if interrupt > 0 {
		c.Halted = false
	}

	if !c.ime {
		return
	}

	if interrupt&0x0F >= 1 {
		c.ime = false
		c.PushSP(c.PC)
	}

	switch {
	case isEnabled(interrupt, V_BLANK):
		*IF = clearBit(*IF, V_BLANK)
		c.PC = VB_VEC
		// read for Joypad input as well
	case isEnabled(interrupt, LCD_STAT):
		*IF = clearBit(*IF, LCD_STAT)
		c.PC = ST_VEC

	case isEnabled(interrupt, TIMER):
		*IF = clearBit(*IF, TIMER)
		c.PC = TM_VEC

	case isEnabled(interrupt, SERIAL):
		*IF = clearBit(*IF, SERIAL)
		c.PC = SE_VEC

	case isEnabled(interrupt, JOYPAD):
		*IF = clearBit(*IF, JOYPAD)
		c.PC = JP_VEC
	}

}

func clearBit(flag, bitNum uint8) uint8 {
	return flag &^ bitNum
}
