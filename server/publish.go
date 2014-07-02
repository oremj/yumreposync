package server

import (
	"log"
	"net/http"
)

func writeFiles(req *http.Request) error {
	mp, err := req.MultipartReader()
	if err != nil {
		return err
	}
	return Storage.PublishMultiPartFiles(mp)
}

func Publish(w http.ResponseWriter, req *http.Request) {
	if err := writeFiles(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error: Publish: %s", err)
		return
	}
}
