package main

import (
	"crypto/md5"
	"encoding/hex"
	"image"
	"image/draw"
	"image/png"
	"io"
	"log"
	"os"
	"sync"
)

func imgMd5(img image.Image) (string, error) {
	md5 := md5.New()
	// io.Copy(md5, file)
	err := png.Encode(md5, img)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(md5.Sum(nil)), nil
}

func fileMd5(file io.Reader) (string, error) {
	md5 := md5.New()
	_, err := io.Copy(md5, file)
	// err := png.Encode(md5, img)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(md5.Sum(nil)), nil
}

func puzzleHash(source, part, output string) error {

	log.Println("[+] puzzleDiffHash")

	_, pWidth, pHeight = readImage(part)
	imgSrc, sWidth, sHeight := readImage(source)
	m := intRound(sWidth / pWidth)
	n := intRound(sHeight / pHeight)
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
			hashString, err := imgMd5(subImg)
			// fmt.Println("imgMatrix hashString :", hashString)
			if err != nil {
				log.Println("[-] Md5 Hash:", err)
			}
			imgMatrix[x][y] = &ImageMatrixData{
				X: x * pWidth, Y: y * pHeight, Img: nil,
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
			// hashString, _ := fileMd5(imgFile)
			img, _ := png.Decode(imgFile)
			hashString, err := imgMd5(img)
			// fmt.Println("flagMatrix hashString :", hashString)
			flagMatrix[n] = &ImageMatrixData{Name: partPath, Data: hashString, Img: img}

		}(wg, n, partPath)
		if n%100 == 0 {
			wg.Wait()
		}
	}
	wg.Wait()

	extBlock := image.NewNRGBA(image.Rect(0, 0, pWidth, pHeight))
	// TODO: color form extColor
	draw.Draw(extBlock, extBlock.Bounds(), image.Black, image.ZP, draw.Src)

	diffNum := 0
	diffMatrix := make([]*ImageMatrixData, len(parts))

	for x := 0; x < m; x++ {
		for y := 0; y < n; y++ {
			im := imgMatrix[x][y]
			skip := true
			for i, fm := range flagMatrix {
				if im.Data == fm.Data {
					im.Img = fm.Img
					flagMatrix = append(flagMatrix[:i], flagMatrix[i+1:]...)
					// fmt.Printf("imgMatrix[%d][%d] good\n", x, y)
					skip = false
					break
				}
			}
			if skip {
				// fmt.Printf("imgMatrix[%d][%d] extBlock\n", x, y)
				im.Img = extBlock
				diffMatrix[diffNum] = &ImageMatrixData{
					X: x, Y: y,
				}
				diffNum++
			}
		}
	}

	log.Printf("[+] Hash Diff Num : %d\n", diffNum)

	// Merge Image
	log.Printf("[+] Merge Image\n")
	err := mergeImage(imgMatrix, output, sWidth, sHeight, m, n)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
