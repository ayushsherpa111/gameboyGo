package frontend

import (
	"log"
	"os"

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
	f, e := os.OpenFile("vram.test", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0o777)
	if e != nil {
		w.lgr.Fatalf("Failed to open log file")
	}
	defer f.Close()
	defer sdl.Quit()
	defer EmuWindow.cleanUp()
	isRunning := true

	go w.listenForInput()
	var frames []uint32

	for isRunning {
		w.tex.UpdateRGBA(nil, frames, WIDTH)
		w.renderer.Copy(w.tex, nil, nil)
		w.renderer.Present()
		select {
		case frames = <-w.bufferChan:
		case key := <-w.SdlInpChan:
			w.inputChan <- key
			switch key {
			case sdl.K_q:
				w.lgr.Infof("Quitting")
				isRunning = false
				return
			}
		default:
		}
	}
}

func (w *window) SetChannels(bufferChan <-chan []uint32, inputChan chan<- sdl.Keycode) {
	w.bufferChan = bufferChan
	w.inputChan = inputChan
}
