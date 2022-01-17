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
var bootLoader []byte

func main() {
	ROM := os.Args[1]
	mem, err := memory.InitMem(bootLoader, ROM)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	cpu := cpu.NewCPU(mem)
	go cpu.ListenIMEChan()
	store := opcodes.NewOpcodeStore(cpu) // LUT for decoding instructions

	for {
		if e := cpu.FetchDecodeExec(store); e != nil {
			cpu.CloseChan <- struct{}{}
			return
		}
	}
}
