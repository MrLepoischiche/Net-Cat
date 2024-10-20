package colors

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"
)

const (
	selectGraphicRendition = "\x1B["

	FgBlack         = 30
	FgRed           = 31
	FgGreen         = 32
	FgYellow        = 33
	FgBlue          = 34
	FgMagenta       = 35
	FgCyan          = 36
	FgWhite         = 37
	FgReset         = 39
	FgBrightBlack   = 90
	FgBrightRed     = 91
	FgBrightGreen   = 92
	FgBrightYellow  = 93
	FgBrightBlue    = 94
	FgBrightMagenta = 95
	FgBrightCyan    = 96
	FgBrightWhite   = 97

	BgBlack         = 40
	BgRed           = 41
	BgGreen         = 42
	BgYellow        = 43
	BgBlue          = 44
	BgMagenta       = 45
	BgCyan          = 46
	BgWhite         = 47
	BgReset         = 49
	BgBrightBlack   = 100
	BgBrightRed     = 101
	BgBrightGreen   = 102
	BgBrightYellow  = 103
	BgBrightBlue    = 104
	BgBrightMagenta = 105
	BgBrightCyan    = 106
	BgBrightWhite   = 107

	ResetFGColorTag = selectGraphicRendition + "39m"
	ResetBGColorTag = selectGraphicRendition + "49m"
	ResetColorsTag  = selectGraphicRendition + "39;49m"
)

var (
	EnglishColors = map[string][3]int{
		"amber":      {255, 191, 0},
		"amethyst":   {153, 102, 204},
		"apricot":    {250, 207, 176},
		"aqua":       {0, 255, 255},
		"argent":     {171, 168, 173},
		"azur":       {0, 127, 255},
		"azure":      {240, 255, 255},
		"banana":     {250, 232, 181},
		"bleu":       {48, 140, 232},
		"blizzard":   {171, 229, 237},
		"blood":      {102, 0, 0},
		"blue":       {0, 0, 127},
		"bonbon":     {250, 66, 158},
		"bronze":     {204, 127, 51},
		"brown":      {165, 42, 42},
		"byzantine":  {189, 51, 163},
		"byzantium":  {112, 41, 99},
		"candy":      {227, 112, 122},
		"cardinal":   {196, 31, 59},
		"carmine":    {150, 0, 23},
		"chocolat":   {122, 64, 0},
		"chocolate":  {210, 105, 31},
		"celeste":    {178, 255, 255},
		"cerise":     {222, 48, 99},
		"champagne":  {247, 232, 207},
		"coffee":     {112, 79, 56},
		"copper":     {184, 115, 51},
		"coral":      {255, 127, 80},
		"cream":      {255, 252, 209},
		"crimson":    {220, 20, 60},
		"egg":        {240, 235, 214},
		"eggshell":   {240, 235, 214},
		"emerald":    {79, 199, 120},
		"fire":       {227, 89, 33},
		"flame":      {227, 89, 33},
		"fuschia":    {255, 0, 255},
		"gold":       {255, 215, 0},
		"grass":      {124, 252, 0},
		"iceberg":    {112, 166, 209},
		"indigo":     {75, 0, 130},
		"ivory":      {255, 255, 240},
		"lava":       {207, 15, 33},
		"lavender":   {230, 230, 250},
		"lilac":      {199, 163, 199},
		"lime":       {0, 255, 0},
		"mahogany":   {191, 64, 0},
		"malachite":  {10, 217, 82},
		"mandarin":   {242, 122, 71},
		"mango":      {252, 191, 3},
		"mauve":      {224, 176, 255},
		"melon":      {255, 186, 173},
		"mint":       {61, 181, 138},
		"mustard":    {255, 219, 89},
		"navy":       {0, 0, 128},
		"nickel":     {115, 115, 115},
		"onyx":       {54, 56, 56},
		"opal":       {168, 194, 189},
		"or":         {212, 176, 56},
		"orange":     {255, 127, 0},
		"orchid":     {217, 112, 214},
		"peach":      {255, 230, 181},
		"pear":       {209, 227, 48},
		"pink":       {255, 192, 203},
		"pistachio":  {148, 196, 115},
		"platinum":   {230, 227, 227},
		"plum":       {143, 69, 133},
		"prune":      {112, 28, 28},
		"pumpkin":    {255, 117, 23},
		"purple":     {97, 0, 128},
		"raspberry":  {227, 10, 92},
		"rose":       {255, 0, 127},
		"ruby":       {224, 18, 94},
		"rust":       {184, 64, 13},
		"salmon":     {250, 128, 114},
		"sand":       {194, 178, 127},
		"sapphire":   {15, 82, 186},
		"scarlet":    {255, 36, 0},
		"sepia":      {112, 66, 20},
		"silver":     {192, 192, 192},
		"skyblue":    {135, 206, 235},
		"strawberry": {250, 79, 84},
		"tangerine":  {242, 133, 0},
		"teal":       {0, 127, 127},
		"tomato":     {255, 99, 71},
		"tourmaline": {135, 161, 168},
		"tumbleweed": {222, 171, 135},
		"turquoise":  {64, 224, 208},
		"vanilla":    {242, 230, 171},
		"violet":     {143, 0, 255},
		"wine":       {115, 46, 56},
		"wood":       {194, 153, 107},
		//"":{},
	}
)

