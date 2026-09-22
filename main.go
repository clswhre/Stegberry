package main

import (
	_ "embed"
	"fmt"
	"image"
	"os"

	flag "github.com/spf13/pflag"
)

type Config struct {
	Input   string
	Message string
	Output  string
	Decode  bool
}

func main() {
	config := parseArgs()

	printHead()

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
		printStyle(colorGreen, "[ * ] Режим - декодування\n")
	} else {
		printStyle(colorMagenta, "[ + ] Режим - кодування\n")
	}

	if !config.Decode {
		hideInImage(img, config.Message, config.Output)
		os.Exit(0)
	}

	result := extractFromImage(img)
	fmt.Printf("Декодоване повідомлення: ")
	printStyle(colorBlue, result)
	fmt.Println()

}

func parseArgs() Config {
	var config Config

	flag.StringVarP(&config.Input, "input", "i", "", "/path/to/image.png (Обов'язково) | string")
	flag.StringVarP(&config.Message, "message", "m", "", "Повідомлення для приховування (Обов'язково для кодування) | string")
	flag.StringVarP(&config.Output, "output", "o", "output.png", "/path/to/image_output.png | string")
	flag.BoolVarP(&config.Decode, "decode", "d", false, "Режим декодування | bool")

	flag.Parse()

	if config.Input == "" {
		fmt.Fprintln(os.Stderr, "Помилка: прапорець `-input` обов'язковий!")
		flag.Usage()
		os.Exit(1)
	}

	if !config.Decode && config.Message == "" {
		fmt.Println("Помилка: прапорець `-message` обов'язковий для кодування!")
		flag.Usage()
		os.Exit(1)
	}

	return config
}
