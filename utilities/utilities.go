package utilities

import (
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func IsMacOs() bool {
	return runtime.GOOS == "darwin"
}

func GetLineCountOfOutput(output *string) int {
	cmd := exec.Command("wc", "-l")
	cmd.Stdin = strings.NewReader(*output)
	lineCountOutput, err := cmd.Output()
	if err != nil {
		return 0
	}
	lineCountStr := strings.TrimSpace(string(lineCountOutput))
	lineCount, err := strconv.Atoi(lineCountStr)
	if err != nil {
		return 0
	}
	return lineCount
}

func GetSystemUsername() string {
	cmd := exec.Command("whoami")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// Get the absolute path of the currently running binary
func GetBinaryPath() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return execPath, nil
}
