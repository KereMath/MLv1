package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func clearDirectory(dirPath string) error {
	// Klasördeki tüm dosya ve dizinleri listele
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %v", dirPath, err)
	}

	// Her dosya/dizin için işlem yap
	for _, file := range files {
		filePath := filepath.Join(dirPath, file.Name())
		err := os.RemoveAll(filePath) // Dosya veya dizini sil
		if err != nil {
			return fmt.Errorf("failed to remove %s: %v", filePath, err)
		}
	}
	return nil
}

func main() {
	// Parent directory'yi al
	baseDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current directory: %v", err)
	}
	parentDir := filepath.Dir(baseDir)

	// tempdata dizinini temizle
	tempDataPath := filepath.Join(parentDir, "scripts", "tempdata")
	err = clearDirectory(tempDataPath)
	if err != nil {
		log.Fatalf("Failed to clear tempdata directory: %v", err)
	}
	log.Println("Successfully cleared tempdata directory")

	// logs dizinini temizle
	logsPath := filepath.Join(parentDir, "logs")
	err = clearDirectory(logsPath)
	if err != nil {
		log.Fatalf("Failed to clear logs directory: %v", err)
	}
	log.Println("Successfully cleared logs directory")
}
