package main

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/ayushsherpa111/gameboyEMU/cpu"
	"github.com/ayushsherpa111/gameboyEMU/memory"
	"github.com/ayushsherpa111/gameboyEMU/opcodes"
)

//go:embed ROMs/BOOTLOADER.gb
var boot_loader []byte

func main() {
	ROM := "./ROMs/Tetris.gb"
	mem, err := memory.InitMem(boot_loader, ROM)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	cpu := cpu.NewCPU(mem)
	store := opcodes.NewOpcodeStore(cpu) // LUT for decoding instructions

	for {
		cpu.Decode(store)
		// fmt.Scanln()
	}
}
