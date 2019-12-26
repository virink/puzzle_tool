package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"time"

	"github.com/hedzr/cmdr"
)

const (
	// APPNAME Name for app
	APPNAME = "Puzzle - CTFTools"
	// VERSION Version for app
	VERSION = "0.0.1"
)

// ImageMatrixData - Image Matrix Data
type ImageMatrixData struct {
	X, Y        int
	Name        string
	Img         image.Image
	Part        string
	Status      bool
	Data        string
	Buff        []uint32
	Up, Down    []color.RGBA
	Left, Right []color.RGBA
	// RGB    color.RGBA
}

// ImageMatrixRelation - Image Matrix Relation
type ImageMatrixRelation struct {
	lImage *ImageMatrixData
	rImage *ImageMatrixData
}

var (
	reSize = uint(32)

	flagMatrix []*ImageMatrixData
	imgMatrix  [][]*ImageMatrixData

	// 碎片宽高
	pWidth  = 0
	pHeight = 0

	// 原图
	imgSrc  = ""
	sWidth  = 0
	sHeight = 0

	// 矩阵坐标
	m = 0
	n = 0
)

func main() {
	root := cmdr.Root(APPNAME, VERSION).
		Header("[$] Puzzle - Virink <virink@outlook.com>").
		Description("", "A tool for puzzle")
	rootCmd := root.RootCommand()

	// 字典处理
	jigsawCmd := root.NewSubCommand().
		Titles("j", "jigsaw").
		Description("", "Jigsaw Puzzles").
		Group("Jigsaw").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			jType := cmdr.GetInt("app.jigsaw.type", 1)
			reSize = cmdr.GetUint("app.jigsaw.resize", 32)
			source := cmdr.GetString("app.jigsaw.source")
			part := cmdr.GetString("app.jigsaw.part")
			output := cmdr.GetString("app.jigsaw.output")
			extColor := cmdr.GetString("app.jigsaw.ext", "#fff")
			if _, e := os.Stat(source); os.IsNotExist(e) {
				log.Println("Source file is not exist!")
				return
			}
			if _, e := os.Stat(part); os.IsNotExist(e) {
				log.Println("Part of puzzles file is not exist!")
				return
			}
			t := time.Now().UnixNano()

			os.Remove(output)
			log.Printf("[+] Source   : %s\n", source)
			log.Printf("[+] Part     : %s\n", part)
			log.Printf("[+] Output   : %s\n", output)
			log.Printf("[+] ExtColor : %s\n", extColor)

			if jType == 1 {
				puzzleDiffHash(source, part, output, extColor)
			} else if jType == 2 {
				// TODO: 无原图拼图，挖坑不填了
				puzzleFit(source, part, output)
			} else if jType == 3 {
				puzzleHash(source, part, output)
			} else if jType == 4 {
				puzzleDiffRGB(source, part, output)
			}

			t = time.Now().UnixNano() - t
			log.Println(fmt.Sprintf("[+] Success! And use %f s", float64(t)/1e9))
			return
		})
	jigsawCmd.NewFlag(cmdr.OptFlagTypeInt).
		Titles("r", "resize").
		Description("resize", ``).
		DefaultValue(8, "size")
	jigsawCmd.NewFlag(cmdr.OptFlagTypeInt).
		Titles("t", "type").
		Description("1:DiffHash 2:- 3:Hash 4:DiffRGB", ``).
		DefaultValue(1, "num")
	jigsawCmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("s", "source").
		Description("Read the source file (*.png)", ``).
		DefaultValue("origin.png", "FILE")
	jigsawCmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("o", "output").
		Description("Write the output file (*.db)", ``).
		DefaultValue("results.png", "FILE")
	jigsawCmd.NewFlag(cmdr.OptFlagTypeBool).
		Titles("p", "part").
		Description("Part of puzzles file", ``).
		DefaultValue("", "FILE")
	jigsawCmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("e", "ext").
		Description("Append ext block", ``).
		DefaultValue("#fff", "color")

	if err := cmdr.Exec(rootCmd); err != nil {
		log.Println("[-]", err)
	}

}
