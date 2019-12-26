package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"sync"
)

func puzzleDiffHash(source, part, output, extColor string) error {

	log.Println("[+] puzzleDiffHash")

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
			hashString, err := differenceHash(subImg, reSize)
			// fmt.Println(hashString, len(hashString))
			if err != nil {
				log.Println("[-] differenceHash:", err)
			}
			imgMatrix[x][y] = &ImageMatrixData{
				X: x * pWidth, Y: y * pHeight, Img: subImg,
				Data: hashString, Status: false,
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
			hasshString, _ := differenceHash(img, reSize)
			flagMatrix[n] = &ImageMatrixData{Name: partPath, Data: hasshString, Img: img}

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

	// Compare Image
	log.Printf("[+] Compare Image\n")
	for x := 0; x < m; x++ {
		for y := 0; y < n; y++ {
			distance := 99999
			im := imgMatrix[x][y]
			for i, fm := range flagMatrix {
				d, err := hammingDistance(im.Data, fm.Data)
				if err != nil {
					fmt.Println(err)
					break
				} else if d < 5 {
					distance = d
					im.Img = fm.Img
					flagMatrix = append(flagMatrix[:i], flagMatrix[i+1:]...)
					break
				} else if d < distance {
					distance = d
					im.Img = fm.Img
				}
			}
			// if distance > 50 {
			// 	im.Img = extBlock
			// }
		}
	}

	// Merge Image
	log.Printf("[+] Merge Image\n")
	err := mergeImage(imgMatrix, output, sWidth, sHeight, m, n)
	if err != nil {
		return err
	}

	return nil
}
