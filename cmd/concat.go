package cmd

import (
	"errors"
	"os"

	"github.com/parisoft/yanct/chr"

	"github.com/spf13/cobra"
)

var concatCmd = &cobra.Command{
	Use:   "concat CHR_1 CHR_2 [...CHR_N]",
	Short: "Concatenate many CHR files into one",
	Long: `Concatenate many CHR files into one.
All files are appended to the first one. After each append, all duplicated tiles are removed.
If a binary metasprite is found on the same path of a CHR file being appended, it's updated to follow the concatenated file.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("concat requires 2 CHR files or more")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateMetasprFmt(); err != nil {
			return err
		}
		if err := validateTileH(); err != nil {
			return err
		}
		if err := validateChrOut(); err != nil {
			return err
		}

		return concat(args...)
	},
}

func init() {
	concatCmd.Flags().Uint8VarP(&flg.tileH, FlgTileH, "t", 8, "Height of the tiles: 8 for 8x8, 16 for 8x16")
	concatCmd.Flags().StringVarP(&flg.chrOut, FlgChrOut, "o", "", "output CHR file name")
	concatCmd.MarkFlagRequired(FlgChrOut)
	rootCmd.AddCommand(concatCmd)
}

func concat(chrlist ...string) error {
	tiledim := chr.Tile8x8
	if flg.tileH == 16 {
		tiledim = chr.Tile8x16
	}
	tilesets := make([]*chr.Tileset, len(chrlist))
	metasprites := make([]*chr.Metasprite, len(chrlist))
	binnames := make([]string, len(chrlist))
	for i, chrfilename := range chrlist {
		chrfile, err := os.Open(chrfilename)
		if err != nil {
			return err
		}
		defer chrfile.Close()

		binfilename := chrfilename[:len(chrfilename)-3] + "bin"
		binnames[i] = binfilename
		binfile, err := os.Open(binfilename)
		if err == nil {
			defer binfile.Close()
			var metaspr *chr.Metasprite
			if metaspr, err = chr.NewMetaspriteFromFile(binfile); err != nil {
				return err
			}
			metasprites[i] = metaspr
		} else if !os.IsNotExist(err) {
			return err
		}

		var tileset *chr.Tileset
		if tileset, err = chr.NewTilesetFromCHR(chrfile, tiledim); err != nil {
			return err
		}
		tilesets[i] = tileset
	}

	output := chr.NewTileset(tiledim)
	for i, tileset := range tilesets {
		chr.ConcatTiles(output, tileset, metasprites[i])
	}

	output.Write(flg.chrOut)

	for i, metasprite := range metasprites {
		if metasprite != nil {
			metasprite.WriteBin(binnames[i])
		}
	}

	return nil
}
