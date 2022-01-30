package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"

	"github.com/ayushsherpa111/gameboyEMU/cpu"
	"github.com/ayushsherpa111/gameboyEMU/memory"
	"github.com/ayushsherpa111/gameboyEMU/opcodes"
)

//go:embed ROMs/BOOTLOADER.gb
var bootLoader []byte

func main() {
	var ROM string

	flag.StringVar(&ROM, "r", "", "ROM file to execute")
	flag.Parse()

	if ROM == "" {
		os.Exit(2)
	}
	// ROM := os.Args[1]
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
