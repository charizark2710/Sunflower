package handler

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/utils"
)

type FileIO struct {
	Name     string
	rFile    *os.File
	wFile    *os.File
	isClosed bool
	m        sync.Mutex
}

var fileStreamCache sync.Map

func (f *FileIO) open(IOtype string, path string) error {
	f.m.Lock()
	defer f.m.Unlock()
	utils.Log(LogConstant.Info, "Start open "+f.Name)

	fileCache, ok := fileStreamCache.Load(f.Name)

	if fileCache == nil || !ok {
		// https://docs.nersc.gov/filesystems/unix-file-permissions/
		// https://man7.org/linux/man-pages/man2/openat.2.html
		file, err := os.OpenFile(path, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0777)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				mkDirErr := os.MkdirAll(filepath.Dir(path), 0777)
				if mkDirErr != nil {
					utils.Log(LogConstant.Error, mkDirErr)
					return mkDirErr
				}
				neWFile, newErr := os.Create(path)
				if newErr != nil {
					utils.Log(LogConstant.Error, newErr)
					return newErr
				}
				file = neWFile
			} else {
				utils.Log(LogConstant.Error, err)
				return err
			}
		}
		f.rFile = file
		f.wFile = file
		go func() {
			time.Sleep(5 * time.Minute)
			// time.Sleep(1 * time.Hour)
			f.m.Lock()
			defer f.m.Unlock()
			i := 0
			for !f.isClosed {
				f.close()
				if !f.isClosed {
					time.Sleep(10 * time.Second)
					i++
				}
				if i > 3 {
					break
				}
			}
		}()
		fileStreamCache.Store(f.Name, f)
	} else {
		f.rFile = fileCache.(*FileIO).rFile
		f.wFile = fileCache.(*FileIO).wFile
	}

	utils.Log(LogConstant.Info, "Open "+f.Name+" success")

	return nil
}

func (f *FileIO) Read(path string) ([]byte, error) {
	utils.Log(LogConstant.Info, "Start reading operation to "+f.Name+" log")

	// if file is not opened
	err := f.open("read", path)
	if err != nil {
		utils.Log(LogConstant.Error, err)
		return []byte{}, err
	}
	f.rFile.Seek(0, io.SeekStart)
	bytes, readErr := io.ReadAll(f.rFile)
	if readErr != nil {
		utils.Log(LogConstant.Error, readErr)
		return []byte{}, readErr
	}
	utils.Log(LogConstant.Info, "Finish reading operation to "+f.Name+" log: ", string(bytes))

	return bytes, nil
}

func (f *FileIO) Write(path string, value string) error {
	utils.Log(LogConstant.Info, "Start writing operation to "+f.Name+" log")

	// if file is not opened
	err := f.open("write", path)
	if err != nil {
		utils.Log(LogConstant.Error, err)
		return err
	}
	f.m.Lock()
	defer f.m.Unlock()
	n, wErr := fmt.Fprintln(f.wFile, value)
	if wErr != nil {
		utils.Log(LogConstant.Error, err)
		return err
	}
	utils.Log(LogConstant.Info, "Finish writing operation to "+f.Name+" log: ", n, "byte")
	return nil
}

// Close after priod of time
func (f *FileIO) close() {
	f.isClosed = true
	utils.Log(LogConstant.Info, "Start closing "+f.Name)

	err := f.rFile.Close()

	if err != nil {
		if !errors.Is(err, os.ErrClosed) {
			f.isClosed = false
			utils.Log(LogConstant.Error, err)
			return
		}
		utils.Log(LogConstant.Warning, "Close Read "+f.Name+" falied")
		utils.Log(LogConstant.Warning, err)
	}

	err = f.wFile.Close()
	if err != nil {
		if !errors.Is(err, os.ErrClosed) {
			f.isClosed = false
			utils.Log(LogConstant.Error, err)
			return
		}
		utils.Log(LogConstant.Warning, "Close Write"+f.Name+" falied")
		utils.Log(LogConstant.Warning, err)
	}
	if _, exist := fileStreamCache.Load(f.Name); exist {
		fileStreamCache.Delete(f.Name)
	}
	utils.Log(LogConstant.Info, "Finish closing "+f.Name)
	_, ok := fileStreamCache.Load(f.Name)
	utils.Log(LogConstant.Info, f.Name, ok)

}
