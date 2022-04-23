package frontend

import (
	"os"

	"github.com/ayushsherpa111/gameboyEMU/logger"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	win_conf = sdl.WINDOW_ALWAYS_ON_TOP | sdl.WINDOW_SHOWN
)

var tex_pixel_format = uint32(sdl.PIXELFORMAT_RGBA32)

type window struct {
	win        *sdl.Window
	renderer   *sdl.Renderer
	tex        *sdl.Texture
	lgr        logger.Logger
	inputChan  chan<- sdl.Event // Write Only Channel for keyboard events
	bufferChan <-chan []uint32  // Read Only Channel for frame buffers
}

func (w *window) cleanUp() {
	w.tex.Destroy()
	w.renderer.Destroy()
	w.win.Destroy()
	close(w.inputChan)
}

func (w *window) createTexture(txWidth, txHeight int32) error {
	tex, err := w.renderer.CreateTexture(tex_pixel_format, sdl.TEXTUREACCESS_STREAMING, txWidth, txHeight)
	if err != nil {
		return err
	}
	w.tex = tex
	return nil
}

func createWindow(width, height int32, bufferChan <-chan []uint32) (*window, error) {
	newWin := &window{
		bufferChan: bufferChan,
		lgr:        logger.NewLogger(os.Stdout, true, "Frontend"),
	}
	win, rend, err := sdl.CreateWindowAndRenderer(width, height, win_conf)
	if err != nil {
		return newWin, err
	}
	newWin.win = win
	newWin.renderer = rend
	return newWin, nil
}
