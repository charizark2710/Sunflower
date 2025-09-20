package handler

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
	"rogchap.com/v8go"
)

func LoadJSDir(dir string, deep int) []map[string]string {
	if deep <= 0 {
		return []map[string]string{}
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	var files []string
	var subResult []map[string]string

	// Collect all JS files and recursively process subdirectories
	for _, entry := range entries {
		if entry.IsDir() {
			deep--
			subResult = append(subResult, LoadJSDir(filepath.Join(dir, entry.Name()), deep)...)
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".js") {
			continue
		}
		files = append(files, filepath.Join(dir, entry.Name()))
	}

	// Process files and bundle them asynchronously
	result := make([]map[string]string, len(files)) // Pre-allocate with exact length
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i, file := range files {
		wg.Add(1)
		go func(index int, filePath string) {
			defer wg.Done()

			data, err := os.ReadFile(filePath)
			if err != nil {
				log.Print(err)
				return
			}

			bundledCode, err := BundleJSCode(string(data))
			if err != nil {
				log.Print(err)
				return
			}

			mu.Lock()
			result[index] = map[string]string{"bundledCode": string(bundledCode), "orgCode": string(data), "name": filePath}
			mu.Unlock()
		}(i, file)
	}

	wg.Wait()

	// Filter out nil elements (failed processing) - FIXED VERSION
	var filteredResult []map[string]string
	for _, r := range result {
		if r != nil && r["bundledCode"] != "" && r["name"] != "" {
			filteredResult = append(filteredResult, r)
		}
	}

	if filteredResult == nil {
		return subResult
	}
	return append(filteredResult, subResult...)
}

func perfEventOpen(attr *unix.PerfEventAttr, pid, cpu, groupFd, flags int) (int, error) {
	r0, _, e1 := unix.Syscall6(
		unix.SYS_PERF_EVENT_OPEN,
		uintptr(unsafe.Pointer(attr)),
		uintptr(pid),
		uintptr(cpu),
		uintptr(groupFd),
		uintptr(flags),
		0,
	)
	fd := int(r0)
	if fd == -1 {
		return -1, e1
	}
	return fd, nil
}

func startCounting() (int, int) {
	attr := unix.PerfEventAttr{
		Type:   unix.PERF_TYPE_HARDWARE,
		Config: unix.PERF_COUNT_HW_INSTRUCTIONS,
		Size:   uint32(unsafe.Sizeof(unix.PerfEventAttr{})),
		Bits:   unix.PerfBitExcludeKernel,
	}

	fd_ic, err := perfEventOpen(&attr, 0, -1, -1, unix.PERF_FLAG_FD_CLOEXEC)
	if err != nil {
		fmt.Println("perf_event_open failed:", err)
		os.Exit(1)
	}

	attr = unix.PerfEventAttr{
		Type:   unix.PERF_TYPE_HARDWARE,
		Config: unix.PERF_COUNT_HW_CPU_CYCLES,
		Size:   uint32(unsafe.Sizeof(unix.PerfEventAttr{})),
		Bits:   unix.PerfBitExcludeKernel,
	}

	fd_cycle, err := perfEventOpen(&attr, 0, -1, -1, unix.PERF_FLAG_FD_CLOEXEC)
	if err != nil {
		fmt.Println("perf_event_open failed:", err)
		os.Exit(1)
	}

	// Reset + enable
	unix.IoctlSetInt(fd_ic, unix.PERF_EVENT_IOC_RESET, 0)
	unix.IoctlSetInt(fd_ic, unix.PERF_EVENT_IOC_ENABLE, 0)
	unix.IoctlSetInt(fd_cycle, unix.PERF_EVENT_IOC_RESET, 0)
	unix.IoctlSetInt(fd_cycle, unix.PERF_EVENT_IOC_ENABLE, 0)

	return fd_ic, fd_cycle
}

func finishCounting(fd int) uint64 {
	unix.IoctlSetInt(fd, unix.PERF_EVENT_IOC_DISABLE, 0)

	// Read
	buf := make([]byte, 8)
	n, err := unix.Read(fd, buf)
	if err != nil || n != 8 {
		fmt.Println("failed to read perf counter:", err)
		return 0
	}
	val := *(*uint64)(unsafe.Pointer(&buf[0]))

	// (optional) reset for reuse
	unix.IoctlSetInt(fd, unix.PERF_EVENT_IOC_RESET, 0)

	unix.Close(fd)
	return val
}

func ExecuteJs(code string, retry int) (*v8go.Value, uint64, uint64, error) {
	iso := v8go.NewIsolate()
	ctx := v8go.NewContext(iso)

	fd_ic, fd_cycle := startCounting()
	value, err := ctx.RunScript(code, "main.js")
	count := finishCounting(fd_ic)
	cycle := finishCounting(fd_cycle)

	if err != nil {
		return nil, count, cycle, err
	}

	if retry > 0 && count == 0 {
		time.Sleep(500 * time.Millisecond)
		return ExecuteJs(code, retry-1)
	}

	return value, count, cycle, err
}
