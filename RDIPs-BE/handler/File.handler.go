package handler

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/utils"
)

type FileIO struct {
	Name     string
	RFile    *os.File
	WFile    *os.File
	isClosed bool
	m        sync.Mutex
}

var fileStreamCache sync.Map

func (f *FileIO) open(IOtype string, date time.Time) error {
	f.m.Lock()
	defer f.m.Unlock()
	utils.Log(LogConstant.Info, "Start open "+f.Name)

	fileCache, ok := fileStreamCache.Load(f.Name)

	if fileCache == nil || !ok {
		Y, M, D := date.Date()
		YMD := fmt.Sprint(Y) + "-" + fmt.Sprint(M) + "-" + fmt.Sprint(D)
		hour := fmt.Sprint(date.Hour())
		path := "logs/" + f.Name + "/" + YMD

		// https://docs.nersc.gov/filesystems/unix-file-permissions/
		// https://man7.org/linux/man-pages/man2/openat.2.html
		file, err := os.OpenFile(path+"/"+hour+".log", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0777)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				os.MkdirAll(path, 0777)
				neWFile, newErr := os.Create(path + "/" + hour + ".log")
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
		f.RFile = file
		f.WFile = file
		go func() {
			time.Sleep(10 * time.Second)
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

		f.RFile = fileCache.(*FileIO).RFile
		f.WFile = fileCache.(*FileIO).WFile
	}

	utils.Log(LogConstant.Info, "Open "+f.Name+" success")

	return nil
}

func (f *FileIO) Read(date time.Time) ([]byte, error) {
	utils.Log(LogConstant.Info, "Start reading operation to "+f.Name+" log")

	// if file is not opened
	err := f.open("read", date)
	if err != nil {
		utils.Log(LogConstant.Error, err)
		return []byte{}, err
	}
	f.RFile.Seek(0, io.SeekStart)
	bytes, readErr := io.ReadAll(f.RFile)
	if readErr != nil {
		utils.Log(LogConstant.Error, readErr)
		return []byte{}, readErr
	}
	utils.Log(LogConstant.Info, "Finish reading operation to "+f.Name+" log: ", string(bytes))

	return bytes, nil
}

func (f *FileIO) Write(date time.Time, value string) error {
	utils.Log(LogConstant.Info, "Start writing operation to "+f.Name+" log")

	// if file is not opened
	err := f.open("write", date)
	if err != nil {
		utils.Log(LogConstant.Error, err)
		return err
	}
	f.m.Lock()
	defer f.m.Unlock()
	n, wErr := fmt.Fprintln(f.WFile, value)
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

	err := f.RFile.Close()

	if err != nil {
		if !errors.Is(err, os.ErrClosed) {
			f.isClosed = false
			utils.Log(LogConstant.Error, err)
			return
		}
		utils.Log(LogConstant.Warning, "Close Read "+f.Name+" falied")
		utils.Log(LogConstant.Warning, err)
	}

	err = f.WFile.Close()
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
