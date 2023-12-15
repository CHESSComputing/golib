package utils

import "fmt"

// Color prints given string in color based on ANSI escape codes, see
// http://www.wikiwand.com/en/ANSI_escape_code#/Colors
func Color(col, text string) string {
	return BOLD + "\x1b[" + col + text + PLAIN
}

// ColorURL returns colored string of given url
func ColorURL(rurl string) string {
	return Color(BLUE, rurl)
}

// Error prints Server error message with given arguments
func Error(args ...interface{}) {
	fmt.Println(Color(RED, "Server ERROR"), args)
}

// Warning prints Server error message with given arguments
func Warning(args ...interface{}) {
	fmt.Println(Color(BROWN, "Server WARNING"), args)
}

// BLACK color
const BLACK = "0;30m"

// RED color
const RED = "0;31m"

// GREEN color
const GREEN = "0;32m"

// BROWN color
const BROWN = "0;33m"

// BLUE color
const BLUE = "0;34m"

// PURPLE color
const PURPLE = "0;35m"

// CYAN color
const CYAN = "0;36m"

// LIGHT_PURPLE color
const LIGHT_PURPLE = "1;35m"

// LIGHT_CYAN color
const LIGHT_CYAN = "1;36m"

// BOLD type
const BOLD = "\x1b[1m"

// PLAIN type
const PLAIN = "\x1b[0m"
