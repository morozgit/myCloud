package filehandler

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func DownloadFile(filepath, destination string) error {
	sourceFile, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл: %w", err)
	}
	defer sourceFile.Close()

	err = os.MkdirAll(destination, os.ModePerm)
	if err != nil {
		return fmt.Errorf("не удалось создать директорию: %w", err)
	}

	destinationFilePath := destination + "/" + filepath[strings.LastIndex(filepath, "/")+1:]

	destinationFile, err := os.Create(destinationFilePath)
	if err != nil {
		return fmt.Errorf("не удалось создать файл: %w", err)
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return fmt.Errorf("не удалось скопировать данные: %w", err)
	}
	return nil
}
