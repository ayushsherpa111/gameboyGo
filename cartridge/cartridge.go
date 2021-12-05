package cartridge

import "os"

type cart struct {
	file *os.File
}

func (a *cart) HeaderInfo() {
}

func (a *cart) LoadROM(romPath string) {
}

func (a *cart) SetInfo() {
}

func NewCart() *cart {
	return &cart{}
}