// RGB functions
func NewRandomFGColorRGB() string {
	color, _ := NewFGColorRGB(rand.IntN(256), rand.IntN(256), rand.IntN(256))
	return color
}

func NewRandomBGColorRGB() string {
	color, _ := NewBGColorRGB(rand.IntN(256), rand.IntN(256), rand.IntN(256))
	return color
}

func NewFGColorRGB(r, g, b int) (string, error) {
	if (r < 0 || r > 255) && (g < 0 || g > 255) && (b < 0 || b > 255) {
		return "", errors.New("invalid parameters: numbers must be between 127 and 255")
	}
	return fmt.Sprintf("%s38;2;%d;%d;%dm", selectGraphicRendition, r, g, b), nil
}

func NewBGColorRGB(r, g, b int) (string, error) {
	if (r < 0 || r > 255) && (g < 0 || g > 255) && (b < 0 || b > 255) {
		return "", errors.New("invalid parameters: numbers must be between 127 and 255")
	}
	return fmt.Sprintf("%s48;2;%d;%d;%dm", selectGraphicRendition, r, g, b), nil
}

func RGBValuesToColor(vals [3]int) (string, error) {
	return fmt.Sprintf("%s38;2;%d;%d;%dm", selectGraphicRendition, vals[0], vals[1], vals[2]), nil
}

func SprintfForegroundRGB(str string, r, g, b int) string {
	rgbColor, errRGB := NewFGColorRGB(r, g, b)
	if errRGB != nil {
		return str
	}

	return fmt.Sprintf("%s%s%s", rgbColor, str, ResetFGColorTag)
}

func SprintfBackgroundRGB(str string, r, g, b int) string {
	rgbColor, errRGB := NewBGColorRGB(r, g, b)
	if errRGB != nil {
		return str
	}

	return fmt.Sprintf("%s%s%s", rgbColor, str, ResetBGColorTag)
}

func SprintfRGB(str string, rFG, gFG, bFG, rBG, gBG, bBG int) string {
	rgbFGColor, errRGB := NewFGColorRGB(rFG, gFG, bFG)
	if errRGB != nil {
		return str
	}
	rgbBGColor, errRGB := NewBGColorRGB(rBG, gBG, bBG)
	if errRGB != nil {
		return str
	}

	return fmt.Sprintf("%s%s%s%s", rgbFGColor, rgbBGColor, str, ResetColorsTag)
}

// ANSI functions
func SprintfANSI(str string, fg, bg int) string {
	if ((fg < 30 || fg > 37) && fg != 39) && (fg < 90 || fg > 97) {
		return str
	}
	if ((bg < 40 || bg > 47) && bg != 49) && (bg < 100 || bg > 107) {
		return str
	}

	return fmt.Sprintf("%s%d;%dm%s%s", selectGraphicRendition, fg, bg, str, ResetColorsTag)
}

func NewFGColorANSI(color string) (string, error) {
	if color == "" {
		return "", errors.New("empty parameter")
	}

	switch strings.ToLower(color) {
	case "black":
		return selectGraphicRendition + strconv.Itoa(FgBlack) + "m", nil
	case "red":
		return selectGraphicRendition + strconv.Itoa(FgRed) + "m", nil
	case "green":
		return selectGraphicRendition + strconv.Itoa(FgGreen) + "m", nil
	case "yellow":
		return selectGraphicRendition + strconv.Itoa(FgYellow) + "m", nil
	case "blue":
		return selectGraphicRendition + strconv.Itoa(FgBlue) + "m", nil
	case "magenta":
		return selectGraphicRendition + strconv.Itoa(FgMagenta) + "m", nil
	case "cyan":
		return selectGraphicRendition + strconv.Itoa(FgCyan) + "m", nil
	case "white":
		return selectGraphicRendition + strconv.Itoa(FgWhite) + "m", nil
	case "bright black", "bright-black", "brightblack":
		return selectGraphicRendition + strconv.Itoa(FgBrightBlack) + "m", nil
	case "bright red", "bright-red", "brightred":
		return selectGraphicRendition + strconv.Itoa(FgBrightRed) + "m", nil
	case "bright green", "bright-green", "brightgreen":
		return selectGraphicRendition + strconv.Itoa(FgBrightGreen) + "m", nil
	case "bright yellow", "bright-yellow", "brightyellow":
		return selectGraphicRendition + strconv.Itoa(FgBrightYellow) + "m", nil
	case "bright blue", "bright-blue", "brightblue":
		return selectGraphicRendition + strconv.Itoa(FgBrightBlue) + "m", nil
	case "bright magenta", "bright-magenta", "brightmagenta":
		return selectGraphicRendition + strconv.Itoa(FgBrightMagenta) + "m", nil
	case "bright cyan", "bright-cyan", "brightcyan":
		return selectGraphicRendition + strconv.Itoa(FgBrightCyan) + "m", nil
	case "bright white", "bright-white", "brightwhite":
		return selectGraphicRendition + strconv.Itoa(FgBrightWhite) + "m", nil
	default:
		return "", errors.New("invalid parameter")
	}
}
