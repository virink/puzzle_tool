package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"math"
	"os"
	"strings"

	"github.com/corona10/goimagehash/transforms"
	"github.com/nfnt/resize"
)

func decodeImage(filePath string) (image.Image, error) {
	reader, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func intRound(i int) int {
	return int(math.Round(float64(i)))
}

func getAllParts(dirPath string) ([]string, error) {
	res := []string{}
	rd, err := ioutil.ReadDir(dirPath)
	for _, fi := range rd {
		if !fi.IsDir() {
			res = append(res, dirPath+"/"+fi.Name())
		}
	}
	return res, err
}

func getParentDirectory(dirctory string) string {
	runes := []rune(dirctory)
	l := strings.LastIndex(dirctory, "/")
	if l > len(runes) {
		l = len(runes)
	} else if l == -1 {
		return ""
	}
	return string(runes[0:l])
}

func readImage(filename string) (image.Image, int, int) {
	img, _ := decodeImage(filename)
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	return img, w, h
}

func clipImage(img image.Image, x0, y0, x1, y1 int) (image.Image, error) {
	switch img.(type) {
	case *image.NRGBA:
		img := img.(*image.NRGBA)
		subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.NRGBA)
		return subImg, nil
	case *image.RGBA:
		img := img.(*image.RGBA)
		subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.RGBA)
		return subImg, nil
	}
	return nil, errors.New("not support image type")
}

func grayImage(m image.Image) *image.RGBA {
	bounds := m.Bounds()
	dx := bounds.Dx()
	dy := bounds.Dy()
	newRgba := image.NewRGBA(bounds)
	for i := 0; i < dx; i++ {
		for j := 0; j < dy; j++ {
			_, g, _, a := m.At(i, j).RGBA()
			gUint8 := uint8((g >> 8))
			aUint8 := uint8(a >> 8)
			newRgba.SetRGBA(i, j, color.RGBA{gUint8, gUint8, gUint8, aUint8})
		}
	}
	return newRgba
}

func hammingDistance(lhash, rhash string) (int, error) {
	if len(lhash) != len(rhash) {
		return -1, errors.New("hash error")
	}
	n := 0
	for i := 0; i < len(lhash); i++ {
		if lhash[i] != rhash[i] {
			n++
		}
	}
	return n, nil
}

func differenceHash(img image.Image, size uint) (string, error) {
	if img == nil {
		return "", errors.New("Image object can not be nil")
	}
	resized := resize.Resize(size+1, size, img, resize.Bilinear)
	var hash bytes.Buffer
	pixels := transforms.Rgb2Gray(resized)
	for i := 0; i < len(pixels); i++ {
		for j := 0; j < len(pixels[i])-1; j++ {
			if pixels[i][j] < pixels[i][j+1] {
				hash.WriteByte('1')
			} else {
				hash.WriteByte('0')
			}
		}
	}
	return hash.String(), nil
}

func mergeImage(imgMatrix [][]*ImageMatrixData, output string, sWidth, sHeight, m, n int) error {
	resImage := image.NewNRGBA(image.Rect(0, 0, sWidth, sHeight))
	for x := 0; x < m; x++ {
		for y := 0; y < n; y++ {
			im := imgMatrix[x][y]
			if im.Img != nil {
				draw.Draw(resImage, resImage.Bounds(), im.Img, im.Img.Bounds().Min.Sub(image.Pt(im.X, im.Y)), draw.Over)
			} else {
				fmt.Println("orz")
			}
		}
	}
	out, err := os.Create(output)
	if err != nil {
		return err
	}
	return png.Encode(out, resImage)
}
