package filehandler

import (
	"fmt"

	"github.com/pkg/browser"
)

func CreateDownloadLink(filepath string) (string, error) {
	const baseURL = "http://localhost:8080/files"

	link := fmt.Sprintf("%s%s", baseURL, filepath)
	return link, nil
}

func DownloadFile(url string) error {
	err := browser.OpenURL(url)
	if err != nil {
		return fmt.Errorf("не удалось открыть ссылку в браузере: %w", err)
	}
	return nil
}
