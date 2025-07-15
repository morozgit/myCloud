package filehandler

import (
	"archive/zip"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"mycloud/file_api/config"
)

func CreateDownloadLink(path string) (string, error) {
	var baseURL = config.GetBaseURL()
	fullPath := filepath.Join("/home", path)
	log.Printf("fullPath: %s", fullPath)
	info, err := os.Stat(fullPath)
	if err != nil {
		log.Printf("не удалось получить информацию о файле/папке: %v", err)
		return "", err
	}

	var downloadLink string
	if info.IsDir() {
		archivePath := fmt.Sprintf("%s.zip", fullPath)

		err := zipFolder(fullPath, archivePath)
		if err != nil {
			log.Printf("не удалось архивировать папку: %v", err)
			return "", err
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
		log.Printf("не удалось создать архив: %v", err)
		return err
	}
	defer outFile.Close()

	zipWriter := zip.NewWriter(outFile)
	defer zipWriter.Close()

	err = filepath.Walk(source, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			log.Printf("ошибка при обработке файла: %v", err)
			return err
		}

		if fi.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(source, file)
		if err != nil {
			log.Printf("не удалось получить относительный путь: %v", err)
			return err
		}

		zipFile, err := zipWriter.Create(relPath)
		if err != nil {
			log.Printf("не удалось добавить файл в архив: %v", err)
			return err
		}

		inFile, err := os.Open(file)
		if err != nil {
			log.Printf("не удалось открыть файл: %v", err)
			return err
		}
		defer inFile.Close()

		_, err = fmt.Fprintln(zipFile, inFile)
		return err
	})

	return err
}
