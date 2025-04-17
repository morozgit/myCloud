package filehandler

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func CreateDownloadLink(filepath string) (string, error) {
	const baseURL = "http://172.17.0.1:8080/files"

	link := fmt.Sprintf("%s%s", baseURL, filepath)
	return link, nil
}

func DownloadFile(url, filename, downloadDir string) error {

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("не удалось скачать файл: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("неожиданный статус ответа: %s", resp.Status)
	}

	err = os.MkdirAll(downloadDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("не удалось создать папку downloads: %w", err)
	}

	out, err := os.Create(downloadDir + filename)
	if err != nil {
		return fmt.Errorf("не удалось создать файл: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("ошибка при записи файла: %w", err)
	}

	return nil
}
