package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ayushsherpa111/gameboyEMU/cpu"
	"github.com/ayushsherpa111/gameboyEMU/frontend"
	_ "github.com/ayushsherpa111/gameboyEMU/frontend"
	"github.com/ayushsherpa111/gameboyEMU/joypad"
	"github.com/ayushsherpa111/gameboyEMU/logger"
	"github.com/ayushsherpa111/gameboyEMU/memory"
	"github.com/ayushsherpa111/gameboyEMU/opcodes"
	"github.com/ayushsherpa111/gameboyEMU/ppu"
	"github.com/ayushsherpa111/gameboyEMU/scheduler"
	"github.com/veandco/go-sdl2/sdl"
)

//go:embed ROMs/BOOTLOADER.gb
var bootLoader []byte

func main() {
	if e := sdl.Init(frontend.CONFIG); e != nil {
		// EmuWindow.lgr.Fatalf("Failed to Initialize SDL")
		log.Fatal("Failed to Initialize SDL")
	}
	var ROM string
	var debug bool

	lgr := logger.NewLogger(os.Stdout, debug, "Main")

	flag.StringVar(&ROM, "r", "", "ROM file to execute.")
	flag.BoolVar(&debug, "d", false, "Debug flag.")

	flag.Parse()

	if ROM == "" {
		lgr.Fatalf("No ROM provided")
		os.Exit(2)
	}
	bufferChan := make(chan []uint32, 10)
	inputChan := make(chan sdl.Keycode, 120)
	joyPadCtx := joypad.NewContext()

	ppu := ppu.NewPPU(bufferChan)
	mem, err := memory.InitMem(bootLoader, ROM, debug, ppu, joyPadCtx)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	cpu := cpu.NewCPU(mem)
	sched := scheduler.NewScheduler(cpu)
	cpu.Scheduler = sched
	mem.SetScheduler(sched)

	frontend.SetupWindow()
	frontend.EmuWindow.SetChannels(bufferChan, inputChan)

	store := opcodes.NewOpcodeStore(cpu) // LUT for decoding instructions

	go func() {
		for {
			if e := cpu.FetchDecodeExec(store); e != nil {
				cpu.CloseChan <- struct{}{}
				return
			}
			select {
			case k := <-inputChan:
				switch k {
				case sdl.K_q:
					return
				}
			default:
			}
		}
	}()

	frontend.EmuWindow.Run()
}
