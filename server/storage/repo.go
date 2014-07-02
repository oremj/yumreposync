package storage

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"sync"
)

type Storage struct {
	RepoDir string
	lock    sync.Mutex
}

func (s *Storage) SanitizeFileName(name string) string {
	return filepath.Base(name)
}

func (s *Storage) writeFile(dir, name string, r io.Reader) error {
	sanitizedName := s.SanitizeFileName(name)
	f, err := os.Create(filepath.Join(dir, sanitizedName))
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

func (s *Storage) PublishMultiPartFiles(reader *multipart.Reader) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}

	params := make(map[string]string)
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if part.FileName() == "" {
			tmp := new(bytes.Buffer)
			_, err = tmp.ReadFrom(part)
			if err != nil {
				return err
			}
			params[part.FormName()] = tmp.String()
		} else if err := s.writeFile(tmpDir, part.FileName(), part); err != nil {
			return err
		}
	}

	fmt.Println(params)
	return nil
}
