package storage

import (
	"io"
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

func (s *Storage) writeFile(name string, r io.Reader) error {
	sanitizedName := s.SanitizeFileName(name)
	f, err := os.Create(filepath.Join(s.RepoDir, sanitizedName))
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
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if err := s.writeFile(part.FileName(), part); err != nil {
			return err
		}
	}
	return nil
}
