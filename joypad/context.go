package joypad

const (
	actionBit    uint8 = 0x20
	directionBit uint8 = 0x10
)

type context struct {

	/*
		Selection bit
		Bit 5. (0xFF00)
	*/
	actionBit bool

	/*
		Selection bit
		Bit 4. (0xFF00)
	*/
	directionBit bool

	controller joypad
}

func (c context) SetSelBit(val uint8) {
	if (^val & actionBit) > 0 {
		c.actionBit = true
	}

	if (^val & directionBit) > 0 {
		c.directionBit = true
	}
}

func (c context) GetGamepadState() *uint8 {
	if c.directionBit {
		return &c.controller.directionBits
	} else {
		return &c.controller.actionBits
	}
	// return nil
}

func NewContext() context {
	return context{
		actionBit:    false,
		directionBit: false,
		controller: joypad{
			directionBits: 0,
			actionBits:    0,
		},
	}
}
