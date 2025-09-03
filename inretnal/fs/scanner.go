package fs

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"new_filesync/proto"
	"os"
	"path/filepath"
)

func ScanDir(root string) (*proto.FileListResponse, error) {
	files := make([]*proto.FileMeta, 0)

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err // ошибка доступа к файлу
		}

		if d.IsDir() {
			return nil // пропускаем директории
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		hash, err := computeSHA1(path)
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		files = append(files, &proto.FileMeta{
			Path:         filepath.ToSlash(relPath),
			Size:         info.Size(),
			ModifiedUnix: info.ModTime().Unix(),
			Hash:         hash,
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &proto.FileListResponse{Files: files}, nil
}

// computeSHA1 считает SHA1-хеш файла
func computeSHA1(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	hasher := sha1.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
