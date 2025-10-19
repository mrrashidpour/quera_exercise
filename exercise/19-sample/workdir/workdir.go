package workdir

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/otiai10/copy"
)

type WorkDir struct {
	root string
}

// InitEmptyWorkDir creates a temporary directory and returns a WorkDir
func InitEmptyWorkDir() *WorkDir {
	tmpDir, _ := ioutil.TempDir("", "workdir")
	return &WorkDir{root: tmpDir}
}

// helper to get full path
func (wd *WorkDir) fullPath(p string) string {
	return filepath.Join(wd.root, p)
}

func (wd *WorkDir) CreateFile(path string) error {
	full := wd.fullPath(path)
	dir := filepath.Dir(full)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	f, err := os.Create(full)
	if err != nil {
		return err
	}
	return f.Close()
}

func (wd *WorkDir) CreateDir(path string) error {
	full := wd.fullPath(path)
	return os.MkdirAll(full, 0755)
}

func (wd *WorkDir) WriteToFile(path, content string) error {
	full := wd.fullPath(path)
	if _, err := os.Stat(full); os.IsNotExist(err) {
		return errors.New("file does not exist")
	}
	return ioutil.WriteFile(full, []byte(content), 0644)
}

func (wd *WorkDir) AppendToFile(path, content string) error {
	full := wd.fullPath(path)
	if _, err := os.Stat(full); os.IsNotExist(err) {
		return errors.New("file does not exist")
	}
	f, err := os.OpenFile(full, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	return err
}

func (wd *WorkDir) CatFile(path string) (string, error) {
	full := wd.fullPath(path)
	b, err := ioutil.ReadFile(full)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (wd *WorkDir) ListFilesRoot() []string {
	var files []string
	_ = filepath.Walk(wd.root, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // ignore
		}
		if !info.IsDir() {
			rel, err := filepath.Rel(wd.root, p)
			if err == nil {
				files = append(files, filepath.ToSlash(rel))
			}
		}
		return nil
	})
	return files
}

func (wd *WorkDir) ListFilesIn(dir string) ([]string, error) {
	fullDir := wd.fullPath(dir)
	var files []string
	err := filepath.Walk(fullDir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			rel, err := filepath.Rel(wd.root, p)
			if err == nil {
				files = append(files, filepath.ToSlash(rel))
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (wd *WorkDir) Clone() *WorkDir {
	tmpDir, _ := ioutil.TempDir("", "workdir_clone")
	_ = copy.Copy(wd.root, tmpDir)
	return &WorkDir{root: tmpDir}
}
