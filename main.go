package main

import (
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"time"

	flag "github.com/spf13/pflag"
)

const (
	colorBlue  = "\033[34m"
	colorReset = "\033[0m"
)

type Config struct {
	Input   string
	Message string
	Output  string
	Decode  bool
}

//go:embed logo.txt
var asciiLogo string

func main() {
	config := parseArgs()

	fmt.Printf("%s%s%s", colorBlue, asciiLogo, colorReset)

	file, err := os.Open(config.Input)
	if err != nil {
		fmt.Printf("Не вдалось відкрити файл: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Printf("Не вдалось декодувати файл: %v\n", err)
		os.Exit(1)
	}

	if config.Decode {
		result := extractFromImage(img)
		fmt.Println(result)
	} else {
		hideInImage(img, config.Message, config.Output)
	}

	time.Sleep(400 * time.Millisecond)

}

func parseArgs() Config {
	var cfg Config

	flag.StringVarP(&cfg.Input, "input", "i", "", "/path/to/image.png (Обов'язково) | string")
	flag.StringVarP(&cfg.Message, "message", "m", "", "Повідомлення для приховування (Обов'язково для кодування) | string")
	flag.StringVarP(&cfg.Output, "output", "o", "output.png", "/path/to/image_output.png | string")
	flag.BoolVarP(&cfg.Decode, "decode", "d", false, "Режим декодування | bool")

	flag.Parse()

	if cfg.Input == "" {
		fmt.Println("Помилка: прапорець `-input` обов'язковий!")
		flag.Usage()
		os.Exit(1)
	}

	if !cfg.Decode && cfg.Message == "" {
		fmt.Println("Помилка: прапорець `-message` обов'язковий для кодування!")
		flag.Usage()
		os.Exit(1)
	}

	return cfg
}

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

	fmt.Printf("X >> %d px\nY >> %d px\nC >> %d bits\n", x, y, capacity)

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

	fmt.Printf("Output file successfully saved to: %s\n", outputPath)
}

func extractFromImage(input image.Image) string {
	bounds := input.Bounds()
	var messageBytes []byte
	var currentByte uint8
	bitCount := 0

	for py := bounds.Min.Y; py < bounds.Max.Y; py++ {
		for px := bounds.Min.X; px < bounds.Max.X; px++ {
			// Конвертуємо піксель у стандартний 8-бітний RGBA
			c := color.RGBAModel.Convert(input.At(px, py)).(color.RGBA)

			// Перевіряємо канали R, G, B
			channels := []uint8{c.R, c.G, c.B}

			for _, ch := range channels {
				// 1. (ch & 1) витягує найменш значущий біт (LSB) з каналу
				// 2. (currentByte << 1) зсуває вже зібрані біти вліво
				// 3. | додає новий біт у кінець
				currentByte = (currentByte << 1) | (ch & 1)
				bitCount++

				// Коли зібрали повний байт (8 бітів)
				if bitCount == 8 {
					// Якщо байт дорівнює 0 (наш маркер кінця \x00), зупиняємось
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

	// Повертаємо те, що вдалося зібрати, якщо нульовий байт не знайдено
	return string(messageBytes)
}
