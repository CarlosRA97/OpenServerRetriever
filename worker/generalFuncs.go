package worker

import (
	"fmt"
	"io/ioutil"
	"log"
)

func Check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func firebaseDB(path string) string {
	return fmt.Sprintf("%s.json", path)
}

func authQuery(key string) string {
	return fmt.Sprintf("auth=%s", key)
}

func getMapsFrom(path string) []string {
	var maps []string
	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		nonMapsFolders := f.Name() != ".DS_Store" &&
			f.Name() != ".idea" &&
			f.Name() != "logs" &&
			f.Name() != "mapas-originales"

		if nonMapsFolders {
			if f.IsDir() {
				maps = append(maps, f.Name())
			}
		}
	}
	return maps
}
