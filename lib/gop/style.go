package gop

import (
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// Style type
type Style struct {
	Set   string
	Unset string
}

var (
	// Bold style
	Bold = Style{"\x1b[1m", "\x1b[21m"}
	// Faint style
	Faint = Style{"\x1b[2m", "\x1b[22m"}
	// Italic style
	Italic = Style{"\x1b[3m", "\x1b[23m"}
	// Underline style
	Underline = Style{"\x1b[4m", "\x1b[24m"}

	// Black color
	Black = Style{"\x1b[30m", "\x1b[39m"}
	// Red color
	Red = Style{"\x1b[31m", "\x1b[39m"}
	// Green color
	Green = Style{"\x1b[32m", "\x1b[39m"}
	// Yellow color
	Yellow = Style{"\x1b[33m", "\x1b[39m"}
	// Blue color
	Blue = Style{"\x1b[34m", "\x1b[39m"}
	// Magenta color
	Magenta = Style{"\x1b[35m", "\x1b[39m"}
	// Cyan color
	Cyan = Style{"\x1b[36m", "\x1b[39m"}
	// White color
	White = Style{"\x1b[37m", "\x1b[39m"}

	// BgBlack color
	BgBlack = Style{"\x1b[40m", "\x1b[49m"}
	// BgRed color
	BgRed = Style{"\x1b[41m", "\x1b[49m"}
	// BgGreen color
	BgGreen = Style{"\x1b[42m", "\x1b[49m"}
	// BgYellow color
	BgYellow = Style{"\x1b[43m", "\x1b[49m"}
	// BgBlue color
	BgBlue = Style{"\x1b[44m", "\x1b[49m"}
	// BgMagenta color
	BgMagenta = Style{"\x1b[45m", "\x1b[49m"}
	// BgCyan color
	BgCyan = Style{"\x1b[46m", "\x1b[49m"}
	// BgWhite color
	BgWhite = Style{"\x1b[47m", "\x1b[49m"}

	// None type
	None = Style{}
)

var styleMap = map[string]Style{
	Bold.Set:      Bold,
	Faint.Set:     Faint,
	Italic.Set:    Italic,
	Underline.Set: Underline,

	Black.Set:   Black,
	Red.Set:     Red,
	Green.Set:   Green,
	Yellow.Set:  Yellow,
	Blue.Set:    Blue,
	Magenta.Set: Magenta,
	Cyan.Set:    Cyan,
	White.Set:   White,

	BgBlack.Set:   BgBlack,
	BgRed.Set:     BgRed,
	BgGreen.Set:   BgGreen,
	BgYellow.Set:  BgYellow,
	BgBlue.Set:    BgBlue,
	BgMagenta.Set: BgMagenta,
	BgCyan.Set:    BgCyan,
	BgWhite.Set:   BgWhite,
}

var regNewline = regexp.MustCompile(`\r?\n`)

// S is the shortcut for Stylize
func S(str string, styles ...Style) string {
	return Stylize(str, styles)
}

// Stylize string
func Stylize(str string, styles []Style) string {
	for _, s := range styles {
		str = stylize(s, str)
	}
	return str
}

func stylize(s Style, str string) string {
	if NoStyle || s == None {
		return str
	}

	newline := regNewline.FindString(str)

	lines := regNewline.Split(str, -1)
	out := []string{}

	for _, l := range lines {
		out = append(out, s.Set+l+s.Unset)
	}

	return strings.Join(out, newline)
}

// NoStyle respects https://no-color.org/ and "tput colors"
var NoStyle = func() bool {
	_, noColor := os.LookupEnv("NO_COLOR")

	b, _ := exec.Command("tput", "colors").CombinedOutput()
	n, _ := strconv.ParseInt(strings.TrimSpace(string(b)), 10, 32)
	return noColor || n == 0
}()

// RegANSI token
var RegANSI = regexp.MustCompile(`\x1b\[(\d+)m`)

// StripANSI tokens
func StripANSI(str string) string {
	return RegANSI.ReplaceAllString(str, "")
}

// VisualizeANSI tokens
func VisualizeANSI(str string) string {
	return RegANSI.ReplaceAllString(str, "<$1>")
}

// FixNestedStyle like
//     <red>1<blue>2<cyan>3</cyan>4</blue>5</red>
// into
//     <red>1</red><blue>2</blue><cyan>3</cyan><blue>4</blue><red>5</red>
func FixNestedStyle(s string) string {
	out := ""
	stack := []string{}
	i := 0
	l := 0
	r := 0

	for i < len(s) {
		loc := RegANSI.FindStringIndex(s[i:])
		if loc == nil {
			break
		}

		l, r = i+loc[0], i+loc[1]
		token := s[l:r]

		out += s[i:l]

		if len(stack) == 0 {
			stack = append(stack, token)
			out += token
		} else if token == styleMap[stack[len(stack)-1]].Unset {
			out += token
			stack = stack[:len(stack)-1]
			if len(stack) > 0 {
				out += stack[len(stack)-1]
			}
		} else {
			out += styleMap[stack[len(stack)-1]].Unset
			stack = append(stack, token)
			out += token
		}

		i = r
	}

	return out + s[i:]
}
