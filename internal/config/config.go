package config

import (
	"encoding/json"
	"log"
	"os"
)

type File struct {
	Files    []string `json:"files"`
	FilesDir string   `json:"filesDir"`
}

const FileName = "ccoco.config.json"
const CcocoDir = ".ccoco"

const CacheDir = CcocoDir + "/cache"
const ConfigsDir = CcocoDir + "/configs"
const PreflightsDir = CcocoDir + "/preflights"

var DefaultFile = File{
	Files:    []string{".env"},
	FilesDir: ".",
}

func GetFile() File {
	data, err := os.ReadFile(FileName)
	if err != nil {
		return DefaultFile
	}

	var configFile File
	if err := json.Unmarshal(data, &configFile); err != nil {
		log.Fatalf("Error unmarshalling %s: %v", FileName, err)
	}

	if len(configFile.Files) == 0 {
		log.Printf("Files are empty. Using default files: %v", DefaultFile.Files)
		configFile.Files = DefaultFile.Files
	}
	if configFile.FilesDir == "" {
		log.Printf("FilesDir is empty. Using default filesDir: %v", DefaultFile.FilesDir)
		configFile.FilesDir = DefaultFile.FilesDir
	}

	return configFile
}
