package memory

type step func()

type Opcode struct {
	Op     byte
	Label  string
	Length uint8
	Steps  []step
}

func NewOpcode(op byte) *Opcode {
	return &Opcode{
		Op:     op,
		Label:  "",
		Length: 1,
		Steps:  make([]step, 0, 4),
	}
}

func (o *Opcode) Execute() {
	for _, step := range o.Steps {
		step()
	}
}
