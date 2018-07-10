package utils

import (
	"log"
	"os"
)

func FindBasePath() string {
	basePath := os.Getenv("TADA_PATH")
	if basePath == "" {
		path, err := ExpandHomeDir("~/.tada")
		if err != nil {
			log.Printf("%v", err)
		}
		basePath = path
	}

	return basePath
}
