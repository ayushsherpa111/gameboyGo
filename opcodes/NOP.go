package opcodes

type nop struct {
}

func NewNOP() *nop {
	return &nop{}
}

func (n *nop) Exec(opcode byte) {}
