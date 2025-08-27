package handler

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/evanw/esbuild/pkg/api"
	"rogchap.com/v8go"
)

func BundleJSCode(jsCode string) ([]byte, error) {
	result := api.Transform(jsCode, api.TransformOptions{
		Loader:            api.LoaderJS,
		MinifySyntax:      true,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
	})
	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("%s", result.Errors[0].Text)
	}
	return result.Code, nil
}

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

func ExecuteJs(code string) (*v8go.Value, int64, uint64, error) {

	iso := v8go.NewIsolate()
	ctx := v8go.NewContext(iso)

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	start := time.Now()
	value, err := ctx.RunScript(code, "main.js")
	elapsed := time.Since(start)

	if err != nil {
		return nil, 0, 0, err
	}

	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)

	memUsage := m.Alloc - m2.Alloc
	runtime.GC()
	return value, elapsed.Microseconds(), memUsage / 1024, err

}
