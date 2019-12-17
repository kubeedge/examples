package logger

import "os"

type LevelLength int

const (
	LevelShort LevelLength = iota
	LevelLong
)

type FormatOptions struct {
	TimeFormat    string
	PackageLength int
	LevelLength   LevelLength
}

func FormatMessage(options FormatOptions, level LogLevel, packageName, msg string, colored bool) string {
	appName := os.Args[0]
	out := metaFmtStr(colored, level, options, appName,
		packageName, msg, "%[1]s [%[3]s] %[4]s  %[5]s")
	return out
}

func (options FormatOptions) GetLevelStr(level LogLevel) string {
	if options.LevelLength == LevelLong {
		return level.LongStr()
	} else {
		return level.ShortStr()
	}
}
