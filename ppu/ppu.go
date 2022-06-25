package ppu

import (
	"errors"
	"fmt"
	"os"

	"github.com/ayushsherpa111/gameboyEMU/logger"
)

type tilePixelValue struct {
	low  uint8
	high uint8
}

const (
	LCD_STAT_HBLANK      uint8 = 0x00
	LCD_STAT_VBLANK      uint8 = 0x01
	LCD_STAT_OAM_RAM     uint8 = 0x02
	LCD_STAT_DATA2DRIVER uint8 = 0x03
	LCD_STAT_COINC       uint8 = 0x04 // R
)

const (
	LCD_STAT_INT_COINC  uint8 = 0x40 // R/W
	LCD_STAT_INT_OAM    uint8 = 0x20 // R/W
	LCD_STAT_INT_VBLANK uint8 = 0x10 // R/W
	LCD_STAT_INT_HBLANK uint8 = 0x08 // R/W
)

// compose buffer using the tileset and the vram
// 1 tile (8x8)px = 16 bytes
const (
	BUF_X = 160
	BUF_Y = 144

	BG_SET_X = 32
	BG_SET_Y = 32
)

var pallete = []uint32{0xFFFFFFFF,
	0x44444444,
	0xAAAAAAAA,
	0x00000000,
}

const (
	BG_WIN_ENABLE uint8 = 0x1 << iota
	OBJ_ENABLE
	OBJ_SIZE
	BG_TILE_MAP
	BG_WIN_DATA
	WIN_ENABLE
	WIN_TILE_MAP
	LCD_PPU_ENABLE
)

type PPU_MODE uint8

var PPU_BASE uint16 = 0xFF40
var V_RAM_START uint16 = 0x8000

const (
	LCD_C = 0xFF40
	LCD_S = 0xFF41
	ScY   = 0xFF42 // Determines the Viewport Y coordinate
	ScX   = 0xFF43 // Determines the Viewport X coordinate
	Ly    = 0xFF44
	LyC   = 0xFF45
	Bgp   = 0xFF47
	Wy    = 0xFF4A
	Wx    = 0xFF4B
)

const (
	V_BLANK_INT uint8 = 0x1
	LCD_INT     uint8 = 0x2
)

// TODO: add a tile data addressing mode struct. To store base pointer and sign mode

// INFO: BG map contains number denoting index on where tile data should be added.

type oam struct {
	yPOS    uint8 // byte 1
	xPOS    uint8 // byte 2
	tileIdx uint8 // byte 3

	/*
	 Bit7   BG and Window over OBJ (0=No, 1=BG and Window colors 1-3 over the OBJ)
	 Bit6   Y flip          (0=Normal, 1=Vertically mirrored)
	 Bit5   X flip          (0=Normal, 1=Horizontally mirrored)
	 Bit4   Palette number  **Non CGB Mode Only** (0=OBP0, 1=OBP1)
	 Bit3   Tile VRAM-Bank  **CGB Mode Only**     (0=Bank 0, 1=Bank 1)
	 Bit2-0 Palette number  **CGB Mode Only**     (OBP0-7)
	*/
	flags uint8 // byte 4
}

type ppu struct {
	vRAM          []uint8
	oam           []uint8
	oam_entries   [10]oam
	canvas_buffer [BUF_X * BUF_Y]uint32
	lgr           logger.Logger
	// background    [BG_SET_X][BG_SET_Y]uint8
	window  [][]uint8
	tileset []uint8

	ppu_regs []uint8 // 0xFF40 - 0xFF4B
	mode     PPU_MODE
	dots     uint16 // ticks to determine the mode of the PPU
	bufChan  chan<- []uint32
	IF       *uint8
}

var defaultVal uint8 = 0xFF

func NewPPU(bufferChan chan<- []uint32) *ppu {
	p := &ppu{
		vRAM:          make([]uint8, 8*1024),
		oam:           make([]uint8, 160),
		oam_entries:   [10]oam{},
		lgr:           logger.NewLogger(os.Stdout, true, "PPU"),
		canvas_buffer: [BUF_X * BUF_Y]uint32{},
		ppu_regs:      make([]uint8, 12),
		bufChan:       bufferChan,
		dots:          0,
	}
	return p
}

func (p *ppu) SortOAM() {
	for i := 0; i < len(p.oam_entries); i++ {
		for j := i + 1; j < len(p.oam_entries); j++ {
			if p.oam_entries[i].xPOS > p.oam_entries[j].xPOS {
				p.oam_entries[i], p.oam_entries[j] = p.oam_entries[j], p.oam_entries[i]
			}
		}
	}
}

func (p *ppu) fetchSprites(ly uint8) {
	spriteCount := 0
	for i := 0; i < len(p.oam) && spriteCount < 10; i += 4 {
		newOAM := oam{
			yPOS:    p.oam[i],
			xPOS:    p.oam[i+1],
			tileIdx: p.oam[i+2],
			flags:   p.oam[i+3],
		}
		if newOAM.yPOS <= ly && newOAM.yPOS+8 <= ly {
			p.oam_entries[spriteCount] = newOAM
			spriteCount++
		}
	}
	p.SortOAM()
}

