package filehandler

import (
	"fmt"
	"net/http"

	"mycloud/file_api/config"
)

func CreateDownloadLink(filepath string) (string, error) {
	var baseURL = config.GetBaseURL()

	link := fmt.Sprintf("%s%s", baseURL, filepath)
	return link, nil
}

func DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем путь файла из запроса
	filePath := r.URL.Query().Get("path")
	if filePath == "" {
		http.Error(w, "Не указан путь к файлу", http.StatusBadRequest)
		return
	}

	// Генерируем ссылку для скачивания
	downloadLink, err := CreateDownloadLink(filePath)
	if err != nil {
		http.Error(w, "Не удалось создать ссылку на файл", http.StatusInternalServerError)
		return
	}

	// Отправляем ссылку на фронтенд
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"download_link": "%s"}`, downloadLink)
}
