package codegen

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/logrusorgru/aurora"
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

func WriteFile(filename string, content string) {
	dir := filepath.Dir(filename)

	if dir != "" {
		os.MkdirAll(dir, os.ModePerm)
	}

	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	n3, err := f.WriteString(content)
	if err != nil {
		panic(err)
	}
	f.Sync()

	pwd, _ := os.Getwd()

	log.Printf(aurora.Sprintf(aurora.Green("Generated file to %s(%d KiB, %d B)"), aurora.Blue(path.Join(pwd, filename)), n3/1024, n3))
}

func WriteGoFile(path string, content string) {
	WriteFile(path, content)
	exec.Command("gofmt", "-w", path).CombinedOutput()
	exec.Command("goimports", "-w", path).CombinedOutput()
}

func WriteJSONFile(path string, data interface{}) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}
	WriteFile(path, string(b))
}
