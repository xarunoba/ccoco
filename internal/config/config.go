package config

import (
	"encoding/json"
	"log"
	"os"
)

type File struct {
	Files    []string
	FilesDir string
	Strict   bool
}

const FileName = "ccoco.config.json"
const CcocoDir = ".ccoco"

const CacheDir = CcocoDir + "/cache"
const ConfigsDir = CcocoDir + "/configs"
const PreflightsDir = CcocoDir + "/preflights"
const CcocoExecutable = CacheDir + "/ccoco"

var defaultFile = File{
	Files:    []string{"env"},
	FilesDir: ".",
	Strict:   false,
}

func GetFile() File {
	data, err := os.ReadFile(FileName)
	if err != nil {
		return defaultFile
	}

	var configFile File
	if err := json.Unmarshal(data, &configFile); err != nil {
		log.Fatalf("Error unmarshalling %s: %v", FileName, err)
	}

	if len(configFile.Files) == 0 {
		configFile.Files = defaultFile.Files
	}
	if configFile.FilesDir == "" {
		configFile.FilesDir = defaultFile.FilesDir
	}

	return configFile
}
