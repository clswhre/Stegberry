package main

import (
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

func hideInImage(input image.Image, message string, outputPath string) {
	bounds := input.Bounds()
	x := bounds.Dx()
	y := bounds.Dy()
	capacity := x * y * 3

	msgBytes := []byte(message + "\x00")
	if len(msgBytes)*8 > capacity {
		fmt.Printf("Помилка: передоз. Ємніст: %d бітів, повідомлення: %d бітів\n", capacity, len(msgBytes)*8)
		os.Exit(1)
	}

	fmt.Printf("X >> %d px\nY >> %d px\nC(лише для непрозорих зображень) >> %d bits\n", x, y, capacity)

	var bits []uint8
	for _, b := range msgBytes {
		for i := 7; i >= 0; i-- {
			bits = append(bits, (b>>i)&1)
		}
	}

	outImg := image.NewRGBA(bounds)
	bitIdx := 0

	for py := bounds.Min.Y; py < bounds.Max.Y; py++ {
		for px := bounds.Min.X; px < bounds.Max.X; px++ {

			c := color.RGBAModel.Convert(input.At(px, py)).(color.RGBA)
			if c.A != 0 {
				if bitIdx < len(bits) {
					c.R = (c.R & 0xFE) | bits[bitIdx]
					bitIdx++
				}
				if bitIdx < len(bits) {
					c.G = (c.G & 0xFE) | bits[bitIdx]
					bitIdx++
				}
				if bitIdx < len(bits) {
					c.B = (c.B & 0xFE) | bits[bitIdx]
					bitIdx++
				}
			}

			outImg.SetRGBA(px, py, c)
		}
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("Не вдалось створити вихідний файл: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	err = png.Encode(outFile, outImg)
	if err != nil {
		fmt.Printf("Не вдалось зберегти PNG: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Вихідне зображення збережено в %s\n", outputPath)
}

func extractFromImage(input image.Image) string {
	bounds := input.Bounds()
	var messageBytes []byte
	var currentByte uint8
	bitCount := 0

	for py := bounds.Min.Y; py < bounds.Max.Y; py++ {
		for px := bounds.Min.X; px < bounds.Max.X; px++ {
			c := color.RGBAModel.Convert(input.At(px, py)).(color.RGBA)
			if c.A == 0 {
				continue
			}
			channels := []uint8{c.R, c.G, c.B}

			for _, ch := range channels {
				// 1. (ch & 1) витягує найменш значущий біт (LSB) з каналу
				// 2. (currentByte << 1) зсуває вже зібрані біти вліво
				// 3. | додає новий біт у кінець
				currentByte = (currentByte << 1) | (ch & 1)
				bitCount++

				if bitCount == 8 {
					if currentByte == 0 {
						return string(messageBytes)
					}

					messageBytes = append(messageBytes, currentByte)
					currentByte = 0
					bitCount = 0
				}
			}
		}
	}

	return string(messageBytes)
}
