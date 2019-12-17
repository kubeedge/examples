package shell

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"sync/atomic"
	"syscall"
)

// ExitCodeOrError keeps exit code from application termination
// either error if application failed in any stage.
type ExitCodeOrError struct {
	ExitCode int
	Error    error
}

// App struct keep everything regarding external application started process
// including command line, wait channel which tracks process completion
// and exit code ether any exception happened in any stage of
// application start up or completion.
type App struct {
	cmd             *exec.Cmd
	waitCh          chan ExitCodeOrError
	exitCodeOrError atomic.Value
}

// NewApp return new application instance defined by executable name
// and arguments, and ready to start by following Run call.
func NewApp(name string, args ...string) *App {
	cmd := exec.Command(name, args...)

	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	app := &App{cmd: cmd}
	return app
}

// AddEnvironments add environments in the form "key=value".
func (app *App) AddEnvironments(env []string) {
	if app.cmd.Env == nil {
		app.cmd.Env = os.Environ()
	}
	app.cmd.Env = append(app.cmd.Env, env...)
}

// Run start application synchronously with link to the process
// stdout/stderr output, to get output.
// Method doesn't return control until the application
// finishes its execution.
func (app *App) Run(stdOut *bytes.Buffer, stdErr *bytes.Buffer) ExitCodeOrError {
	_, err := app.Start(stdOut, stdErr)
	if err != nil {
		return ExitCodeOrError{0, err}
	}
	/*
		err = syscall.Setpriority(1, app.cmd.Process.Pid, 19)
		if err != nil {
			return ExitCodeOrError{0, err}
		}
	*/
	st := app.Wait()
	return st
}

func (app *App) sendExitCodeOrError(exitCode int, err error) {
	state := &ExitCodeOrError{ExitCode: exitCode, Error: err}
	// log.Printf("Exit status: %+v", state)
	app.exitCodeOrError.Store(state)
	app.waitCh <- *state
}

func readFromIo(read io.Reader, buf *bytes.Buffer) error {
	var b [4096]byte
	for {
		n, err := read.Read(b[:])
		if n > 0 {
			buf.Write(b[:n])
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (app *App) asyncWait(stdOut, stdErr *bytes.Buffer,
	readOut, readErr io.ReadCloser) {
	defer close(app.waitCh)

	var wg sync.WaitGroup
	if readOut != nil {
		func(wg *sync.WaitGroup) {
			wg.Add(1)
			defer wg.Done()
			err2 := readFromIo(readOut, stdOut)
			if err2 != nil {
				app.sendExitCodeOrError(0, err2)
				return
			}
		}(&wg)
	}
	if readErr != nil {
		func(wg *sync.WaitGroup) {
			wg.Add(1)
			defer wg.Done()
			err2 := readFromIo(readErr, stdErr)
			if err2 != nil {
				app.sendExitCodeOrError(0, err2)
				return
			}
		}(&wg)
	}
	wg.Wait()
	err := app.cmd.Wait()
	var exitCode int
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if stat, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				exitCode = stat.ExitStatus()
				// reset error, since exitCode already not equal to zero
				err = nil
			}
		}
	}
	app.sendExitCodeOrError(exitCode, err)
}

// Start run application asynchronously and
// return channel to wait/track exit state and status.
// If application failed to run, error returned,
func (app *App) Start(stdOut *bytes.Buffer,
	stdErr *bytes.Buffer) (chan ExitCodeOrError, error) {
	var readOut io.ReadCloser
	var readErr io.ReadCloser
	var err error
	if stdOut != nil {
		readOut, err = app.cmd.StdoutPipe()
		if err != nil {
			return nil, err
		}
	}
	if stdErr != nil {
		readErr, err = app.cmd.StderrPipe()
		if err != nil {
			return nil, err
		}
	}
	err = app.cmd.Start()
	if err != nil {
		return nil, err
	}
	app.waitCh = make(chan ExitCodeOrError)
	go app.asyncWait(stdOut, stdErr, readOut, readErr)
	return app.waitCh, nil
}

// CheckIsInstalled use Linux utility [which] to find
// that executable installed or not in the system.
func (app *App) CheckIsInstalled() error {
	// Can't use [whereis], because it doesn't return correct exit code
	// based on search results. Can use [type], as an option.
	whApp := NewApp("which", app.cmd.Path)
	st := whApp.Run(nil, nil)
	if st.Error != nil {
		return st.Error
	}
	if st.ExitCode != 0 {
		return fmt.Errorf("App \"%s\" does not exist", app.cmd.Path)
	}
	return nil
}

// ExitCodeOrError return exit status once application has been finished.
func (app *App) ExitCodeOrError() *ExitCodeOrError {
	ref := app.exitCodeOrError.Load()
	return ref.(*ExitCodeOrError)
}

// Wait switch from asynchronous mode to synchronous
// and wait until application is finished.
func (app *App) Wait() ExitCodeOrError {
	st, ok := <-app.waitCh
	if ok {
		return st
	} else {
		return ExitCodeOrError{ExitCode: 0, Error: fmt.Errorf("Exited already")}
	}
}

// Kill terminate application started asynchronously.
func (app *App) Kill() error {
	//log.Println(fmt.Sprintf("Start killing app: %v", app.cmd))
	if IsLinuxMacOSFreeBSD() {
		// Kill not only main but all child processes,
		// so extract for this purpose group id.
		pgid, err := syscall.Getpgid(app.cmd.Process.Pid)
		if err != nil {
			return err
		}
		// Specifying gid with negative sign also results in the killing of child processes.
		err = syscall.Kill(-pgid, syscall.SIGKILL)
		if err != nil {
			return err
		}
	} else {
		// Kill only mother process
		err := app.cmd.Process.Kill()
		if err != nil {
			return err
		}
	}
	state := app.Wait()
	//log.Println(fmt.Sprintf("Done killing app: %v", app.cmd))
	return state.Error
}
