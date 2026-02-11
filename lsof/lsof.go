package lsof

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"martinshaw.co/ejecting/utilities"
)

type OpenFile struct {
	Command string
	PID     int
	User    string
	Name    string
}

func executeCommandForOpenFilesByDiskMountPrefix(mountPathPrefix string) (string, error) {
	cmd := exec.Command("lsof", "+D", mountPathPrefix)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func parseLineOfOutput(line string) (OpenFile, error) {
	fields := strings.FieldsFunc(line, func(r rune) bool {
		return r == '\t'
	})
	if len(fields) < 9 {
		return OpenFile{}, fmt.Errorf("unexpected output format: %s", line)
	}

	pid, error := strconv.Atoi(fields[1])
	if error != nil {
		return OpenFile{}, fmt.Errorf("error parsing PID for %s: %v", fields[0], error)
	}

	return OpenFile{
		Command: fields[0],
		PID:     pid,
		User:    fields[2],
		Name:    fields[8],
	}, nil
}

func GetOpenFilesByDiskMountPrefix(mountPathPrefix string) ([]OpenFile, error) {
	output, error := executeCommandForOpenFilesByDiskMountPrefix(mountPathPrefix)
	if error != nil {
		return nil, error
	}

	openFilesCount := utilities.GetLineCountOfOutput(output) - 1 // Subtract 1 for the header line
	openFiles := make([]OpenFile, 0, openFilesCount)
	isScanningHeaderRow := true

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		if isScanningHeaderRow {
			isScanningHeaderRow = false
			continue
		}

		if openFile, error := parseLineOfOutput(scanner.Text()); error == nil {
			openFiles = append(openFiles, openFile)
		} else {
			fmt.Printf("Warning: Error parsing line of lsof output: %v\n", error)
			continue
		}
	}

	return openFiles, nil
}
