package opcodes

import (
	"github.com/ayushsherpa111/gameboyEMU/cpu"
)

type ld struct {
	c *cpu.CPU
}

// Load into a 16 bit register
func (i *ld) r16_u16(r1, r2 uint8, val uint16) {
	i.c.SetRegister(r1, uint8(val>>8))
	i.c.SetRegister(r2, uint8(val))
}

// Load into memory
func (i *ld) u16_u8(addr uint16, v uint8) {
	i.c.SetMem(addr, v)
}

func (i *ld) r16_u8(addr uint16, val uint8) {
	i.c.SetMem(addr, val)
}

// Load into 8bit register
func (i *ld) r8_u8(reg uint8, val uint8) {
	i.c.SetRegister(reg, uint8(val))
}

// Load into memory from Stack pointer
func (i *ld) u16_SP(mem uint16) {
	i.c.SetMem(mem, uint8(i.c.SP))
	i.c.SetMem(mem+1, uint8(i.c.SP>>8))
}

// Load from memory into register
func (i *ld) r8_u16(reg uint8, addr uint16) {
	i.c.SetRegister(reg, *i.c.GetMem(addr))
}

// Load uint16 into Stack pointer
func (i *ld) SP_u16(val uint16) {
	i.c.SP = val
}

func (i *ld) r8_r8(to, from uint8) {
	i.c.SetRegister(to, *i.c.GetRegister(from))
}

func (i *ld) r16_sp_u8(highReg, lowReg uint8, val int8) {
	i.c.SET_ZERO(false)
	i.c.SET_NEG(false)
}

func ld_tar(opcode uint8) uint8 {
	var tar uint8
	lower := opcode & 0x0F
	upper := opcode & 0xF0
	switch {
	case upper == 0x40:
		if lower >= 0x00 && lower <= 0x07 {
			tar = cpu.B
		} else if lower >= 0x08 && lower <= 0x0F {
			tar = cpu.C
		}
	case upper == 0x50:
		if lower >= 0x00 && lower <= 0x07 {
			tar = cpu.D
		} else if lower >= 0x08 && lower <= 0x0F {
			tar = cpu.E
		}
	case upper == 0x60:
		if lower >= 0x00 && lower <= 0x07 {
			tar = cpu.H
		} else if lower >= 0x08 && lower <= 0x0F {
			tar = cpu.L
		}
	case upper == 0x70:
		tar = cpu.A
	}
	return tar
}

func ld_src(opcode uint8) uint8 {
	var src uint8
	lower := opcode & 0x0F
	switch lower % 8 {
	case 0x00:
		src = cpu.B
	case 0x01:
		src = cpu.C
	case 0x02:
		src = cpu.D
	case 0x03:
		src = cpu.E
	case 0x04:
		src = cpu.H
	case 0x05:
		src = cpu.L
	case 0x07:
		src = cpu.A
	}
	return src
}

func (i *ld) LD_HL_SPi8() {
	u8, _ := i.c.Fetch()
	v1 := int32(int8(u8))
	v2 := int32(i.c.SP)
	sum := int32(v2 + v1)

	i.c.SET_ZERO(false)
	i.c.SET_NEG(false)

	// i.c.SET_HALF_CARRY((uint8(i.c.SP)&0x0F)+(u8&0x0F) > 0x0F)
	// i.c.SET_CARRY(int16(i.c.SP)&0xFF+int16(u8)&0xFF > 0xFF)

	tVal := int32(i.c.SP) ^ v1 ^ sum

	i.c.SET_HALF_CARRY((tVal & 0x10) == 0x10)
	i.c.SET_CARRY((tVal & 0x100) == 0x100)

	i.c.SetRegister(cpu.H, uint8(uint16(sum)>>8))
	i.c.SetRegister(cpu.L, uint8(uint16(sum)))
}

