package chr

import (
	"fmt"
)

// Sprite defines a sprite
type Sprite struct {
	X   int8
	Y   int8
	Opt byte
	Idx byte
}

// Bytes transform a Sprite into an array of bytes in the format [x, y, opt, idx]
func (spr *Sprite) Bytes() [4]byte {
	bytes := [4]byte{}
	bytes[0] = byte(spr.X)
	bytes[1] = byte(spr.Y)
	bytes[2] = spr.Opt
	bytes[3] = spr.Idx
	return bytes
}

func (spr *Sprite) String() string {
	return fmt.Sprintf("%d, %d, 0x%x, %d", spr.X, spr.Y, spr.Idx, spr.Opt)
}
