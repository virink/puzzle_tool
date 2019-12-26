package main

import (
	"image"
	"image/color"
	"log"
	"math"
	"sort"
)

func hashPartImage(part string) []*ImageMatrixData {
	pathDir := getParentDirectory(part)
	parts, _ := getAllParts(pathDir)
	partsNum := len(parts)
	matrix := make([]*ImageMatrixData, partsNum)
	for n, partPath := range parts {
		img, _, _ := readImage(partPath)
		up, down, left, right := getRGBA(img)
		matrix[n] = &ImageMatrixData{
			Name:  partPath,
			Up:    up,
			Down:  down,
			Left:  left,
			Right: right,
		}
	}
	return matrix
}

func pixelsDistance(left, right color.RGBA) float64 {
	// NOTE: 基于RGB颜色空间相似度判断
	// TODO: LAB ? https://blog.csdn.net/qq_16564093/article/details/80698479
	return math.Sqrt(math.Pow(float64(left.R-right.R), 2) +
		math.Pow(float64(left.G-right.G), 2) +
		math.Pow(float64(left.B-right.B), 2))
}

func getRGBA(img image.Image) (up, down, left, right []color.RGBA) {
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	up = make([]color.RGBA, w)
	down = make([]color.RGBA, w)
	left = make([]color.RGBA, h)
	right = make([]color.RGBA, h)
	for x := 0; x < w; x++ {
		r, g, b, _ := img.At(x, 0).RGBA()
		up[x] = color.RGBA{uint8(r), uint8(g), uint8(b), 0}
		r, g, b, _ = img.At(x, h).RGBA()
		down[x] = color.RGBA{uint8(r), uint8(g), uint8(b), 0}
	}
	for y := 0; y < h; y++ {
		r, g, b, _ := img.At(0, y).RGBA()
		left[y] = color.RGBA{uint8(r), uint8(g), uint8(b), 0}
		r, g, b, _ = img.At(w, y).RGBA()
		right[y] = color.RGBA{uint8(r), uint8(g), uint8(b), 0}
	}
	return
}

func duplicateString(list []string) []string {
	set := make(map[string]int, 1)
	for _, v := range list {
		set[v] = 1
	}
	keys := make([]string, 0, len(set))
	for k := range set {
		keys = append(keys, k)
	}
	return keys
}

// 取 distance 最小的前 N 个
func getDistanceLimit(imgList map[string]float64, limitNum int) []string {
	distanceList := make([]float64, 0, len(imgList))
	distanceLimit := float64(0)
	names := []string{}
	// sort
	for _, v := range imgList {
		distanceList = append(distanceList, v)
	}
	sort.Float64s(distanceList)
	// fmt.Println(distanceList)
	// fmt.Println("limitNum = ", limitNum)
	// limit
	for i, k := range distanceList {
		if i > limitNum {
			distanceLimit = k
			break
		}
	}
	// fmt.Println("distanceLimit = ", distanceLimit)
	// names
	for k, v := range imgList {
		if v <= distanceLimit {
			names = append(names, k)
		}
	}
	return names
}

func puzzleFit(source, part, output string) error {

	_, pWidth, pHeight = readImage(part)
	_, sWidth, sHeight := readImage(source)
	m := intRound(sWidth / pWidth)
	n := intRound(sHeight / pHeight)
	log.Printf("[+] m = %d, n = %d\n", m, n)

	imgMatrix = make([][]*ImageMatrixData, m, m)
	for i := 0; i < m; i++ {
		imgMatrix[i] = make([]*ImageMatrixData, n, n)
	}

	// 碎片 - 哈希
	log.Printf("[+] Hash Parts Image\n")
	// flagMatrix := hashPartImage(part)

	// 寻找图片左边 - 左 != 右
	log.Printf("[+] Load and Compare Image First Left\n")

	return nil
}