func (i *ld) Exec(opcode byte) {
	A := i.c.GetRegister(cpu.A)
	C := i.c.GetRegister(cpu.C)
	HL := i.c.HL()
	switch opcode {
	case 0x01:
		// LD BC, u16
		i.r16_u16(cpu.B, cpu.C, i.c.Fetch16())
	case 0x02:
		// LD (BC), A
		i.u16_u8(i.c.BC(), *A)
	case 0x06:
		// LD B, u8
		arg, _ := i.c.Fetch()
		i.r8_u8(cpu.B, arg)
	case 0x08:
		// LD (u16), SP
		i.u16_SP(i.c.Fetch16())
	case 0x0A:
		// LD A, (BC)
		i.r8_u16(cpu.A, i.c.BC())
	case 0x0E:
		// LD C, u8
		arg, _ := i.c.Fetch()
		i.r8_u8(cpu.C, arg)
	case 0x11:
		// LD DE, u16
		i.r16_u16(cpu.D, cpu.E, i.c.Fetch16())
	case 0x12:
		// LD DE, A
		i.u16_u8(i.c.DE(), *A)
	case 0x16:
		// LD D, u8
		arg, _ := i.c.Fetch()
		i.r8_u8(cpu.D, arg)
	case 0x1A:
		i.r8_u16(cpu.A, i.c.DE())
	case 0x1E:
		// LD E, u8
		arg, _ := i.c.Fetch()
		i.r8_u8(cpu.E, arg)
	case 0x21:
		// LD HL, u16
		i.r16_u16(cpu.H, cpu.L, i.c.Fetch16())
	case 0x22:
		// LD HL+, A
		i.r16_u8(HL, *A)
		i.c.SetHL(HL + 1)
	case 0x26:
		// LD H, u8
		arg, _ := i.c.Fetch()
		i.r8_u8(cpu.H, arg)
	case 0x2A:
		// LD A, (HL+)
		i.r8_u8(cpu.A, *i.c.GetMem(HL))
		i.c.SetHL(HL + 1)
	case 0x2E:
		// LD L, u8
		arg, _ := i.c.Fetch()
		i.r8_u8(cpu.L, arg)
	case 0x31:
		// LD SP, u16
		i.SP_u16(i.c.Fetch16())
	case 0x32:
		// LD (HL-), A
		i.u16_u8(HL, *A)
		i.c.SetHL(HL - 1)
	case 0x36:
		// LD (HL), u8
		arg, _ := i.c.Fetch()
		i.u16_u8(i.c.HL(), arg)
	case 0x3A:
		// LD A, (HL-)
		i.r8_u16(cpu.A, HL)
		i.c.SetHL(HL - 1)
	case 0x3E:
		// LD A, u8
		arg, _ := i.c.Fetch()
		i.r8_u8(cpu.A, arg)
	case 0xE0:
		// LD (0xFF00+u8), A
		arg, _ := i.c.Fetch()
		i.u16_u8(0xFF00+uint16(arg), *A)
	case 0xE2:
		// LD (0xFF00+C), A
		i.u16_u8(0xFF00+uint16(*C), *A)
	case 0xEA:
		// LD (u16), A
		i.u16_u8(i.c.Fetch16(), *A)
	case 0xF0:
		// LD A, (0xFF00+u8)
		arg, _ := i.c.Fetch()
		val := 0xFF00 + uint16(arg) // 44
		i.r8_u16(cpu.A, val)
		return
	case 0xF2:
		// LD A, (0xFF00+C)
		i.r8_u16(cpu.A, (0xFF00 + uint16(*C)))
	case 0xF8:
		// LD HL, SP+i8
		i.LD_HL_SPi8()
	case 0xF9:
		i.SP_u16(HL)
	case 0xFA:
		// LD A, (u16)
		i.r8_u16(cpu.A, i.c.Fetch16())
	}

	switch {
	case opcode >= 0x40 && opcode <= 0x7f:
		// 0x40 - 0x47 = LD B
		// 0x48 - 0x4F = LD C
		// 0x50 - 0x57 = LD D
		// 0x58 - 0x5F = LD E
		// 0x60 - 0x57 = LD H
		// 0x68 - 0x5F = LD L
		// 0x70 - 0x77 = LD (HL)
		// 0x78 - 0x7F = LD A

		// Not LD (HL)
		if !(opcode >= 0x70 && opcode <= 0x77) {
			r1 := ld_tar(opcode)
			if (opcode&0x0F)%8 != 0x06 {
				r2 := ld_src(opcode)
				i.r8_r8(r1, r2)
			} else {
				i.r8_u16(r1, i.c.HL())
			}
		} else {
			// LD (HL)
			reg := i.c.GetRegister(ld_src(opcode))
			i.u16_u8(i.c.HL(), *reg)
		}
	}
}

func NewLD(cpu *cpu.CPU) *ld {
	return &ld{
		c: cpu,
	}
}
