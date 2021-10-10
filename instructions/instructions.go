package instructions


type Instruction interface {
	Exec()
	Label() string
}

