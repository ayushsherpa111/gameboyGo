package frontend

import (
	"fmt"
	"log"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	SCALE = 4

	WIDTH  = 160 // Width of the Texture
	HEIGHT = 144 // Height of the Texture

	CONFIG = sdl.INIT_VIDEO | sdl.INIT_EVENTS | sdl.INIT_AUDIO
	TITLE  = "GB EMU"
)

var EmuWindow *window
var DebugWindow *window

func SetupWindow() {
	var err error
	if EmuWindow, err = createWindow(WIDTH*SCALE, HEIGHT*SCALE, nil); err != nil {
		log.Fatalln(err.Error())
	} else {
		EmuWindow.lgr.Infof("Initiating Frontend\n")
		EmuWindow.lgr.Printf("Window Initialized with dimensions %d X %d\n", WIDTH, HEIGHT)

		if err = EmuWindow.createTexture(WIDTH, HEIGHT); err != nil {
			EmuWindow.lgr.Fatalf("Failed to Initialize Texture %s", err.Error())
		}
	}

}

func (w *window) Run() {
	defer sdl.Quit()
	defer EmuWindow.cleanUp()
	isRunning := true

	w.tex.UpdateRGBA(nil, w.winBuf, WIDTH)
	for isRunning {
		w.renderer.Copy(w.tex, nil, nil)
		w.renderer.Present()
		v := <-w.bufferChan
		fmt.Println("received")
		w.tex.UpdateRGBA(nil, v, WIDTH)
		// select {
		// case
		// default:
		// }
		// // handle key presses.
		// for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		// }
	}
}

func (w *window) SetChannels(bufferChan <-chan []uint32, inputChan chan<- sdl.Event) {
	w.bufferChan = bufferChan
	w.inputChan = inputChan
}
