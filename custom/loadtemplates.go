package custome

import (
	"path/filepath"
)

func LoadTemplates(templatelist []string) []string {

	files := []string{}
	for _, dir := range templatelist {
		ff, err := filepath.Glob(dir)
		if err != nil {
			panic(err)
		}
		files = append(files, ff...)
	}

	return files

}
