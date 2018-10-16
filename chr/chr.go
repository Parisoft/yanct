package chr

// Cleanup8x8Tiles removes empty and duplicated tiles
func Cleanup8x8Tiles(tileset Tileset, metasprite Metasprite) {
	removeEmpty8x8Tiles(tileset, metasprite)
	removeDuplicated8x8Tiles(tileset, metasprite)
}

// Cleanup8x16Tiles removes empty and duplicated tiles
func Cleanup8x16Tiles(tileset Tileset, metasprite Metasprite) (Tileset, Metasprite) {
	tileset, metasprite = removeEmpty8x16Tiles(tileset, metasprite)
	tileset = removeDuplicated8x16Tiles(tileset, metasprite)

	return tileset, metasprite
}

func removeEmpty8x8Tiles(tileset Tileset, metasprite Metasprite) {
	for idx := len(tileset) - 1; idx >= 0; idx-- {
		if tileset[idx].Empty() {
			tileset.RemoveAt(idx)

			for i, spr := range metasprite {
				if spr.Idx == byte(idx) {
					metasprite.RemoveAt(i)
				} else if spr.Idx > byte(idx) {
					spr.Idx--
				}
			}
		}
	}
}

func removeEmpty8x16Tiles(tileset Tileset, metasprite Metasprite) (Tileset, Metasprite) {
	for idx := len(tileset) - 2; idx >= 0; idx -= 2 {
		if tileset[idx].Empty() && tileset[idx+1].Empty() {
			tileset = tileset.RemoveAt(idx + 1)
			tileset = tileset.RemoveAt(idx)

			for i := len(metasprite)-1;i>=0; i-- {
				spr := metasprite[i]
				if spr.Idx == byte(idx) {
					metasprite = metasprite.RemoveAt(i)
				} else if spr.Idx > byte(idx) {
					spr.Idx -= 2
				}
			}
		}
	}

	return tileset, metasprite
}

func removeDuplicated8x8Tiles(tileset Tileset, metasprites ...Metasprite) {
	for i := len(tileset) - 1; i >= 0; i-- {
		for j := 0; j < i; j++ {
			if tileset[i].Equals(&tileset[j]) {
				tileset.RemoveAt(i)

				for _, metasprite := range metasprites {
					for _, spr := range metasprite {
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

func removeDuplicated8x16Tiles(tileset Tileset, metasprites ...Metasprite) Tileset {
	for i := len(tileset) - 2; i >= 0; i -= 2 {
		for j := 0; j < i; j += 2 {
			if tileset[i].Equals(&tileset[j]) && tileset[i+1].Equals(&tileset[j+1]) {
				tileset = tileset.RemoveAt(i + 1)
				tileset = tileset.RemoveAt(i)

				for _, metasprite := range metasprites {
					for _, spr := range metasprite {
						if spr.Idx == byte(i) {
							spr.Idx = byte(j)
						}
					}
				}

				break
			}
		}
	}

	return tileset
}
