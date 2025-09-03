package fs

import (
	"os"
	"path/filepath"
)

func FileWriter(basePath string, relPath string) (func([]byte) error, func() error, error) {
	fullPath := filepath.Join(basePath, relPath)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0775); err != nil {
		return nil, nil, err
	}

	file, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, nil, err
	}

	write := func(content []byte) error {
		_, err := file.Write(content)
		if err != nil {
			return err
		}
		return nil
	}

	close := func() error {
		err := file.Close()
		if err != nil {
			return err
		}
		return nil
	}

	return write, close, nil
}
