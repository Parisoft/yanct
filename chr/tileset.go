package chr

import (
	"image"
	"os"
)

const (
	//TilesetRows is the number of rows
	TilesetRows = 16
	//TilesetCols is the number of cols
	TilesetCols = 16
)

//Tileset is a table of tiles
type Tileset []Tile

//NewTileset builds a tileset from an indexed PNG image
func NewTileset(img image.PalettedImage, bgColorIdx byte) Tileset {
	tileset := make(Tileset, 256)
	for i := 0; i < 256; i++ {
		tileset[i] = Tile{}
	}
	h := byte(img.Bounds().Dy())
	w := byte(img.Bounds().Dx())

	for row, y := TilesetRows-h/8, byte(0); y < h; row, y = row+1, y+8 {
		for col, x := byte(0), byte(0); x < w; col, x = col+1, x+8 {
			tile := &tileset[row*TilesetCols+col]
			pixels := pixels(x, y, img)

			//http://wiki.nesdev.com/w/index.php/PPU_pattern_tables
			for i := byte(0); i < 8; i++ {
				for j := byte(0); j < 8; j++ {
					pixel := pixels[i*8+j]
					if pixel == bgColorIdx {
						pixel = 0
					} else if pixel == 0 {
						pixel = bgColorIdx
					}

					tile.Plane[0][i] |= (pixel & 1) << (7 - j)
					tile.Plane[1][i] |= ((pixel & 2) >> 1) << (7 - j)
				}
			}
		}
	}

	return tileset
}

//To8x16 returns a tileset with tiles converted to 8x16 pixels
func (tileset Tileset) To8x16() Tileset {
	// move sprites at odd lines to 1 line above at odd column
	tmp := make(Tileset, len(tileset))

	for row, i := 0, -1; row < len(tileset)/TilesetRows; row += 2 {
		for col := 0; col < TilesetCols; col++ {
			idx := row*TilesetCols + col
			i++
			tmp[i] = tileset[idx]
			i++
			tmp[i] = tileset[idx+TilesetCols]
		}
	}

	return tmp
}

func (tileset Tileset) Write(filename string) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	for _, tile := range tileset {
		file.Write(tile.Plane[0][:])
		file.Write(tile.Plane[1][:])
	}
}

//RemoveAt remove an element at position i
func (tileset Tileset) RemoveAt(i int) Tileset {
	return append(tileset[:i], tileset[i+1:]...)
}

func pixels(x, y byte, img image.PalettedImage) []byte {
	pixels := make([]byte, 64)

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			pixels[i*8+j] = img.ColorIndexAt(int(x)+j, int(y)+i)
		}
	}

	return pixels
}
