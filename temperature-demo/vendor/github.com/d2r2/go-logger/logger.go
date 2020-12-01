package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"
)

type LogLevel int

const (
	FatalLevel LogLevel = iota
	PanicLevel
	ErrorLevel
	WarnLevel
	NotifyLevel
	InfoLevel
	DebugLevel
)

func (v LogLevel) String() string {
	switch v {
	case FatalLevel:
		return "Fatal"
	case PanicLevel:
		return "Panic"
	case ErrorLevel:
		return "Error"
	case WarnLevel:
		return "Warning"
	case NotifyLevel:
		return "Notice"
	case InfoLevel:
		return "Information"
	case DebugLevel:
		return "Debug"
	default:
		return "<undef>"
	}
}

func (v LogLevel) LongStr() string {
	return v.String()
}

func (v LogLevel) ShortStr() string {
	switch v {
	case FatalLevel:
		return "Fatal"
	case PanicLevel:
		return "Panic"
	case ErrorLevel:
		return "Error"
	case WarnLevel:
		return "Warn"
	case NotifyLevel:
		return "Notice"
	case InfoLevel:
		return "Info"
	case DebugLevel:
		return "Debug"
	default:
		return "<undef>"
	}
}

type Log struct {
	log     *log.Logger
	colored bool
	level   LogLevel
}

func NewLog(log *log.Logger, colored bool, level LogLevel) *Log {
	v := &Log{log: log, colored: colored, level: level}
	return v
}

type Logger struct {
	sync.RWMutex
	logs           []*Log
	packages       []*Package
	options        FormatOptions
	logFile        *File
	rotateMaxSize  int64
	rotateMaxCount int
	enableSyslog   bool
}

func NewLogger() *Logger {
	stdout := NewLog(log.New(os.Stdout, "", 0), true, DebugLevel)
	logs := []*Log{stdout}
	options := FormatOptions{TimeFormat: "2006-01-02T15:04:05.000", LevelLength: LevelShort, PackageLength: 8}
	l := &Logger{
		logs:           logs,
		options:        options,
		rotateMaxSize:  1024 * 1024 * 512,
		rotateMaxCount: 3,
	}
	return l
}

func (v *Logger) Close() error {
	v.Lock()
	defer v.Unlock()

	for _, pack := range v.packages {
		pack.Close()
	}
	v.packages = nil

	if v.logFile != nil {
		v.logFile.Close()
	}
	return nil
}

func (v *Logger) SetRotateParams(rotateMaxSize int64, rotateMaxCount int) {
	v.Lock()
	defer v.Unlock()
	v.rotateMaxSize = rotateMaxSize
	v.rotateMaxCount = rotateMaxCount
}

func (v *Logger) GetRotateMaxSize() int64 {
	v.Lock()
	defer v.Unlock()
	return v.rotateMaxSize
}

func (v *Logger) GetRotateMaxCount() int {
	v.Lock()
	defer v.Unlock()
	return v.rotateMaxCount
}

/*
func (v *Logger) SetApplicationName(appName string) {
	v.Lock()
	defer v.Unlock()
	v.appName = appName
}

func (v *Logger) GetApplicationName() string {
	v.RLock()
	defer v.RUnlock()
	return v.appName
}
*/

func (v *Logger) EnableSyslog(enable bool) {
	v.Lock()
	defer v.Unlock()
	v.enableSyslog = enable
}

func (v *Logger) GetSyslogEnabled() bool {
	v.RLock()
	defer v.RUnlock()
	return v.enableSyslog
}

func (v *Logger) SetFormatOptions(options FormatOptions) {
	v.Lock()
	defer v.Unlock()
	v.options = options
}

func (v *Logger) GetFormatOptions() FormatOptions {
	v.RLock()
	defer v.RUnlock()
	return v.options
}

func (v *Logger) SetLogFileName(logFilePath string) error {
	if path.Ext(logFilePath) == "" {
		logFilePath += ".log"
	}
	fp, err := filepath.Abs(logFilePath)
	if err != nil {
		return err
	}
	v.Lock()
	defer v.Unlock()
	lf := &File{Path: fp}
	v.logFile = lf
	return nil
}

func (v *Logger) GetLogFileInfo() *File {
	v.RLock()
	defer v.RUnlock()
	return v.logFile
}

func (v *Logger) NewPackageLogger(packageName string, level LogLevel) PackageLog {
	v.Lock()
	defer v.Unlock()
	p := &Package{parent: v, packageName: packageName, level: level}
	v.packages = append(v.packages, p)
	return p
}

func (v *Logger) ChangePackageLogLevel(packageName string, level LogLevel) error {
	var p *Package
	for _, item := range v.packages {
		if item.packageName == packageName {
			p = item
			break
		}
	}
	if p != nil {
		p.SetLogLevel(level)
	} else {
		err := fmt.Errorf("Package log %q is not found", packageName)
		return err
	}
	return nil
}

func (v *Logger) AddCustomLog(writer io.Writer, colored bool, level LogLevel) {
	v.Lock()
	defer v.Unlock()
	log := NewLog(log.New(writer, "", 0), colored, level)
	v.logs = append(v.logs, log)
}

func (v *Logger) getLogs() []*Log {
	v.RLock()
	defer v.RUnlock()
	lst := []*Log{}
	for _, item := range v.logs {
		lst = append(lst, item)
	}
	return lst
}

var (
	globalLock sync.RWMutex
	lgr        *Logger
)

func SetFormatOptions(format FormatOptions) {
	globalLock.RLock()
	defer globalLock.RUnlock()
	lgr.SetFormatOptions(format)
}

func SetRotateParams(rotateMaxSize int64, rotateMaxCount int) {
	globalLock.RLock()
	defer globalLock.RUnlock()
	lgr.SetRotateParams(rotateMaxSize, rotateMaxCount)
}

func NewPackageLogger(module string, level LogLevel) PackageLog {
	globalLock.RLock()
	defer globalLock.RUnlock()
	return lgr.NewPackageLogger(module, level)
}

func ChangePackageLogLevel(packageName string, level LogLevel) error {
	globalLock.RLock()
	defer globalLock.RUnlock()
	return lgr.ChangePackageLogLevel(packageName, level)
}

func SetLogFileName(logFilePath string) error {
	globalLock.RLock()
	defer globalLock.RUnlock()
	return lgr.SetLogFileName(logFilePath)
}

/*
func SetApplicationName(appName string) {
	globalLock.RLock()
	defer globalLock.RUnlock()
	lgr.SetApplicationName(appName)
}
*/

func EnableSyslog(enable bool) {
	globalLock.RLock()
	defer globalLock.RUnlock()
	lgr.EnableSyslog(enable)
}

func AddCustomLog(writer io.Writer, colored bool, level LogLevel) {
	globalLock.RLock()
	defer globalLock.RUnlock()
	lgr.AddCustomLog(writer, colored, level)
}

func FinalizeLogger() error {
	var err error
	if lgr != nil {
		err = lgr.Close()
	}
	globalLock.Lock()
	defer globalLock.Unlock()
	lgr = nil
	return err
}

func init() {
	lgr = NewLogger()
}
