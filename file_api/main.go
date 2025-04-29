package main

import (
	"log"
	"mycloud/file_api/config"
	"mycloud/file_api/rabbithandler"
	"net/http"
)

const serverAddr = ":8080"

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	config.LoadEnv()
	rabbithandler.StartFileServer()

	go rabbithandler.HandleMessages()

	log.Printf("Сервер запущен на http://localhost%s\n", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, nil))
}
