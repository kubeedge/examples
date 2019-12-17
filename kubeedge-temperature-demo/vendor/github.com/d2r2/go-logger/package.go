package logger

import (
	"fmt"
	"log/syslog"
	"os"
	"sync"

	"github.com/davecgh/go-spew/spew"
)

type PackageLog interface {
	Printf(level LogLevel, format string, args ...interface{})
	Print(level LogLevel, args ...interface{})
	Debugf(format string, args ...interface{})
	Debug(args ...interface{})
	Infof(format string, args ...interface{})
	Info(args ...interface{})
	Notifyf(format string, args ...interface{})
	Notify(args ...interface{})
	Warningf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Warning(args ...interface{})
	Warn(args ...interface{})
	Errorf(format string, args ...interface{})
	Error(args ...interface{})
	Panicf(format string, args ...interface{})
	Panic(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatal(args ...interface{})
}

type Package struct {
	sync.RWMutex
	parent      *Logger
	packageName string
	level       LogLevel
	syslog      *syslog.Writer
}

// Static cast to verify that object implement interface.
var _ PackageLog = &Package{}

func (v *Package) Close() error {
	v.Lock()
	defer v.Unlock()
	if v.syslog != nil {
		err := v.syslog.Close()
		v.syslog = nil
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *Package) SetLogLevel(level LogLevel) {
	v.Lock()
	defer v.Unlock()
	v.level = level
}

func (v *Package) GetLogLevel() LogLevel {
	v.RLock()
	defer v.RUnlock()
	return v.level
}

func (v *Package) getSyslog(level LogLevel, options FormatOptions,
	appName string) (*syslog.Writer, error) {
	v.Lock()
	defer v.Unlock()
	if v.syslog == nil {
		tag := metaFmtStr(false, level, options, appName,
			v.packageName, "", "%[2]s-%[3]s")
		sl, err := syslog.New(syslog.LOG_DEBUG, tag)
		if err != nil {
			err = spew.Errorf("Failed to connect to syslog: %v\n", err)
			return nil, err
		}
		v.syslog = sl
	}
	return v.syslog, nil
}

func (v *Package) writeToSyslog(options FormatOptions,
	level LogLevel, appName string, msg string) error {

	sl, err := v.getSyslog(level, options, appName)
	if err != nil {
		return err
	}
	switch level {
	case DebugLevel:
		return sl.Debug(msg)
	case InfoLevel:
		return sl.Info(msg)
	case WarnLevel:
		return sl.Warning(msg)
	case ErrorLevel:
		return sl.Err(msg)
	case PanicLevel:
		return sl.Crit(msg)
	case FatalLevel:
		return sl.Emerg(msg)
	default:
		return sl.Debug(msg)
	}
}

type printLog func(log *Log, msg interface{})
type getMessage func(colored bool) interface{}

func printLogs(logs []*Log, level LogLevel, prnt printLog, getMsg getMessage) {
	// Console and custom logs output
	for _, log := range logs {
		if log.level >= level {
			prnt(log, getMsg(log.colored))
		}
	}
}

func (v *Package) print(level LogLevel, msg string) {
	lvl := v.GetLogLevel()
	if lvl >= level {
		appName := getApplicationName()
		logs := v.parent.getLogs()
		options := v.parent.GetFormatOptions()
		out1 := FormatMessage(options, level, v.packageName, msg, false)
		// File output
		if lf := v.parent.GetLogFileInfo(); lf != nil {
			rotateMaxSize := v.parent.GetRotateMaxSize()
			rotateMaxCount := v.parent.GetRotateMaxCount()
			if err := lf.writeToFile(out1, rotateMaxSize, rotateMaxCount); err != nil {
				err = spew.Errorf("Failed to report syslog message %q: %v\n", out1, err)
				printLogs(logs, FatalLevel,
					func(log *Log, msg interface{}) {
						log.log.Fatal(msg)
					},
					func(colored bool) interface{} {
						return err
					})
			}
		}
		// Syslog output
		if v.parent.GetSyslogEnabled() {
			if err := v.writeToSyslog(options, level, appName, msg); err != nil {
				err = spew.Errorf("Failed to report syslog message %q: %v\n", msg, err)
				printLogs(logs, FatalLevel,
					func(log *Log, msg interface{}) {
						log.log.Fatal(msg)
					},
					func(colored bool) interface{} {
						return err
					})
			}
		}
		// Console and custom logs output
		outColored1 := FormatMessage(options, level, v.packageName, msg, true)
		printLogs(logs, level,
			func(log *Log, msg interface{}) {
				log.log.Print(msg)
			},
			func(colored bool) interface{} {
				if colored {
					return outColored1 + fmt.Sprintln()
				} else {
					return out1 + fmt.Sprintln()
				}
			})
		// Check critical events
		if level == PanicLevel {
			panic(out1)
		} else if level == FatalLevel {
			os.Exit(1)
		}
	}
}

func (v *Package) Printf(level LogLevel, format string, args ...interface{}) {
	lvl := v.GetLogLevel()
	if lvl >= level {
		msg := spew.Sprintf(format, args...)
		v.print(level, msg)
	}
}

func (v *Package) Print(level LogLevel, args ...interface{}) {
	lvl := v.GetLogLevel()
	if lvl >= level {
		msg := fmt.Sprint(args...)
		v.print(level, msg)
	}
}

func (v *Package) Debugf(format string, args ...interface{}) {
	v.Printf(DebugLevel, format, args...)
}

func (v *Package) Debug(args ...interface{}) {
	v.Print(DebugLevel, args...)
}

func (v *Package) Infof(format string, args ...interface{}) {
	v.Printf(InfoLevel, format, args...)
}

func (v *Package) Info(args ...interface{}) {
	v.Print(InfoLevel, args...)
}

func (v *Package) Notifyf(format string, args ...interface{}) {
	v.Printf(NotifyLevel, format, args...)
}

func (v *Package) Notify(args ...interface{}) {
	v.Print(NotifyLevel, args...)
}

func (v *Package) Warningf(format string, args ...interface{}) {
	v.Printf(WarnLevel, format, args...)
}

func (v *Package) Warnf(format string, args ...interface{}) {
	v.Printf(WarnLevel, format, args...)
}

func (v *Package) Warning(args ...interface{}) {
	v.Print(WarnLevel, args...)
}

func (v *Package) Warn(args ...interface{}) {
	v.Print(WarnLevel, args...)
}

func (v *Package) Errorf(format string, args ...interface{}) {
	v.Printf(ErrorLevel, format, args...)
}

func (v *Package) Error(args ...interface{}) {
	v.Print(ErrorLevel, args...)
}

func (v *Package) Panicf(format string, args ...interface{}) {
	v.Printf(PanicLevel, format, args...)
}

func (v *Package) Panic(args ...interface{}) {
	v.Print(PanicLevel, args...)
}

func (v *Package) Fatalf(format string, args ...interface{}) {
	v.Printf(FatalLevel, format, args...)
}

func (v *Package) Fatal(args ...interface{}) {
	v.Print(FatalLevel, args...)
}
