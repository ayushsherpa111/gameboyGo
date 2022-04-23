package frontend

import (
	"github.com/veandco/go-sdl2/sdl"
)

const (
	SC_WIDTH  = 256 // Width of the window
	SC_HEIGHT = 256 // Height of the window
	SCALE     = 2

	TX_WIDTH  = 160 // Width of the Texture
	TX_HEIGHT = 144 // Height of the Texture

	CONFIG = sdl.INIT_VIDEO | sdl.INIT_EVENTS | sdl.INIT_AUDIO
	TITLE  = "GB EMU"
)

var EmuWindow *window
var DebugWindow *window

func init() {

	if e := sdl.Init(CONFIG); e != nil {
		EmuWindow.lgr.Fatalf("Failed to Initialize SDL")
	}

	if EmuWindow, err := createWindow(SC_WIDTH*SCALE, SC_HEIGHT*SCALE, nil); err != nil {
		EmuWindow.lgr.Fatalf(err.Error())
	} else {
		EmuWindow.lgr.Infof("Initiating Frontend\n")
		EmuWindow.lgr.Printf("Window Initialized with dimensions %d X %d\n", SC_WIDTH, SC_HEIGHT)

		if err := EmuWindow.createTexture(TX_WIDTH, TX_HEIGHT); err != nil {
			EmuWindow.lgr.Fatalf("Failed to Initialize Texture %s", err.Error())
		}
	}

	// go emuWindow.Run()
}

func (w *window) Run() {
	defer sdl.Quit()
	defer EmuWindow.cleanUp()
	isRunning := true

	for isRunning {
		// handle key presses.
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.KeyboardEvent:
				isRunning = false
				w.lgr.Printf("Quitting...\n")
			}
		}
	}
}

func (w *window) SetChannels(bufferChan <-chan []uint32, inputChan chan<- sdl.Event) {
	w.bufferChan = bufferChan
	w.inputChan = inputChan
}