func (p *ppu) PrintDetails() {
	file, err := os.OpenFile("vram.test", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0o777)
	if err != nil {
		p.lgr.Fatalf("Failed to open vram.test %s", err.Error())
	}
	fmt.Fprintf(file, "%v", p.vRAM)
	p.lgr.Infof(" LCDC: 0x%x\n", p.ppu_regs[parseIdx(LCD_C, PPU_BASE)])
}

func ParsePx(low, high uint8) []uint32 {
	tileRow := make([]uint32, 8)
	for i := 7; i >= 0; i-- {
		bitlow, bitHigh := (low&(0x1<<i))>>i, (high&(0x1<<i))>>i
		tileRow[7-i] = pallete[bitHigh<<1|bitlow]
	}
	return tileRow
}

func SetPx(x, y int, color uint32, buffer []uint32) {
	newPX := (BUF_X * y) + x
	buffer[newPX] = color
}

func (p *ppu) UpdateGPU() {
	// TODO: Check if the PPU is currently in a V_BLANK mode before entering mode 0

	lcd_c := &p.ppu_regs[parseIdx(LCD_C, PPU_BASE)]
	lcd_s := &p.ppu_regs[parseIdx(LCD_S, PPU_BASE)]
	lY := &p.ppu_regs[parseIdx(Ly, PPU_BASE)]
	lYc := &p.ppu_regs[parseIdx(LyC, PPU_BASE)]
	wY := &p.ppu_regs[parseIdx(Wy, PPU_BASE)]
	wX := &p.ppu_regs[parseIdx(Wx, PPU_BASE)]
	scY := &p.ppu_regs[parseIdx(ScY, PPU_BASE)]
	scX := &p.ppu_regs[parseIdx(ScX, PPU_BASE)]

	if *lcd_c&LCD_PPU_ENABLE == 0 {
		// LCD off
		p.dots = 0
		*lY = 0
		*lcd_s = setMode(*lcd_s, LCD_STAT_HBLANK)
		return
	}

	if *lY == *lYc {
		*lcd_s |= LCD_STAT_COINC

		if *lcd_s&LCD_STAT_INT_COINC != 0 {
			p.setInterrupt(LCD_INT)
		}
	} else {
		*lcd_s &= ^LCD_STAT_COINC
	}

	if p.dots >= 456 {
		*lY++
		*lY %= 154
		// checking for V_BLANK
		if *lY >= 144 {
			if *lY == 144 {
				p.setInterrupt(V_BLANK_INT)
			}
			*lcd_s = setMode(*lcd_s, LCD_STAT_VBLANK)
			if *lcd_s&LCD_STAT_INT_VBLANK != 0 {
				p.setInterrupt(LCD_INT)
			}
			// send frame buffer
			p.bufChan <- p.canvas_buffer[:]
		}
	}

	if *lY < 144 {
		if p.dots == 80 {
			// build sprite array from OAM
			*lcd_s = setMode(*lcd_s, LCD_STAT_OAM_RAM)
			if *lcd_s&LCD_STAT_INT_OAM != 0 {
				p.setInterrupt(LCD_INT)
			}
			p.fetchSprites(*lY)
		} else if p.dots == (80 + 172) {
			*lcd_s = setMode(*lcd_s, LCD_STAT_DATA2DRIVER)
			p.scanLine(lcd_c, wY, wX, scY, scX, lY)
		} else if p.dots == 80+172+204 {
			*lcd_s = setMode(*lcd_s, LCD_STAT_HBLANK)
			if *lcd_s&LCD_STAT_INT_HBLANK != 0 {
				p.setInterrupt(LCD_INT)
			}
		}
	}
	p.dots++
	p.dots = p.dots % 457
}

func (p *ppu) getSlice(start, end uint16) []uint8 {
	if start > end {
		p.lgr.Fatalf("Invalid slice indexes START:%d END:%d", start, end)
	}

	newStart := parseIdx(start, V_RAM_START)
	newEnd := parseIdx(end, V_RAM_START)
	return p.vRAM[newStart : newEnd+1]
}

