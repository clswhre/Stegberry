package main

import (
	_ "embed"
	"fmt"
)

func printHead() {
	printStyle(colorRed, asciiLogo)
	fmt.Printf("%20s%sv%s%s\n", "", styleUnderline, version, colorReset)
	printStyle(colorRed, "-------------------------------------------------------\n")
}

func printStyle(style string, text string) {
	fmt.Printf("%s%s%s", style, text, colorReset)
}
