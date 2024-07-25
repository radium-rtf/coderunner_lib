package files

import (
	"github.com/radium-rtf/coderunner_lib/file"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const (
	data = "data"
)

func Parse(lang, t string) ([]file.File, error) {
	dir := strings.Join([]string{data, lang, t, ""}, string(os.PathSeparator))

	var files []file.File
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		filename := strings.TrimPrefix(path, dir)
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		files = append(files, file.NewFile(filename, file.BytesContent(data)))
		return nil
	})
	return files, err
}

func Find(files []file.File, name string) file.File {
	index := slices.IndexFunc(files, func(f file.File) bool {
		return f.Name == name
	})
	return files[index]
}
