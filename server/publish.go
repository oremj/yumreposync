package server

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func santizeFileName(name string) string {
	return filepath.Base(name)
}

func writeFile(name string, r io.Reader) error {
	f, err := os.Create("/tmp/testrepo/" + name)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = io.Copy(f, r); err != nil {
		return err
	}

	log.Print("Wrote: ", f.Name())
	return nil
}

func writeFiles(req *http.Request) error {
	mp, err := req.MultipartReader()
	if err != nil {
		return err
	}

	for {
		part, err := mp.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if err := writeFile(santizeFileName(part.FileName()), part); err != nil {
			return err
		}
	}
	return nil
}

func Publish(w http.ResponseWriter, req *http.Request) {
	if err := writeFiles(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error: Publish: %s", err)
		return
	}
}
