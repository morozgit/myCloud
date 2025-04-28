package filehandler

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"

	"mycloud/file_api/config"
)

func CreateDownloadLink(path string) (string, error) {
	var baseURL = config.GetBaseURL()
	fullPath := filepath.Join("/home", path)
	info, err := os.Stat(fullPath)
	if err != nil {
		return "", fmt.Errorf("не удалось получить информацию о файле/папке: %w", err)
	}

	var downloadLink string
	if info.IsDir() {
		archivePath := fmt.Sprintf("%s.zip", fullPath)

		err := zipFolder(fullPath, archivePath)
		if err != nil {
			return "", fmt.Errorf("не удалось архивировать папку: %w", err)
		}

		downloadLink = fmt.Sprintf("%s%s.zip", baseURL, path)
	} else {
		downloadLink = fmt.Sprintf("%s%s", baseURL, path)
	}

	return downloadLink, nil
}

func zipFolder(source string, target string) error {
	outFile, err := os.Create(target)
	if err != nil {
		return fmt.Errorf("не удалось создать архив: %w", err)
	}
	defer outFile.Close()

	zipWriter := zip.NewWriter(outFile)
	defer zipWriter.Close()

	err = filepath.Walk(source, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("ошибка при обработке файла: %w", err)
		}

		if fi.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(source, file)
		if err != nil {
			return fmt.Errorf("не удалось получить относительный путь: %w", err)
		}

		zipFile, err := zipWriter.Create(relPath)
		if err != nil {
			return fmt.Errorf("не удалось добавить файл в архив: %w", err)
		}

		inFile, err := os.Open(file)
		if err != nil {
			return fmt.Errorf("не удалось открыть файл: %w", err)
		}
		defer inFile.Close()

		_, err = fmt.Fprintln(zipFile, inFile)
		return err
	})

	return err
}
