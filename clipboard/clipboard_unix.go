//go:build !windows
package clipboard

import (
	"crypto/md5"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

var lastHash string

func Read() (string, error) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbpaste")
	case "linux":
		if _, err := exec.LookPath("xclip"); err == nil {
			cmd = exec.Command("xclip", "-selection", "clipboard", "-o")
		} else {
			cmd = exec.Command("xsel", "--clipboard", "--output")
		}
	default:
		return "", nil
	}

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

// HasChanged returns true if clipboard changed since last check
func HasChanged() bool {
	content, err := Read()
	if err != nil {
		return false
	}
	
	currentHash := fmt.Sprintf("%x", md5.Sum([]byte(content)))
	
	if currentHash != lastHash {
		lastHash = currentHash
		return true
	}
	
	return false
}
