package main

import (
	"fmt"
	"log"
	"mycloud/file_api/config"
	"mycloud/file_api/rabbithandler"
	"net/http"
)

const serverAddr = ":8080"

func main() {
	config.LoadEnv()
	rabbithandler.StartFileServer()

	go rabbithandler.HandleMessages()

	fmt.Printf("Сервер запущен на http://localhost%s\n", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, nil))
}
