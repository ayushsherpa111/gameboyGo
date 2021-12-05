package interfaces

type Instruction interface {
	Exec(byte)
}
