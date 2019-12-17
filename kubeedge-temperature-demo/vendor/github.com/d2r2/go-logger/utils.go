package logger

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type logFile struct {
	FileInfo os.FileInfo
	Index    int
}

type sortLogFiles struct {
	Items []logFile
}

func (sf *sortLogFiles) Len() int {
	return len(sf.Items)
}

func (sf *sortLogFiles) Less(i, j int) bool {
	return sf.Items[j].Index < sf.Items[i].Index
}

func (sf *sortLogFiles) Swap(i, j int) {
	item := sf.Items[i]
	sf.Items[i] = sf.Items[j]
	sf.Items[j] = item
}

func findStringSubmatchIndexes(r *regexp.Regexp, s string) map[string][2]int {
	captures := make(map[string][2]int)
	ind := r.FindStringSubmatchIndex(s)
	names := r.SubexpNames()
	for i, name := range names {
		if name != "" && i < len(ind)/2 {
			if ind[i*2] != -1 && ind[i*2+1] != -1 {
				captures[name] = [2]int{ind[i*2], ind[i*2+1]}
			}
		}
	}
	return captures
}

func extractIndex(item os.FileInfo) int {
	r := regexp.MustCompile(`.+\.log(\.(?P<index>\d+))?`)
	fileName := path.Base(item.Name())
	m := findStringSubmatchIndexes(r, fileName)
	if v, ok := m["index"]; ok {
		i, _ := strconv.Atoi(fileName[v[0]:v[1]])
		return i
	} else {
		return 0
	}
}

const (
	nocolor = 0
	red     = 31
	green   = 32
	yellow  = 33
	blue    = 34
	gray    = 37
)

type IndentKind int

const (
	LeftIndent = iota
	CenterIndent
	RightIndent
)

func cutOrIndentText(text string, length int, indent IndentKind) string {
	if length < 0 {
		return text
	} else if len(text) > length {
		text = text[:length]
	} else {
		switch indent {
		case LeftIndent:
			text = text + strings.Repeat(" ", length-len(text))
		case RightIndent:
			text = strings.Repeat(" ", length-len(text)) + text
		case CenterIndent:
			text = strings.Repeat(" ", (length-len(text))/2) + text +
				strings.Repeat(" ", length-len(text)-(length-len(text))/2)

		}
	}
	return text
}

func metaFmtStr(colored bool, level LogLevel, options FormatOptions, appName string,
	packageName string, message string, format string) string {
	var colorPfx, colorSfx string
	if colored {
		var levelColor int
		switch level {
		case DebugLevel:
			levelColor = gray
		case InfoLevel:
			levelColor = blue
		case NotifyLevel, WarnLevel:
			levelColor = yellow
		case ErrorLevel, PanicLevel, FatalLevel:
			levelColor = red
		default:
			levelColor = nocolor
		}
		colorPfx = "\x1b[" + strconv.Itoa(levelColor) + "m"
		colorSfx = "\x1b[0m"
	}
	arg1 := time.Now().Format(options.TimeFormat)
	arg2 := appName
	arg3 := cutOrIndentText(packageName, options.PackageLength, RightIndent)
	lvlStr := options.GetLevelStr(level)
	lvlLen := len([]rune(lvlStr))
	arg4 := colorPfx + cutOrIndentText(strings.ToUpper(lvlStr), lvlLen, LeftIndent) + colorSfx
	arg5 := message
	out := fmt.Sprintf(format, arg1, arg2, arg3, arg4, arg5)
	return out
}

func getApplicationName() string {
	appName := os.Args[0]
	return appName
}