func (p *ppu) drawBackgroundAndWin(lcdc, ly, wY, wX, scY, scX *uint8) {
	var tileYpos uint8
	var winBgTileData []uint8
	var signed bool
	var bgWinMap []uint8
	var drawWindow bool
	// objSize := (*lcdc & OBJ_SIZE) != 0

	if (*lcdc&WIN_ENABLE > 0) && (*ly > *wY) {
		drawWindow = true
		if *lcdc&WIN_TILE_MAP == WIN_TILE_MAP {
			bgWinMap = p.getSlice(0x9C00, 0x9FFF)
		} else {
			bgWinMap = p.getSlice(0x9800, 0x9BFF)
		}
		tileYpos = *ly - *wY
	} else {
		if *lcdc&BG_TILE_MAP == BG_TILE_MAP {
			bgWinMap = p.getSlice(0x9C00, 0x9FFF)
		} else {
			bgWinMap = p.getSlice(0x9800, 0x9BFF)
		}
		tileYpos = *ly + *scY
	}

	if *lcdc&BG_WIN_DATA == BG_WIN_DATA {
		winBgTileData = p.getSlice(0x8000, 0x8FFF)
	} else {
		winBgTileData = p.getSlice(0x8800, 0x97FF)
		signed = true
	}

	var x uint8

	for idx := byte(0); idx < BUF_X; idx += 8 {
		if drawWindow && idx > (*wX-7) {
			x = idx - (*wX - 7)
		} else {
			x = idx + *scX
		}

		var tileNum uint16 = ((uint16(tileYpos) / 8) * 32) + (uint16(x) / 8) // convert PX value to tile number
		tileIDX := bgWinMap[tileNum]                                         // tile number that is supposed to be drawn
		var tileDataAddr uint16
		// fetch the tile from tile dataset
		if signed {
			//  0x8800 + tileDataAddr
			tileDataAddr = uint16((int16(int8(tileIDX)) + 128) * 16) // by adding 128 you map -128 -> 0
		} else {
			// 0x8000 + tileDataAddr
			tileDataAddr = uint16(tileIDX) * 16
		}
		// tileDataAddr points to the first byte of the tile. a tile contains 16 byte (each line of the tile is made up of 2 bytes).
		// Next we find out which line of the current tile are we on
		tileDataY := uint16(uint8(tileYpos%8) * 2)
		low := winBgTileData[tileDataAddr+tileDataY]
		high := winBgTileData[tileDataAddr+tileDataY+1]

		for i := uint8(0); i < 8; i++ {
			idxCoord := BUF_X*uint(*ly) + uint(idx) + uint(i)
			p.canvas_buffer[idxCoord] = constructPixel(low, high, 7-i)
		}
	}
}

func (p *ppu) drawObjects() {
	p.lgr.Printf("DRAW SPRITES")
}

func (p *ppu) scanLine(lcdc, wY, wX, scY, scX, ly *uint8) {
	// render Background
	// p.lgr.Printf("Scanline started %d\n", *ly)
	if *lcdc&BG_WIN_ENABLE == BG_WIN_ENABLE {
		// draw either the background or the window
		p.drawBackgroundAndWin(lcdc, ly, wY, wX, scY, scX)
	}

	if *lcdc&OBJ_ENABLE == OBJ_ENABLE {
		p.drawObjects()
	}
}

func (p *ppu) Read_VRAM(addr uint16) *uint8 {
	lcd_s := p.ppu_regs[parseIdx(LCD_S, PPU_BASE)]
	if lcd_s&LCD_STAT_VBLANK != 0 {
		return &defaultVal
	}
	return &p.vRAM[addr]
}

func constructPixel(low, high, bitNum uint8) uint32 {
	var pxVal uint32
	bitlow, bitHigh := (low&(0x1<<bitNum))>>bitNum, (high&(0x1<<bitNum))>>bitNum
	pxVal = pallete[bitHigh<<1|bitlow]
	return pxVal
}

func (p *ppu) Write_VRAM(addr uint16, val uint8) {
	lcd_s := p.ppu_regs[parseIdx(LCD_S, PPU_BASE)]
	if lcd_s&LCD_STAT_DATA2DRIVER != 0 {
		// prevent CPU from writing to memory
		return
	}
	p.vRAM[addr] = val
}

func (p *ppu) Read_OAM(addr uint16) *uint8 {
	return &p.oam[addr]
}

func (p *ppu) Write_OAM(addr uint16, val uint8) {
	lcd_s := p.ppu_regs[parseIdx(LCD_S, PPU_BASE)]
	if lcd_s&LCD_STAT_OAM_RAM != 0 {
		return
	}
	p.oam[addr] = val
}

func (p *ppu) Read_Regs(regAddr uint16) *uint8 {
	newIdx := parseIdx(regAddr, PPU_BASE)
	return &p.ppu_regs[newIdx]
}

func (p *ppu) Write_Regs(regAddr uint16, val uint8) error {

	newIdx := parseIdx(regAddr, PPU_BASE)
	if int(newIdx) > len(p.ppu_regs) {
		return errors.New(fmt.Sprintf("Invalid memory address %X\n", regAddr))
	}

	p.ppu_regs[newIdx] = val
	return nil
}

func (p *ppu) clearInterrupt(interrupt uint8) {
	*p.IF &= ^interrupt
}

func (p *ppu) setInterrupt(interrupt uint8) {
	*p.IF |= interrupt
}

func (p *ppu) RefInterruptFlag(IF *uint8) {
	p.IF = IF
}

func setMode(flag, mode uint8) uint8 {
	flag &= 0xFC
	flag |= mode
	return flag
}

func parseIdx(idx uint16, baseAddr uint16) uint32 {
	return uint32(idx - baseAddr)
}
