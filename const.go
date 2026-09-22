package main

import (
	_ "embed"
)

//go:embed logo.txt
var asciiLogo string

const (
	version = "1.1.0"

	styleBold      = "\033[1m"
	styleItalic    = "\033[3m"
	styleUnderline = "\033[4m"

	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
)
