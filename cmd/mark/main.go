package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/yanzyuyu/mark/pkg/bookmark"
	"github.com/yanzyuyu/mark/pkg/editor"
	"github.com/yanzyuyu/mark/pkg/git"
)

const version = "0.1.0"

const usage = `mark — version-controlled code bookmarks

Usage:
  mark add <label> <file>:<line>   add a bookmark
  mark ls                          list all bookmarks
  mark go <label>                  open bookmark in $EDITOR
  mark rm <label>                  remove a bookmark
  mark check                       verify bookmarks are still valid
  mark version                     print version

Examples:
  mark add "entry point" cmd/main.go:1
  mark add "hot path"    pkg/vm/vm.go:142
  mark ls
  mark go "hot path"
  mark rm "entry point"
  mark check

Bookmarks are stored in .marks at your git repo root.
Commit .marks to share bookmarks with your team.
`

func main() {
	if len(os.Args) < 2 {
		fmt.Print(usage)
		os.Exit(0)
	}

	switch os.Args[1] {
	case "add":
		cmdAdd(os.Args[2:])
	case "ls", "list":
		cmdList()
	case "go", "open":
		cmdGo(os.Args[2:])
	case "rm", "remove", "del", "delete":
		cmdRemove(os.Args[2:])
	case "check":
		cmdCheck()
	case "version", "--version", "-v":
		fmt.Printf("mark v%s\n", version)
	default:
		fmt.Fprintf(os.Stderr, "error: unknown command %q\n\n", os.Args[1])
		fmt.Print(usage)
		os.Exit(1)
	}
}

func repoRoot() string {
	root, err := git.Root()
	if err != nil {
		cwd, _ := os.Getwd()
		return cwd
	}
	return root
}

func loadStore() (*bookmark.Store, string) {
	root := repoRoot()
	s, err := bookmark.Load(root)
	if err != nil {
		die("failed to load .marks: %v", err)
	}
	return s, root
}

func cmdAdd(args []string) {
	if len(args) < 2 {
		die("usage: mark add <label> <file>:<line>")
	}

	label := args[0]
	ref := args[1]

	file, line, ok := parseRef(ref)
	if !ok {
		die("invalid reference %q — expected <file>:<line> (e.g. main.go:42)", ref)
	}

	s, root := loadStore()

	relFile, err := filepath.Rel(root, absPath(file))
	if err != nil {
		relFile = file
	}
	relFile = filepath.ToSlash(relFile)

	absFile := filepath.Join(root, relFile)

	bm, err := s.Add(label, relFile, absFile, line, git.Author())
	if err != nil {
		die("%v", err)
	}

	if err := s.Save(); err != nil {
		die("failed to save .marks: %v", err)
	}

	fmt.Printf("%s  %s  %s:%d\n", green("✓"), bold(label), absFile, line)
	fmt.Printf("   %s\n", dim(bm.Content))
}

func cmdList() {
	s, root := loadStore()

	if len(s.Bookmarks) == 0 {
		fmt.Println(dim("no bookmarks yet — use: mark add <label> <file>:<line>"))
		return
	}

	maxLabel := 0
	for _, b := range s.Bookmarks {
		if len(b.Label) > maxLabel {
			maxLabel = len(b.Label)
		}
	}

	fmt.Printf("\n  %s\n\n", bold("bookmarks"))
	for _, b := range s.Bookmarks {
		pad := strings.Repeat(" ", maxLabel-len(b.Label))
		loc := fmt.Sprintf("%s:%d", filepath.Join(root, b.File), b.Line)
		fmt.Printf("  %s%s  %s\n", cyan(b.Label), pad, dim(loc))
		if b.Content != "" {
			fmt.Printf("  %s%s  %s\n", strings.Repeat(" ", maxLabel), "  ", yellow(b.Content))
		}
	}
	fmt.Println()
}

func cmdGo(args []string) {
	if len(args) < 1 {
		die("usage: mark go <label>")
	}

	label := strings.Join(args, " ")
	s, root := loadStore()

	bm := s.Find(label)
	if bm == nil {
		die("bookmark %q not found", label)
	}

	absFile := filepath.Join(root, bm.File)
	if err := editor.Open(absFile, bm.Line); err != nil {
		die("failed to open editor: %v", err)
	}
}

func cmdRemove(args []string) {
	if len(args) < 1 {
		die("usage: mark rm <label>")
	}

	label := strings.Join(args, " ")
	s, _ := loadStore()

	if !s.Remove(label) {
		die("bookmark %q not found", label)
	}

	if err := s.Save(); err != nil {
		die("failed to save .marks: %v", err)
	}

	fmt.Printf("%s  removed %s\n", red("✕"), bold(label))
}

func cmdCheck() {
	s, root := loadStore()

	if len(s.Bookmarks) == 0 {
		fmt.Println(dim("no bookmarks to check"))
		return
	}

	results := s.Check(root)

	ok := 0
	for _, r := range results {
		switch {
		case r.OK:
			ok++
			fmt.Printf("  %s  %s\n", green("✓"), bold(r.Bookmark.Label))
		case r.Missing:
			fmt.Printf("  %s  %s — file not found or too short\n", red("✕"), bold(r.Bookmark.Label))
		case r.Drifted:
			fmt.Printf("  %s  %s — content changed at line %d\n", yellow("~"), bold(r.Bookmark.Label), r.Bookmark.Line)
			fmt.Printf("       expected: %s\n", dim(r.Bookmark.Content))
			fmt.Printf("       current:  %s\n", dim(r.Current))
		}
	}

	fmt.Printf("\n  %d/%d bookmarks valid\n", ok, len(results))
}

func parseRef(ref string) (string, int, bool) {
	idx := strings.LastIndex(ref, ":")
	if idx < 0 {
		return "", 0, false
	}
	file := ref[:idx]
	lineStr := ref[idx+1:]
	n, err := strconv.Atoi(lineStr)
	if err != nil || n < 1 {
		return "", 0, false
	}
	return file, n, true
}

func absPath(p string) string {
	if filepath.IsAbs(p) {
		return p
	}
	cwd, _ := os.Getwd()
	return filepath.Join(cwd, p)
}

func die(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
	os.Exit(1)
}

func bold(s string) string  { return "\033[1m" + s + "\033[0m" }
func dim(s string) string   { return "\033[2m" + s + "\033[0m" }
func green(s string) string { return "\033[32m" + s + "\033[0m" }
func red(s string) string   { return "\033[31m" + s + "\033[0m" }
func yellow(s string) string { return "\033[33m" + s + "\033[0m" }
func cyan(s string) string  { return "\033[36m" + s + "\033[0m" }
