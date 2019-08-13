package main

import (
	"handler"
	"log"
	"net/http"
	"fmt"
)

func main() {

	http.HandleFunc("/listall", handler.ListAllHandler)
	http.HandleFunc("/", handlefile)
	fmt.Printf("Starting server for testing HTTP POST...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func handlefile(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handler.DownloadFileHandler(w, r)
	case "POST":
		handler.UploadFileHandler(w, r)
	case "DELETE":
		handler.DeleteFileHandler(w, r)
	case "HEAD":
		handler.GetMetaHandler(w, r)
	case "PUT":
		handler.UpdatePathHandler(w, r)
	default:
		fmt.Fprintf(w, "Sorry, method supported.")
	}
}

