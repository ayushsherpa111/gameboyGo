package opcodes

import "github.com/ayushsherpa111/gameboyEMU/cpu"

type rst struct {
}

func NewRST(c *cpu.CPU) *rst {
	return &rst{}
}
