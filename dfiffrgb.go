package main

import (
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"sync"
)

func readRGB(img image.Image) []uint32 {
	dx := img.Bounds().Dx()
	dy := img.Bounds().Dy()
	mx := img.Bounds().Min.X
	my := img.Bounds().Min.Y
	buf := make([]uint32, dx*dy*3)
	for y := 0; y < dy; y++ {
		for x := 0; x < dx; x++ {
			p := (x*y + x) * 3
			r, g, b, _ := img.At(mx+x, my+y).RGBA()
			// fmt.Println(x, y, r, g, b)
			buf[p] = r
			buf[p+1] = g
			buf[p+2] = b
		}
	}
	// fmt.Println(img)
	// fmt.Println(buf)
	return buf
}

func diffRGB(left, right []uint32) float64 {
	j := float64(0)
	for i := 0; i < len(left); i++ {
		if left[i] == right[i] {
			j++
		}
	}
	return j / float64(len(left))
}

func puzzleDiffRGB(source, part, output string) error {

	log.Println("[+] puzzleDiffRGB")

	_, pWidth, pHeight := readImage(part)
	imgSrc, sWidth, sHeight := readImage(source)
	m = intRound(sWidth / pWidth)
	n = intRound(sHeight / pHeight)
	log.Printf("[+] m = %d, n = %d\n", m, n)

	// 原图 - 切割 - 哈希
	imgMatrix = make([][]*ImageMatrixData, m, m)
	for i := 0; i < m; i++ {
		imgMatrix[i] = make([]*ImageMatrixData, n, n)
	}
	log.Printf("[+] Cut and Hash Source Image\n")
	for x := 0; x < m; x++ {
		for y := 0; y < n; y++ {
			subImg, err := clipImage(imgSrc, x*pWidth, y*pHeight, (x+1)*pWidth, (y+1)*pHeight)
			if err != nil {
				log.Println("[-] clipImage:", err)
			}
			buff := readRGB(subImg)
			imgMatrix[x][y] = &ImageMatrixData{
				X: x * pWidth, Y: y * pHeight,
				Img: subImg, Buff: buff,
			}
		}
	}

	wg := &sync.WaitGroup{}

	// 碎片 - 哈希
	log.Printf("[+] Hash Part Image\n")
	pathDir := getParentDirectory(part)
	parts, _ := getAllParts(pathDir)
	flagMatrix = make([]*ImageMatrixData, len(parts))
	for n, partPath := range parts {
		wg.Add(1)
		go func(wg *sync.WaitGroup, n int, partPath string) {
			defer wg.Done()
			imgFile, err := os.Open(partPath)
			if err != nil {
				log.Println(n, partPath, err)
			}
			defer imgFile.Close()
			img, _ := png.Decode(imgFile)
			buff := readRGB(img)
			flagMatrix[n] = &ImageMatrixData{
				Name: partPath, Buff: buff, Img: img,
			}
		}(wg, n, partPath)
		if n%100 == 0 {
			wg.Wait()
		}
	}
	wg.Wait()

	extBlock := image.NewNRGBA(image.Rect(0, 0, pWidth, pHeight))
	// TODO: color form extColor
	draw.Draw(extBlock, extBlock.Bounds(), image.Black, image.ZP, draw.Src)
	// blockDistance := float64(pWidth*pHeight) * 0.8

	var ch = make(chan bool, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		p := 0
		for {
			select {
			case <-ch:
				p++
				log.Println("[+] finish :", p, "/", m)
				if p == m {
					close(ch)
					return
				}
			}
		}
	}()

	// Compare Image
	log.Printf("[+] Compare Image\n")
	for x := 0; x < m; x++ {
		imX := imgMatrix[x]
		wg.Add(1)
		go func(x int, imX []*ImageMatrixData) {
			defer wg.Done()
			for y := 0; y < n; y++ {
				distance := float64(0)
				im := imX[y]
				for _, fm := range flagMatrix {
					similar := diffRGB(im.Buff, fm.Buff)
					if similar > distance {
						distance = similar
						im.Img = fm.Img
						im.Name = fm.Name
					}
				}
			}
			ch <- true
		}(x, imX)
	}
	wg.Wait()
	log.Println("[+] Finish Compare")

	// Merge Image
	log.Printf("[+] Merge Image\n")
	err := mergeImage(imgMatrix, output, sWidth, sHeight, m, n)
	if err != nil {
		return err
	}

	return nil
}
