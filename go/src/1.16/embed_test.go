package __16

import (
	"embed"
	"fmt"
	"io/fs"
	"testing"
)

//go:embed assets
var desc embed.FS

func TestEmbed(t *testing.T) {
	// readDir(desc, "assets")
	err := fs.WalkDir(desc, "assets", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		content, err := desc.ReadFile(path)
		if err != nil {
			return err
		}

		fmt.Printf("Content of file[%s]:%s\n", path, content)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

// Deprecated： Use fs.ReadDir instead
func readDir(desc embed.FS, path string) {
	entryList, err := fs.ReadDir(desc, path)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, entry := range entryList {
		completePath := path + "/" + entry.Name()

		if entry.IsDir() {
			readDir(desc, completePath)
			continue
		}

		content, err := desc.ReadFile(completePath)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Content of file[%s]:%s\n", completePath, content)
	}
}
