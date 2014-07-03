package storage

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"os/exec"
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

func (s *Storage) updateMetadata() error {
	c := exec.Command("createrepo", ".")
	c.Dir = s.RepoDir
	out, err := c.CombinedOutput()
	if err != nil {
		return fmt.Errorf("createrepo: %s, %s", err, out)
	}
	return nil
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

func (s *Storage) PublishMultiPartFiles(form *multipart.Form) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, files := range form.File {
		for _, file := range files {
			tmp, err := file.Open()
			if err != nil {
				return err
			}
			if err := s.writeFile(file.Filename, tmp); err != nil {
				return err
			}
		}
	}

	if err := s.updateMetadata(); err != nil {
		return err
	}

	return nil
}
