package helpers

import (
	"encoding/json"
	"github.com/logrusorgru/aurora"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func OpenFile(path string) string {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	if len(data) == 0 {
		panic("empty file data")
	}

	return string(data)
}

func WriteFile(path string, content string) {
	dir := filepath.Dir(path)

	if dir != "" {
		os.MkdirAll(dir, os.ModePerm)
	}

	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	n3, err := f.WriteString(content)
	if err != nil {
		panic(err)
	}
	f.Sync()

	log.Printf(aurora.Sprintf(aurora.Green("Generated file %s(%d KiB, %d B)"), aurora.Blue(path), n3/1024, n3))
}

func WriteGoFile(path string, content string) {
	WriteFile(path, content)
	exec.Command("gofmt", "-w", path).CombinedOutput()
}

func WriteJSONFile(path string, data interface{}) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}
	WriteFile(path, string(b))
}
