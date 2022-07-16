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
	}
	if c.actionBit {
		return &c.controller.actionBits
	}
	return nil
}

func NewContext() context {
	return context{
		actionBit:    false,
		directionBit: false,
		controller: joypad{
			directionBits: 0x0F,
			actionBits:    0x0F,
		},
	}
}

func (c context) HandleEvent(key uint8, state bool) {
	if (key & directionBit) != 0 {
		c.controller.SetDirection(key&0xF, state)
	} else if (key & actionBit) != 0 {
		c.controller.SetAction(key&0xF, state)
	}
}
