package chr

import (
	"fmt"
)

// Metasprite is a table of sprites
type Metasprite []*Sprite

// NewMetasprite builds a metasprite from a Tileset
func NewMetasprite(tileset Tileset) Metasprite {
	var metasprite Metasprite //make(Metasprite, len(tileset))
	rows := len(tileset) / TilesetRows

	for row := 0; row < rows; row++ {
		for col := 0; col < TilesetCols; col++ {
			metasprite = append(metasprite, &Sprite{
				X:   int8(col * 8),
				Y:   int8((row - rows) * 8),
				Opt: 0,
				Idx: byte(row*TilesetCols + col),
			})
		}
	}

	return metasprite
}

//To8x16 returns a metasprite whith tiles converted to 8x16 pixels
func (metasprite Metasprite) To8x16() Metasprite {
	// remove odd lines
	for i := len(metasprite) - 1; i >= 0; i-- {
		if (metasprite[i].Idx/TilesetRows)%2 != 0 {
			metasprite = metasprite.RemoveAt(i)
		}
	}

	for _, spr := range metasprite {
		spr.Idx = (spr.Idx/16)*16 + (spr.Idx%16)*2
	}

	return metasprite
}

//Print a metasprite as a C array
func (metasprite Metasprite) Print() {
	fmt.Printf("const s8_t metasprite[%d] = {\n", len(metasprite)*4+1)
	for _, spr := range metasprite {
		fmt.Printf("%s, \n", spr.String())
	}
	fmt.Println("0x80};")
}

//RemoveAt remove an element at position i
func (metasprite Metasprite) RemoveAt(i int) Metasprite {
	return append(metasprite[:i], metasprite[i+1:]...)
}
