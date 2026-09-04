package bookmark_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yanzyuyu/mark/pkg/bookmark"
)

func TestAddAndFind(t *testing.T) {
	dir := t.TempDir()

	f := filepath.Join(dir, "main.go")
	os.WriteFile(f, []byte("package main\n\nfunc main() {}\n"), 0644)

	s, err := bookmark.Load(dir)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	_, err = s.Add("entry", "main.go", f, 1, "testuser")
	if err != nil {
		t.Fatalf("Add: %v", err)
	}

	found := s.Find("entry")
	if found == nil {
		t.Fatal("Find: expected bookmark, got nil")
	}
	if found.Label != "entry" {
		t.Errorf("label = %q, want %q", found.Label, "entry")
	}
	if found.Line != 1 {
		t.Errorf("line = %d, want 1", found.Line)
	}
	if found.Content != "package main" {
		t.Errorf("content = %q, want %q", found.Content, "package main")
	}
}

func TestAddDuplicateLabel(t *testing.T) {
	dir := t.TempDir()

	f := filepath.Join(dir, "a.go")
	os.WriteFile(f, []byte("hello\nworld\n"), 0644)

	s, _ := bookmark.Load(dir)
	s.Add("dup", "a.go", f, 1, "x")

	_, err := s.Add("dup", "a.go", f, 2, "x")
	if err == nil {
		t.Fatal("expected error for duplicate label, got nil")
	}
}

func TestRemove(t *testing.T) {
	dir := t.TempDir()

	f := filepath.Join(dir, "b.go")
	os.WriteFile(f, []byte("alpha\nbeta\n"), 0644)

	s, _ := bookmark.Load(dir)
	s.Add("alpha line", "b.go", f, 1, "x")

	if !s.Remove("alpha line") {
		t.Fatal("Remove returned false, expected true")
	}
	if s.Find("alpha line") != nil {
		t.Fatal("bookmark still exists after Remove")
	}
}

func TestSaveLoad(t *testing.T) {
	dir := t.TempDir()

	f := filepath.Join(dir, "c.go")
	os.WriteFile(f, []byte("line one\nline two\n"), 0644)

	s, _ := bookmark.Load(dir)
	s.Add("first", "c.go", f, 1, "author")
	s.Add("second", "c.go", f, 2, "author")

	if err := s.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	s2, err := bookmark.Load(dir)
	if err != nil {
		t.Fatalf("Load after Save: %v", err)
	}
	if len(s2.Bookmarks) != 2 {
		t.Errorf("expected 2 bookmarks, got %d", len(s2.Bookmarks))
	}
}

func TestCheck(t *testing.T) {
	dir := t.TempDir()

	f := filepath.Join(dir, "d.go")
	os.WriteFile(f, []byte("original\n"), 0644)

	s, _ := bookmark.Load(dir)
	s.Add("check this", "d.go", f, 1, "x")

	results := s.Check(dir)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].OK {
		t.Errorf("expected OK=true, got drifted=%v missing=%v", results[0].Drifted, results[0].Missing)
	}

	os.WriteFile(f, []byte("changed content\n"), 0644)

	results2 := s.Check(dir)
	if !results2[0].Drifted {
		t.Errorf("expected Drifted=true after file change")
	}
}
