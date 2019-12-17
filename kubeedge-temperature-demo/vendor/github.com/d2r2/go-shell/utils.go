package shell

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

// IsLinuxMacOSFreeBSD determine that running OS belong to Linux, BSD or macOS.
func IsLinuxMacOSFreeBSD() bool {
	return runtime.GOOS == "linux" || runtime.GOOS == "darwin" ||
		runtime.GOOS == "freebsd"
}

// CheckRunAsRoot verify that current context
// is running under root privileges.
func CheckRunAsRoot() bool {
	uid := os.Geteuid()
	if uid == 0 {
		return true
	}
	return false
}

// GetFreeSpace use syscall to find free space for path specified.
func GetFreeSpace(path string) (uint64, error) {
	const errMsg = "can't detect free space available on the system"
	var space uint64
	if IsLinuxMacOSFreeBSD() {
		var stat syscall.Statfs_t
		err := syscall.Statfs(path, &stat)
		if err != nil {
			return 0, err
		}
		// stat.Bavail type is not the same on Linux and FreeBSD, so
		// check that it valid, than cast it to UINT64.
		if stat.Bavail < 0 {
			return 0, errors.New(errMsg)
		}
		// Available blocks * size per block = available space in bytes.
		space = uint64(stat.Bavail) * uint64(stat.Bsize)
	} else {
		return 0, errors.New(errMsg)
	}
	/* else {
		h := syscall.LoadLibrary("kernel32.dll")
		c := h.MustFindProc("GetDiskFreeSpaceExW")
		var freeBytes int64
		_, _, err := c.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(path))),
			uintptr(unsafe.Pointer(&freeBytes)), nil, nil)
		if err != nil {
			return 0, err
		}
		space = uint64(freeBytes)
	}*/
	return space, nil
}

// CopyFile copy regular file.
// Code taken from: https://opensource.com/article/18/6/copying-files-go
func CopyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

// CloseChannelOnSignals close  channel once signals received.
func CloseChannelOnSignals(kill chan struct{}, quit chan struct{}, signals ...os.Signal) {
	// Set up channel on which to send signal notifications
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, signals...)
	// run gorutine and block until a signal is received
	go func() {
		select {
		case <-c:
			// send signal to threads about pending to close
			log.Println("Signal received, close kill channel")
			close(kill)
		// if quit is not null, it can be used as an exit from gorutine
		case <-quit:
			// exit
		}
	}()
}

// CloseContextOnSignals call cancel method once signals received.
func CloseContextOnSignals(cancel context.CancelFunc, quit chan struct{}, signals ...os.Signal) {
	// Set up channel on which to send signal notifications
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, signals...)
	// run gorutine and block until a signal is received
	go func() {
		select {
		case <-c:
			// send pending signal to threads to close
			log.Println("Signal received, cancel context")
			if cancel != nil {
				cancel()
			}
		// if quit is not null, it can be used as an exit from gorutine
		case <-quit:
			// exit
		}
	}()
}
