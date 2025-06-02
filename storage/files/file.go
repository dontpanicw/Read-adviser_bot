package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"tgBot/lib/e"
	"tgBot/storage"
	"time"
)

type Storage struct {
	basePath string
}

const (
	basePerm = 0774
)

func New(basePath string) *Storage {
	return &Storage{basePath: basePath}
}

func (s *Storage) Save(page *storage.Page) error {
	filePath := filepath.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(filePath, basePerm); err != nil {
		return e.Wrap(err, "can't create directory")
	}
	fName, err := fileName(page)
	if err != nil {
		return e.Wrap(err, "can't create filename")
	}

	fPath := filepath.Join(filePath, fName)

	file, err := os.Create(fPath)
	if err != nil {
		return e.Wrap(err, "can't create file")
	}
	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return e.Wrap(err, "can't encode file")
	}

	return nil

}

func (s *Storage) PickRandom(userName string) (page *storage.Page, err error) {
	filePath := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(filePath)
	if err != nil {
		return nil, e.Wrap(err, "can't read directory")
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))

	file := files[n]

	result, err := s.decodePage(filepath.Join(filePath, file.Name()))
	if err != nil {
		return nil, e.Wrap(err, "can't decode page")
	}

	return result, nil

}

func (s *Storage) Remove(page *storage.Page) error {
	fileName, err := fileName(page)
	if err != nil {
		return e.Wrap(err, "can't remove filename")
	}

	path := filepath.Join(s.basePath, page.UserName, fileName)

	if err := os.Remove(path); err != nil {
		msg := fmt.Sprintf("can't remove file %s:", fileName)
		return e.Wrap(err, msg)
	}
	return nil
}

func (s *Storage) IsExists(page *storage.Page) (bool, error) {
	fileName, err := fileName(page)
	if err != nil {
		return false, e.Wrap(err, "can't read directory")
	}

	path := filepath.Join(s.basePath, page.UserName, fileName)

	switch _, err := os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		return false, e.Wrap(err, "can't check if file exists")
	}

	return true, nil
}

func (s *Storage) decodePage(filePath string) (*storage.Page, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap(err, "can't open file")
	}
	defer func() { _ = file.Close() }()

	var p storage.Page

	if err := gob.NewDecoder(file).Decode(&p); err != nil {
		return nil, e.Wrap(err, "can't decode file")
	}
	return &p, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
