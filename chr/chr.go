package chr

import "strings"

const (
	spriteMirrorOpt = (1 << 6)
	spriteFlipOpt   = (1 << 7)
)

//CleanupTiles removes empty and duplicated tiles
func CleanupTiles(tileset *Tileset, metasprite *Metasprite, delMirror, delFlip bool) {
	if tileset.tiledim == Tile8x16 {
		removeEmpty8x16Tiles(tileset, metasprite)
		removeDuplicated8x16Tiles(tileset, metasprite, delMirror, delFlip)
	} else {
		removeEmpty8x8Tiles(tileset, metasprite)
		removeDuplicated8x8Tiles(tileset, metasprite, delMirror, delFlip)
	}
}

//ConcatTiles concatenate the 2nd tileset onto the 1st tileset, updating those respective metrasprites
func ConcatTiles(tileset1, tileset2 *Tileset, metasprite2 *Metasprite, delMirror, delFlip bool) {
	if metasprite2 != nil {
		len := byte(tileset1.Size())
		for _, spr := range metasprite2.sprites {
			spr.Idx += len
		}
	}

	tileset1.tiles = append(tileset1.tiles, tileset2.tiles...)

	if tileset1.tiledim == Tile8x16 {
		removeDuplicated8x16Tiles(tileset1, metasprite2, delMirror, delFlip)
	} else {
		removeDuplicated8x8Tiles(tileset1, metasprite2, delMirror, delFlip)
	}
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

func removeEmpty8x16Tiles(tileset *Tileset, metasprite *Metasprite) {
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

func removeDuplicated8x8Tiles(tileset *Tileset, metasprite *Metasprite, delMirror, delFlip bool) {
	for i := tileset.Size() - 1; i >= 0; i-- {
		for j := 0; j < i; j++ {
			if tileset.At(i).Equals(tileset.At(j)) {
				tileset.RemoveAt(i)

				if metasprite != nil {
					for _, spr := range metasprite.sprites {
						if spr.Idx == byte(i) {
							spr.Idx = byte(j)
						} else if spr.Idx > byte(i) {
							spr.Idx--
						}
					}
				}

				break
			}

			if delMirror && tileset.At(i).Mirrored(tileset.At(j)) {
				tileset.RemoveAt(i)

				if metasprite != nil {
					for _, spr := range metasprite.sprites {
						if spr.Idx == byte(i) {
							spr.Idx = byte(j)
							spr.Opt ^= spriteMirrorOpt
						} else if spr.Idx > byte(i) {
							spr.Idx--
						}
					}
				}

				break
			}

			if delFlip && tileset.At(i).Flipped(tileset.At(j)) {
				tileset.RemoveAt(i)

				if metasprite != nil {
					for _, spr := range metasprite.sprites {
						if spr.Idx == byte(i) {
							spr.Idx = byte(j)
							spr.Opt ^= spriteFlipOpt
						} else if spr.Idx > byte(i) {
							spr.Idx--
						}
					}
				}

				break
			}

			if delMirror && delFlip && tileset.At(i).MirrorFlipped(tileset.At(j)) {
				tileset.RemoveAt(i)

				if metasprite != nil {
					for _, spr := range metasprite.sprites {
						if spr.Idx == byte(i) {
							spr.Idx = byte(j)
							spr.Opt ^= (spriteFlipOpt | spriteMirrorOpt)
						} else if spr.Idx > byte(i) {
							spr.Idx--
						}
					}
				}

				break
			}
		}
	}

}

func removeDuplicated8x16Tiles(tileset *Tileset, metasprite *Metasprite, delMirror, delFlip bool) {
	for i := tileset.Size() - 2; i >= 0; i -= 2 {
		for j := 0; j < i; j += 2 {
			if tileset.At(i).Equals(tileset.At(j)) && tileset.At(i+1).Equals(tileset.At(j+1)) {
				tileset.RemoveAt(i + 1)
				tileset.RemoveAt(i)

				if metasprite != nil {
					for _, spr := range metasprite.sprites {
						if spr.Idx == byte(i) {
							spr.Idx = byte(j)
						} else if spr.Idx > byte(i) {
							spr.Idx -= 2
						}
					}
				}

				break
			}

			if delMirror && tileset.At(i).Mirrored(tileset.At(j)) && tileset.At(i+1).Mirrored(tileset.At(j+1)) {
				tileset.RemoveAt(i + 1)
				tileset.RemoveAt(i)

				if metasprite != nil {
					for _, spr := range metasprite.sprites {
						if spr.Idx == byte(i) {
							spr.Idx = byte(j)
							spr.Opt ^= spriteMirrorOpt
						} else if spr.Idx > byte(i) {
							spr.Idx -= 2
						}
					}
				}

				break
			}

			if delFlip && tileset.At(i).Flipped(tileset.At(j+1)) && tileset.At(i+1).Flipped(tileset.At(j)) {
				tileset.RemoveAt(i + 1)
				tileset.RemoveAt(i)

				if metasprite != nil {
					for _, spr := range metasprite.sprites {
						if spr.Idx == byte(i) {
							spr.Idx = byte(j)
							spr.Opt ^= spriteFlipOpt
						} else if spr.Idx > byte(i) {
							spr.Idx -= 2
						}
					}
				}

				break
			}

			if delMirror && delFlip && tileset.At(i).MirrorFlipped(tileset.At(j+1)) && tileset.At(i+1).MirrorFlipped(tileset.At(j)) {
				tileset.RemoveAt(i + 1)
				tileset.RemoveAt(i)

				if metasprite != nil {
					for _, spr := range metasprite.sprites {
						if spr.Idx == byte(i) {
							spr.Idx = byte(j)
							spr.Opt ^= (spriteFlipOpt | spriteMirrorOpt)
						} else if spr.Idx > byte(i) {
							spr.Idx -= 2
						}
					}
				}

				break
			}
		}
	}
}

func changeFileExtension(name, extension string) string {
	if dot := strings.LastIndex(name, "."); dot > -1 {
		return name[:dot] + "." + extension
	}
	return name + "." + extension
}
