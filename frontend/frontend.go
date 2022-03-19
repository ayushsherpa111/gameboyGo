package frontend

import (
	"os"

	"github.com/ayushsherpa111/gameboyEMU/logger"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	SC_WIDTH  = 500 // Width of the window
	SC_HEIGHT = 450 // Height of the window

	TX_WIDTH  = 160 // Width of the Texture
	TX_HEIGHT = 144 // Height of the Texture

	CONFIG = sdl.INIT_VIDEO | sdl.INIT_EVENTS | sdl.INIT_AUDIO
	TITLE  = "GB EMU"
)

type window struct {
	win        *sdl.Window
	renderer   *sdl.Renderer
	lgr        logger.Logger
	inputChan  chan<- sdl.Event // Write Only Channel for keyboard events
	bufferChan <-chan []uint8   // Read Only Channel for frame buffers
}

var EmuWindow window

func init() {
	EmuWindow = window{}
	EmuWindow.lgr = logger.NewLogger(os.Stdout, true, "Frontend")
	EmuWindow.lgr.Infof("Initiating Frontend\n")
	if e := sdl.Init(CONFIG); e != nil {
		EmuWindow.lgr.Fatalf("Failed to Initialize SDL")
	}
	sdl.SetHint(sdl.HINT_RENDER_VSYNC, "1")
	if win, rend, err := sdl.CreateWindowAndRenderer(TX_WIDTH, TX_HEIGHT, sdl.WINDOW_ALWAYS_ON_TOP|sdl.WINDOW_SHOWN); err != nil {
		EmuWindow.lgr.Fatalf(err.Error())
	} else {
		EmuWindow.lgr.Printf("Window Initialized with dimensions %d X %d\n", SC_WIDTH, SC_HEIGHT)

		win.SetSize(SC_WIDTH, SC_HEIGHT)
		win.SetBordered(true)
		win.SetResizable(true)
		win.SetTitle(TITLE)

		EmuWindow.win = win
		EmuWindow.renderer = rend
	}

	// go emuWindow.Run()
}

func (w *window) Run() {
	defer sdl.Quit()
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

func (w *window) SetChannels(bufferChan <-chan []uint8, inputChan chan<- sdl.Event) {
	w.bufferChan = bufferChan
	w.inputChan = inputChan
}
