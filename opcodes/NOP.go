package opcodes

type nop struct {
	label string
	op    byte
}

func NewNOP(lbl string, op byte) *nop {
	return &nop{lbl, op}
}

func (n *nop) Exec(opcode byte) {}

// Label returns the string label of the opcode
func (n *nop) Label() string {
	return "NOP"
}
