package server

import (
	"log"
	"net/http"
)

func writeFiles(req *http.Request) error {
	// Fit 1 GB in to memory
	if err := req.ParseMultipartForm(1024 << 20); err != nil {
		return err
	}
	defer req.MultipartForm.RemoveAll()

	return Storage.PublishMultiPartFiles(req.MultipartForm)
}

func Publish(w http.ResponseWriter, req *http.Request) {
	if err := writeFiles(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error: Publish: %s", err)
		return
	}
}
