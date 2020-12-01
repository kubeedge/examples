package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"sync"
)

type File struct {
	sync.RWMutex
	Path string
	File *os.File
}

func (v *File) Flush() error {
	v.Lock()
	defer v.Unlock()
	if v.File != nil {
		err := v.File.Sync()
		v.File = nil
		return err
	}
	return nil
}

func (v *File) Close() error {
	v.Lock()
	defer v.Unlock()
	if v.File != nil {
		err := v.File.Close()
		v.File = nil
		return err
	}
	return nil
}

func (v *File) getRotatedFileList() ([]logFile, error) {
	var list []logFile
	err := filepath.Walk(path.Dir(v.Path), func(p string,
		info os.FileInfo, err error) error {
		pattern := "*" + path.Base(v.Path) + "*"
		if ok, err := path.Match(pattern, path.Base(p)); ok && err == nil {
			i := extractIndex(info)
			list = append(list, logFile{FileInfo: info, Index: i})
		} else if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	s := &sortLogFiles{Items: list}
	sort.Sort(s)
	return s.Items, nil
}

func (v *File) doRotate(items []logFile, rotateMaxCount int) error {
	if len(items) > 0 {
		// delete last files
		deleteCount := len(items) - rotateMaxCount + 1
		if deleteCount > 0 {
			for i := 0; i < deleteCount; i++ {
				err := os.Remove(items[i].FileInfo.Name())
				if err != nil {
					return err
				}
			}
			items = items[deleteCount:]
		}
		// change names of rest files
		baseFilePath := items[len(items)-1].FileInfo.Name()
		movs := make([]int, len(items))
		// 1st round to change names
		for i, item := range items {
			movs[i] = i + 100000
			err := os.Rename(item.FileInfo.Name(),
				fmt.Sprintf("%s.%d", baseFilePath, movs[i]))
			if err != nil {
				return err
			}
		}
		// 2nd round to change names
		for i, item := range movs {
			err := os.Rename(fmt.Sprintf("%s.%d", baseFilePath, item),
				fmt.Sprintf("%s.%d", baseFilePath, len(items)-i))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (v *File) rotateFiles(rotateMaxSize int64, rotateMaxCount int) error {
	fs, err := v.File.Stat()
	if err != nil {
		return err
	}
	if fs.Size() > rotateMaxSize {
		if v.File != nil {
			err := v.File.Close()
			if err != nil {
				return err
			}
			v.File = nil
		}
		list, err := v.getRotatedFileList()
		if err != nil {
			return err
		}
		if err = v.doRotate(list, rotateMaxCount); err != nil {
			return err
		}
	}
	return nil
}

func (v *File) getFile() (*os.File, error) {
	v.Lock()
	defer v.Unlock()
	if v.File == nil {
		file, err := os.OpenFile(v.Path, os.O_RDWR|os.O_APPEND, 0660)
		if err != nil {
			file, err = os.Create(v.Path)
			if err != nil {
				return nil, err
			}
		}
		v.File = file
	}
	return v.File, nil
}

func (v *File) writeToFile(msg string, rotateMaxSize int64, rotateMaxCount int) error {
	file, err := v.getFile()
	if err != nil {
		return err
	}
	v.Lock()
	defer v.Unlock()
	var buf bytes.Buffer
	buf.WriteString(msg)
	buf.WriteString(fmt.Sprintln())
	if _, err := io.Copy(file, &buf); err != nil {
		return err
	}
	//	if err = file.Sync(); err != nil {
	//		return err
	//	}
	if err := v.rotateFiles(rotateMaxSize, rotateMaxCount); err != nil {
		return err
	}

	return nil
}
