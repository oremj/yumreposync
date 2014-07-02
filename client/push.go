package client

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func copyFile(fileName string, w io.Writer) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(w, f)
	return err
}

func Push(url string, fileNames []string) error {
	body := new(bytes.Buffer)

	fw := multipart.NewWriter(body)

	for i, fileName := range fileNames {
		baseFileName := filepath.Base(fileName)
		part, err := fw.CreateFormFile(fmt.Sprintf("file%d", i), baseFileName)
		if err != nil {
			return err
		}

		err = copyFile(fileName, part)
		if err != nil {
			return err
		}
	}
	fw.Close()

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", fw.FormDataContentType())

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Push returned status: %d", resp.StatusCode)
	}

	return nil
}
