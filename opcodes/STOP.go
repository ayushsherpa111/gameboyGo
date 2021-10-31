package opcodes

import "os"

type stop struct {
}

func (s *stop) Exec(op byte) {
	os.Exit(-1)
}

func NewStop() *stop {
	return &stop{}
}
