package handler

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	LogConstant "github.com/charizark2710/Sunflower/RDIPs-BE/constant/LogConst"
	"github.com/charizark2710/Sunflower/RDIPs-BE/utils"
)

type FileIO struct {
	Name     string
	rFile    *os.File
	wFile    *os.File
	isClosed bool
}

var fileStreamArr = make(map[string]*FileIO, 0)
var m sync.Mutex

func (f *FileIO) open(IOtype string, date time.Time) error {
	m.Lock()
	defer m.Unlock()
	utils.Log(LogConstant.Info, "Start open "+f.Name)

	if fileStreamArr[f.Name] == nil ||
		fileStreamArr[f.Name].wFile == nil ||
		IOtype == "read" {
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
				newFile, newErr := os.Create(path + "/" + hour + ".log")
				if newErr != nil {
					utils.Log(LogConstant.Error, newErr)
					return newErr
				}
				file = newFile
			} else {
				utils.Log(LogConstant.Error, err)
				return err
			}
		}
		switch IOtype {
		case "read":
			f.rFile = file
		case "write":
			f.wFile = file
			go func() {
				// time.Sleep(10 * time.Second)

				time.Sleep(1 * time.Hour)
				m.Lock()
				defer m.Unlock()
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
		}
		fileStreamArr[f.Name] = nil
		fileStreamArr[f.Name] = f
	} else {
		f.rFile = fileStreamArr[f.Name].rFile
		f.wFile = fileStreamArr[f.Name].wFile
	}

	utils.Log(LogConstant.Info, "Open "+f.Name+" success")

	return nil
}

func (f *FileIO) Read(date time.Time) ([]byte, error) {
	utils.Log(LogConstant.Info, "Start reading operation to "+f.Name+" log")

	// if file is not opened
	err := f.open("read", date)
	defer func() {
		f.rFile.Close()
	}()
	if err != nil {
		utils.Log(LogConstant.Error, err)
		return []byte{}, err
	}

	bytes, readErr := io.ReadAll(f.rFile)
	if readErr != nil {
		utils.Log(LogConstant.Error, readErr)
		return []byte{}, readErr
	}
	utils.Log(LogConstant.Info, "Finish reading operation to "+f.Name+" log: ", string(bytes))

	return bytes, nil
}

func (f *FileIO) Write(date time.Time, bytes []byte) error {
	utils.Log(LogConstant.Info, "Start writing operation to "+f.Name+" log")

	// if file is not opened
	err := f.open("write", date)
	if err != nil {
		utils.Log(LogConstant.Error, err)
		return err
	}
	w := bufio.NewWriter(f.wFile)
	n, wErr := w.Write(bytes)
	if wErr != nil {
		utils.Log(LogConstant.Error, err)
		return err
	}
	err = w.Flush()
	if err != nil {
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

	// rFile is already close after reading process is done so this is just for safe
	err := f.rFile.Close()

	if err != nil {
		if !errors.Is(err, os.ErrClosed) {
			f.isClosed = false
		}
		utils.Log(LogConstant.Warning, "Close Read "+f.Name+" falied")
		utils.Log(LogConstant.Warning, err)
	}

	// only need to close write file
	err = f.wFile.Close()
	if err != nil {
		if !errors.Is(err, os.ErrClosed) {
			f.isClosed = false
		}
		utils.Log(LogConstant.Warning, "Close Write"+f.Name+"falied")
		utils.Log(LogConstant.Warning, err)
	}
	if _, exist := fileStreamArr[f.Name]; exist {
		fileStreamArr[f.Name] = nil
		delete(fileStreamArr, f.Name)
	}
	utils.Log(LogConstant.Info, "Finish closing "+f.Name)

}
