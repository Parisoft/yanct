package chr

// Cleanup8x8Tiles removes empty and duplicated tiles
func Cleanup8x8Tiles(tileset *Tileset, metasprite *Metasprite) {
	removeEmpty8x8Tiles(tileset, metasprite)
	removeDuplicated8x8Tiles(tileset, metasprite)
}

// Cleanup8x16Tiles removes empty and duplicated tiles
func Cleanup8x16Tiles(tileset *Tileset, metasprite *Metasprite)  {
	removeEmpty8x16Tiles(tileset, metasprite)
	removeDuplicated8x16Tiles(tileset, metasprite)
}

func removeEmpty8x8Tiles(tileset *Tileset, metasprite *Metasprite) {
	for idx := tileset.Size() - 1; idx >= 0; idx-- {
		if tileset.At(idx).Empty() {
			tileset.RemoveAt(idx)

			for i := metasprite.Size() - 1; i >= 0; i-- {
				spr := metasprite.At(i)
				if spr.Idx == byte(idx) {
					metasprite.RemoveAt(i)
				} else if spr.Idx > byte(idx) {
					spr.Idx--
				}
			}
		}
	}
}

func removeEmpty8x16Tiles(tileset *Tileset, metasprite *Metasprite)  {
	for idx := tileset.Size() - 2; idx >= 0; idx -= 2 {
		if tileset.At(idx).Empty() && tileset.At(idx+1).Empty() {
			tileset.RemoveAt(idx + 1)
			tileset.RemoveAt(idx)

			for i := metasprite.Size() - 1; i >= 0; i-- {
				spr := metasprite.At(i)
				if spr.Idx == byte(idx) {
					metasprite.RemoveAt(i)
				} else if spr.Idx > byte(idx) {
					spr.Idx -= 2
				}
			}
		}
	}
}

func removeDuplicated8x8Tiles(tileset *Tileset, metasprites ...*Metasprite) {
	for i := tileset.Size() - 1; i >= 0; i-- {
		for j := 0; j < i; j++ {
			if tileset.At(i).Equals(tileset.At(j)) {
				tileset.RemoveAt(i)

				for _, metasprite := range metasprites {
					for _, spr := range metasprite.sprites {
						if spr.Idx == byte(i) {
							spr.Idx = byte(j)
						}
					}
				}

				break
			}
		}
	}

}

func removeDuplicated8x16Tiles(tileset *Tileset, metasprites ...*Metasprite) {
	for i := tileset.Size() - 2; i >= 0; i -= 2 {
		for j := 0; j < i; j += 2 {
			if tileset.At(i).Equals(tileset.At(j)) && tileset.At(i+1).Equals(tileset.At(j+1)) {
				tileset.RemoveAt(i + 1)
				tileset.RemoveAt(i)

				for _, metasprite := range metasprites {
					for _, spr := range metasprite.sprites {
						if spr.Idx == byte(i) {
							spr.Idx = byte(j)
						}
					}
				}

				break
			}
		}
	}
}
