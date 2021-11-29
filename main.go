package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type parsedFile struct {
	Path string `json:"path"`
	Hash string `json:"hash"`
}

func getFileMD5(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		panic(err)
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Please specify a root directory and a file to write to")
		fmt.Println("./calculatehash $(pwd) map.json")
		return
	}

	if os.Args[1] == "" || (len(os.Args) < 3) {
		fmt.Println("Please specify a root directory and a file to write to")
		fmt.Println("./calculatehash $(pwd) map.json")
		return
	}

	var filesList []string
	var files []parsedFile
	var root string = os.Args[1]

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		filename := strings.Replace(path, root+"/", "", -1)
		if info.IsDir() || strings.HasPrefix(filename, ".") {
			return nil
		}

		filesList = append(filesList, path)
		if err != nil {
			panic(err)
		}
		return nil
	})

	for _, file := range filesList {
		files = append(files, parsedFile{Path: file, Hash: getFileMD5(file)})
	}

	filesJson, jsonErr := json.Marshal(files)
	if jsonErr != nil {
		panic(jsonErr)
	}

	writePath := os.Args[2] + "/map.json"
	writingFileErr := ioutil.WriteFile(writePath, filesJson, 0644)
	if writingFileErr != nil {
		panic(writingFileErr)
	}
}
