package bookmark

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const filename = ".marks"

type Store struct {
	Version   int        `json:"version"`
	Bookmarks []Bookmark `json:"bookmarks"`
	path      string
}

func Load(root string) (*Store, error) {
	p := filepath.Join(root, filename)
	s := &Store{Version: 1, path: p}

	data, err := os.ReadFile(p)
	if errors.Is(err, os.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, s); err != nil {
		return nil, fmt.Errorf("corrupt .marks file: %w", err)
	}
	s.path = p
	return s, nil
}

func (s *Store) Save() error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, append(data, '\n'), 0644)
}

func (s *Store) Add(label, file, absFile string, line int, author string) (*Bookmark, error) {
	for _, b := range s.Bookmarks {
		if strings.EqualFold(b.Label, label) {
			return nil, fmt.Errorf("bookmark %q already exists (use rm to remove it first)", label)
		}
	}

	content, err := readLine(absFile, line)
	if err != nil {
		return nil, fmt.Errorf("cannot read %s:%d: %w", file, line, err)
	}

	bm := Bookmark{
		ID:        randomID(),
		Label:     label,
		File:      file,
		Line:      line,
		Content:   strings.TrimSpace(content),
		CreatedAt: time.Now().UTC(),
		Author:    author,
	}
	s.Bookmarks = append(s.Bookmarks, bm)
	return &bm, nil
}

func (s *Store) Remove(label string) bool {
	for i, b := range s.Bookmarks {
		if strings.EqualFold(b.Label, label) {
			s.Bookmarks = append(s.Bookmarks[:i], s.Bookmarks[i+1:]...)
			return true
		}
	}
	return false
}

func (s *Store) Find(label string) *Bookmark {
	for i := range s.Bookmarks {
		if strings.EqualFold(s.Bookmarks[i].Label, label) {
			return &s.Bookmarks[i]
		}
	}
	return nil
}

type CheckResult struct {
	Bookmark Bookmark
	OK       bool
	Current  string
	Drifted  bool
	Missing  bool
}

func (s *Store) Check(root string) []CheckResult {
	results := make([]CheckResult, 0, len(s.Bookmarks))
	for _, b := range s.Bookmarks {
		r := CheckResult{Bookmark: b}
		full := filepath.Join(root, b.File)
		line, err := readLine(full, b.Line)
		if err != nil {
			r.Missing = true
		} else {
			r.Current = strings.TrimSpace(line)
			if r.Current == b.Content {
				r.OK = true
			} else {
				r.Drifted = true
			}
		}
		results = append(results, r)
	}
	return results
}

func readLine(file string, n int) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	cur := 0
	for sc.Scan() {
		cur++
		if cur == n {
			return sc.Text(), nil
		}
	}
	return "", fmt.Errorf("file has fewer than %d lines", n)
}

func randomID() string {
	return fmt.Sprintf("%08x", rand.Int63())
}
