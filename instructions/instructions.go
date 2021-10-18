package instructions

type Instruction interface {
	Exec(byte)
}
