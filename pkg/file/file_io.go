package file

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type fileSystem struct {
	path string
}

func NewFileSystem(path string) IFileSystem {
	return &fileSystem{path}
}

func NewStorageFile(config map[string]map[string]interface{}) string {
	if fileConf, ok := config["file_storage"]; ok {
		if _, ok := fileConf["path"]; ok {
			return fileConf["path"].(string)
		}
	}
	return ""
}

func (fs *fileSystem) CreateFile(nameFile string) (path *string, errors error) {
	err := os.MkdirAll(fs.path, os.ModePerm)
	if err != nil {
		panic(err)
	}

	filePath := filepath.Join(fs.path, nameFile+".json")
	f, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return &filePath, nil
}

func (fs *fileSystem) SaveData(nameFile string, data interface{}) (bool, error) {
	path := fs.path + nameFile + ".json"
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return false, err
	}
	err = os.WriteFile(path, jsonData, 0644)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (fs *fileSystem) IsFileExisting(nameFile string) bool {
	path := fs.path + nameFile + ".json"
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func (fs *fileSystem) LoadFile(fileName string) ([]byte, error) {
	filePath := fs.path + fileName + ".json"

	jsonFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)

	if err != nil {
		return nil, err
	}
	return byteValue, nil

}
