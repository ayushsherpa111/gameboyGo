package opcodes

import (
	"github.com/ayushsherpa111/gameboyEMU/cpu"
)

type ld struct {
	label string
	op    byte
	c     *cpu.CPU
}

// Load into a 16 bit register
func (i *ld) r16_u16(r1, r2 uint8, val uint16) {
	i.c.SetRegister(r1, cpu.Register(val&0x00FF))
	i.c.SetRegister(r2, cpu.Register((val&0xFF00)>>8))
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
	i.c.SetRegister(reg, cpu.Register(val))
}

// Load into memory from Stack pointer
func (i *ld) u16_SP(mem uint16) {
	i.c.SetMem(mem, uint8(i.c.SP))
	i.c.SetMem(mem+1, uint8(i.c.SP>>8))
}

// Load from memory into register
func (i *ld) r8_u16(reg uint8, addr uint16) {
	i.c.SetRegister(reg, cpu.Register(i.c.GetMem(addr)))
}

// Load uint16 into Stack pointer
func (i *ld) SP_u16(val uint16) {
	i.c.SP = val
}

func (i *ld) r8_r8(to, from uint8) {
	i.c.SetRegister(to, i.c.GetRegister(from))
}

func (i *ld) Exec() {
	switch i.op {
	case 0x01:
		// LD BC, u16
		i.r16_u16(cpu.B, cpu.C, i.c.Fetch16())
	case 0x02:
		// LD (BC), A
		i.u16_u8(i.c.BC(), uint8(i.c.GetRegister(cpu.A)))
	case 0x06:
		// LD B, u8
		i.r8_u8(cpu.B, i.c.Fetch())
	case 0x08:
		// LD (u16), SP
		i.u16_SP(i.c.Fetch16())
	case 0x0A:
		// LD A, (BC)
		i.r8_u16(cpu.A, i.c.BC())
	case 0x0E:
		// LD C, u8
		i.r8_u8(cpu.C, i.c.Fetch())
	case 0x11:
		// LD DE, u16
		i.r16_u16(cpu.D, cpu.E, i.c.Fetch16())
	case 0x12:
		// LD DE, A
		i.u16_u8(i.c.DE(), uint8(i.c.GetRegister(cpu.A)))
	case 0x16:
		// LD D, u8
		i.r8_u8(cpu.D, i.c.Fetch())
	case 0x1A:
		i.r8_u16(cpu.A, i.c.DE())
	case 0x1E:
		// LD E, u8
		i.r8_u8(cpu.E, i.c.Fetch())
	case 0x21:
		// LD HL, u16
		i.r16_u16(cpu.H, cpu.L, i.c.Fetch16())
	case 0x22:
		mem := i.c.HL() + 1
		i.r16_u8(mem, uint8(i.c.GetRegister(cpu.A)))
	case 0x26:
		// LD H, u8
		i.r8_u8(cpu.H, i.c.Fetch())
	case 0x2A:
		// LD A, (HL+)
		HL := i.c.HL() + 1
		i.r8_u8(cpu.A, i.c.GetMem(HL))
	case 0x2E:
		// LD L, u8
		i.r8_u8(cpu.L, i.c.Fetch())
	case 0x31:
		// LD SP, u16
		i.SP_u16(i.c.Fetch16())
	case 0x32:
		// LD (HL-), A
		mem := i.c.HL() - 1
		i.u16_u8(mem, uint8(i.c.GetRegister(cpu.A)))
	case 0x36:
		// LD (HL), u8
		i.u16_u8(i.c.HL(), i.c.Fetch())
	case 0x3A:
		// LD A, (HL-)
		mem := i.c.HL() - 1
		i.r8_u16(cpu.A, mem)
	case 0x3E:
		// LD A, u8
		i.r8_u8(cpu.A, i.c.Fetch())
	case 0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47:
		// LD B,[B, C, D, F, H, L, (HL), A]
		if i.op != 0x46 {
			// TODO map registers to numbers in the above order
			// i.r8_r8(cpu.B, )
		} else {
			// LD B, (HL)
			i.r8_u16(cpu.B, i.c.HL())
		}
	default:
		panic("invalid opcode")
	}
}

func (i *ld) Label() string {
	return i.label
}

func NewLD(label string, op byte, cpu *cpu.CPU) *ld {
	return &ld{
		label,
		op,
		cpu,
	}
}
