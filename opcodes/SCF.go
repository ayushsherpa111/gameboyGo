package opcodes

import "github.com/ayushsherpa111/gameboyEMU/cpu"

type scf struct {
	c *cpu.CPU
}

func (s *scf) Exec(op byte) {
	s.c.SET_NEG(false)
	s.c.SET_HALF_CARRY(false)
	s.c.SET_CARRY(true)
}

func NewSCF(c *cpu.CPU) *scf {
	return &scf{c}
}
