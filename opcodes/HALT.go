package opcodes

type halt struct {
}

func (h *halt) Exec(opcode byte) {

}

func NewHalt() *halt {
	return &halt{}
}
