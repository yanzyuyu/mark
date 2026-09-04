package editor

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Open(file string, line int) error {
	ed := os.Getenv("EDITOR")
	if ed == "" {
		ed = os.Getenv("VISUAL")
	}
	if ed == "" {
		if p, err := exec.LookPath("code"); err == nil && p != "" {
			ed = "code"
		}
	}
	if ed == "" {
		fmt.Printf("%s:%d\n", file, line)
		return nil
	}

	var cmd *exec.Cmd

	switch {
	case strings.Contains(ed, "code"):
		cmd = exec.Command(ed, "-g", fmt.Sprintf("%s:%d", file, line))
	case strings.Contains(ed, "subl"):
		cmd = exec.Command(ed, fmt.Sprintf("%s:%d", file, line))
	default:
		cmd = exec.Command(ed, fmt.Sprintf("+%d", line), file)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
