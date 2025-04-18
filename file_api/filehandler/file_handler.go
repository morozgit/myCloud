package filehandler

import (
	"fmt"

	"mycloud/file_api/config"

	"github.com/pkg/browser"
)

func CreateDownloadLink(filepath string) (string, error) {
	var baseURL = config.GetBaseURL()

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
