package iofile

import (
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/deweppro/go-sdk/errors"
)

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func Search(dir, filename string) ([]string, error) {
	files := make([]string, 0)
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || info.Name() != filename {
			return nil
		}
		files = append(files, path)
		return nil
	})
	return files, err
}

func Rewrite(filename string, call func([]byte) ([]byte, error)) error {
	var perm os.FileMode = 0777
	if ls, err := os.Lstat(filename); err != nil {
		perm = ls.Mode().Perm()
	}
	b, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	if b, err = call(b); err != nil {
		return err
	}
	return os.WriteFile(filename, b, perm)
}

func IsValidHash(filename string, h hash.Hash, valid string) error {
	r, err := os.Open(filename)
	if err != nil {
		return err
	}
	if _, err = io.Copy(h, r); err != nil {
		return errors.Wrapf(err, "calculate file hash")
	}
	result := hex.EncodeToString(h.Sum(nil))
	h.Reset()
	if result != valid {
		return fmt.Errorf("invalid hash: expected[%s] actual[%s]", valid, result)
	}
	return nil
}

func Hash(filename string, h hash.Hash) (string, error) {
	r, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	if _, err = io.Copy(h, r); err != nil {
		return "", errors.Wrapf(err, "calculate file hash")
	}
	result := hex.EncodeToString(h.Sum(nil))
	h.Reset()
	return result, nil
}

func Copy(dst, src string, mode os.FileMode) error {
	source, err := os.OpenFile(src, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer source.Close() //nolint: errcheck

	if mode == 0 {
		fi, err0 := source.Stat()
		if err0 != nil {
			return err0
		}
		mode = fi.Mode()
	}

	dist, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer dist.Close() //nolint: errcheck

	_, err = io.Copy(dist, source)
	return err
}
