package main

import (
	"io/ioutil"
	"log"

	"github.com/ayushsherpa111/gameboyEMU/cpu"
	"github.com/ayushsherpa111/gameboyEMU/opcodes"
)

func main() {
	// TODO
	_, err := ioutil.ReadFile("./ROMs/Pokemon - Blue Version (USA, Europe) (SGB Enhanced).gb")
	if err != nil {
		log.Fatalf("Invalid ROM %s\n", err.Error())
	}
	cpu := cpu.NewCPU()
	store := opcodes.NewOpcodeStore(cpu) // LUT for decoding instructions
	for {
		cpu.Decode(store)
	}
}
