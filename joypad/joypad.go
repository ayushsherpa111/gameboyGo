package joypad

const (
	joypadMask  uint8 = 0x0F
	downOrStart uint8 = 0x08
	upOrSelect  uint8 = 0x04
	leftOrB     uint8 = 0x02
	rightOrA    uint8 = 0x01
)

type joypad struct {
	/*
		Get Directions when bit 4 is set
		0 = pressed
	*/
	directionBits uint8

	/*
		Get Actions when bit 5 is set
		0 = pressed
	*/
	actionBits uint8
}

func unmask(val, mask uint8) bool {
	if ^val&mask == mask {
		return true
	}
	return false
}

func (j joypad) SetDirection(direction uint8, state bool) {
	if state {
		j.directionBits &= direction
	} else {
		j.directionBits |= ((^direction) & 0x0F)
	}
}

func (j joypad) SetAction(action uint8, state bool) {
	if state {
		j.actionBits &= action
	} else {
		j.actionBits |= ((^action) & 0x0F)
	}
}

func (j joypad) GetStart() bool {
	return unmask(actionBit, downOrStart)
}

func (j joypad) GetDown() bool {
	return unmask(directionBit, downOrStart)
}

func (j joypad) GetSelect() bool {
	return unmask(actionBit, upOrSelect)
}

func (j joypad) GetUp() bool {
	return unmask(directionBit, upOrSelect)
}

func (j joypad) GetLeft() bool {
	return unmask(directionBit, leftOrB)
}

func (j joypad) GetB() bool {
	return unmask(actionBit, leftOrB)
}

func (j joypad) GetRight() bool {
	return unmask(directionBit, rightOrA)
}

func (j joypad) GetA() bool {
	return unmask(actionBit, rightOrA)
}
