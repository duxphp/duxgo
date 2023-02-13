package storage

import (
	"context"
	"io"
	"os"
)

type LocalStorage struct {
	root   string
	domain string
}

func NewLocalStorage(root string, domain string) *LocalStorage {
	return &LocalStorage{root, domain}
}

func (s *LocalStorage) write(ctx context.Context, path string, contents string, config map[string]interface{}) error {
	fullPath := s.root + "/" + path
	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(contents)
	if err != nil {
		return err
	}
	return nil
}

func (s *LocalStorage) writeStream(ctx context.Context, path string, stream io.Reader, config map[string]interface{}) error {
	fullPath := s.root + "/" + path
	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, stream)
	if err != nil {
		return err
	}
	return nil
}

func (s *LocalStorage) read(ctx context.Context, path string) (string, error) {
	fullPath := s.root + "/" + path
	contents, err := os.ReadFile(fullPath)
	if err != nil {
		return "", err
	}
	return string(contents), nil
}

func (s *LocalStorage) readStream(ctx context.Context, path string) (io.Reader, error) {
	fullPath := s.root + "/" + path
	f, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (s *LocalStorage) delete(ctx context.Context, path string) error {
	fullPath := s.root + "/" + path
	return os.Remove(fullPath)
}

func (s *LocalStorage) publicUrl(ctx context.Context, path string) (string, error) {
	return s.domain + "/" + path, nil
}
