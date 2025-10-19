package commands

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"vc/workdir"
)

type VCStatus struct {
	ModifiedFiles []string
	StagedFiles   []string
}

type commit struct {
	message string
	files   map[string]string // snapshot of repo at this commit
}

type VC struct {
	workDir *workdir.WorkDir
	commits []*commit
	staged  map[string]string // staged file -> content
}

func Init(wd *workdir.WorkDir) *VC {
	return &VC{
		workDir: wd,
		commits: []*commit{},
		staged:  map[string]string{},
	}
}

func (vc *VC) GetWorkDir() *workdir.WorkDir {
	return vc.workDir
}

func (vc *VC) Add(files ...string) {
	for _, f := range files {
		content, err := vc.workDir.CatFile(f)
		if err == nil {
			vc.staged[f] = content
		}
	}
}

func (vc *VC) AddAll() {
	for _, f := range vc.workDir.ListFilesRoot() {
		content, err := vc.workDir.CatFile(f)
		if err == nil {
			vc.staged[f] = content
		}
	}
}

func cloneMap(m map[string]string) map[string]string {
	n := make(map[string]string, len(m))
	for k, v := range m {
		n[k] = v
	}
	return n
}

func (vc *VC) Commit(msg string) {
	// create snapshot based on last commit snapshot, then apply staged changes
	var base map[string]string
	if len(vc.commits) > 0 {
		base = cloneMap(vc.commits[len(vc.commits)-1].files)
	} else {
		base = make(map[string]string)
	}
	for f, content := range vc.staged {
		base[f] = content
	}
	c := &commit{
		message: msg,
		files:   base,
	}
	vc.commits = append(vc.commits, c)
	vc.staged = make(map[string]string)
}

func (vc *VC) Status() VCStatus {
	modified := []string{}
	staged := []string{}

	// لیست فایل‌های استیج شده
	for f := range vc.staged {
		staged = append(staged, f)
	}

	// اگر هنوز commitی وجود ندارد
	if len(vc.commits) == 0 {
		return VCStatus{
			ModifiedFiles: modified, // خالی
			StagedFiles:   staged,
		}
	}

	// آخرین commit
	last := vc.commits[len(vc.commits)-1].files

	// بررسی فایل‌های workdir
	for _, f := range vc.workDir.ListFilesRoot() {
		curContent, err := vc.workDir.CatFile(f)
		if err != nil {
			continue
		}

		// اگر فایل استیج شده است
		stagedContent, isStaged := vc.staged[f]
		lastContent, inCommit := last[f]

		if isStaged {
			// اگر بعد از Add تغییر کرده باشد
			if stagedContent != curContent {
				modified = append(modified, f)
			}
			continue
		}

		if !inCommit || lastContent != curContent {
			modified = append(modified, f)
		}
	}

	return VCStatus{
		ModifiedFiles: modified,
		StagedFiles:   staged,
	}
}

func (vc *VC) Log() []string {
	out := []string{}
	for i := len(vc.commits) - 1; i >= 0; i-- {
		out = append(out, vc.commits[i].message)
	}
	return out
}

func parseTilde(rev string) (int, error) {
	// ~N  where N is integer
	if !strings.HasPrefix(rev, "~") {
		return 0, errors.New("not tilde")
	}
	numStr := strings.TrimPrefix(rev, "~")
	n, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func parseCarets(rev string) (int, bool) {
	// ^^^ -> count carets
	for i := 0; i < len(rev); i++ {
		if rev[i] != '^' {
			return 0, false
		}
	}
	return len(rev), true
}

func (vc *VC) Checkout(rev string) (*workdir.WorkDir, error) {
	if len(vc.commits) == 0 {
		return nil, errors.New("no commits")
	}
	index := -1
	if strings.HasPrefix(rev, "~") {
		n, err := parseTilde(rev)
		if err != nil {
			return nil, err
		}
		index = len(vc.commits) - 1 - n
	} else if strings.HasPrefix(rev, "^") {
		count, ok := parseCarets(rev)
		if !ok {
			return nil, errors.New("invalid rev")
		}
		index = len(vc.commits) - 1 - count
	} else {
		return nil, errors.New("unsupported rev format")
	}

	if index < 0 || index >= len(vc.commits) {
		return nil, fmt.Errorf("revision out of range")
	}

	c := vc.commits[index]

	// clone current workdir to a new WorkDir
	newWD := vc.workDir.Clone()

	// Ensure directories for files; create/truncate files and write snapshot content
	for f, content := range c.files {
		dir := filepath.Dir(f)
		if dir != "." && dir != "" {
			_ = newWD.CreateDir(dir)
		}
		_ = newWD.CreateFile(f)
		_ = newWD.WriteToFile(f, content)
	}

	return newWD, nil
}
