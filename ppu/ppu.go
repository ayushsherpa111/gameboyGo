package ppu

// compose buffer using the tileset and the vram
// 1 tile (8x8)px = 16 bytes
const (
	TILE_SET_X = 20
	TILE_SET_Y = 18

	BG_SET_X = 32
	BG_SET_Y = 32
)

type PPU_MODE uint8

const (
	MODE_0 PPU_MODE = iota
	MODE_1
	MODE_2
	MODE_3
)

type ppu struct {
	vRAM          []uint8
	oam           []uint8
	canvas_buffer [TILE_SET_X][TILE_SET_Y]uint8
	background    [BG_SET_X][BG_SET_Y]uint8
	window        [][]uint8
	tileset       []uint8

	LCD_CONTROL uint8 // 0xFF40
	LCD_STATUS  uint8 // 0xFF41
	ScY         uint8 // 0xFF42
	ScX         uint8 // 0xFF43
	Ly          uint8 // 0xFF44
	LyC         uint8 // 0xFF45
	Bgp         uint8 // 0xFF47
	Wy          uint8 // 0xFF4A
	Wx          uint8 // 0xFF4B
	mode        PPU_MODE
	dots        uint16 // ticks to determine the mode of the PPU
}

func NewPPU() *ppu {
	p := &ppu{
		vRAM:          make([]uint8, 8*1024),
		oam:           make([]uint8, 159),
		canvas_buffer: [TILE_SET_X][TILE_SET_Y]uint8{},
		background:    [BG_SET_X][BG_SET_Y]uint8{},
	}
	return p
}

func (p *ppu) UpdateGPU() {
	panic("not implemented") // TODO: Implement
}

func (p *ppu) Read_VRAM(addr uint16) *uint8 {
	if p.mode == MODE_3 {
		return nil
	}
	return &p.vRAM[addr]
}

func (p *ppu) Write_VRAM(_ uint16, _ uint8) {
}

func (p *ppu) Read_OAM(addr uint16) *uint8 {
}

func (p *ppu) Write_OAM(_ uint16, _ uint8) {
}

func (p *ppu) Read_Regs() {

}

func (p *ppu) Write_Regs() {

}
