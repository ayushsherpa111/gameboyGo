package opcodes

import "github.com/ayushsherpa111/gameboyEMU/cpu"

const (
	vec00 uint16 = 0x00
	vec08        = 0x08
	vec10        = 0x10
	vec18        = 0x18
	vec20        = 0x20
	vec28        = 0x28
	vec30        = 0x30
	vec38        = 0x38
)

type rst struct {
	c      *cpu.CPU
	regMap map[uint8]map[uint8]uint16
}

func (r *rst) _RST(vec uint16) {
	r.c.PushSP(r.c.PC)
	r.c.PC = vec
}

func (r *rst) Exec(op byte) {
	var vec uint16 = r.regMap[op&0x0F][op&0xF0]
	r._RST(vec)
}

func NewRST(c *cpu.CPU) *rst {
	return &rst{
		c: c,
		regMap: map[uint8]map[uint8]uint16{
			0x07: {
				0xC0: vec00,
				0xD0: vec10,
				0xE0: vec20,
				0xF0: vec30,
			},
			0x0F: {
				0xC0: vec08,
				0xD0: vec18,
				0xE0: vec28,
				0xF0: vec38,
			},
		},
	}
}
