package main

import (
	"image"
	"image/png"
	"os"

	"github.com/yanct/chr"
)

func main() {
	pngfile, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	defer pngfile.Close()

	img, err := png.Decode(pngfile)
	if err != nil {
		panic(err)
	}

	pngimg, ok := img.(image.PalettedImage)
	if !ok {
		panic("input file must be a PNG indexed with a 4 color pallete")
	}

	tileset := chr.NewTileset(pngimg, 3)
	metasprite := chr.NewMetasprite(tileset)
	tileset = tileset.To8x16()
	metasprite = metasprite.To8x16()
	tileset, metasprite = chr.Cleanup8x16Tiles(tileset, metasprite)
	tileset.Write("/tmp/sprite-out-go.chr")
	metasprite.Print()
}
