package main

import (
	"fmt"
	"os"

	"github.com/ayushsherpa111/gameboyEMU/cpu"
	"github.com/ayushsherpa111/gameboyEMU/opcodes"
)

func main() {
	ROM := "./ROMs/07-jr,jp,call,ret,rst.gb"
	cpu := cpu.NewCPU()
	if ok, err := cpu.Load_ROM(ROM); !ok {
		fmt.Println((err.Error()))
		os.Exit(-1)
	}
	store := opcodes.NewOpcodeStore(cpu) // LUT for decoding instructions

	for {
		cpu.Decode(store)
		fmt.Scanln()
	}
}
