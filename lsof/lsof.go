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
	CommandName string
	CommandPath string
	PID         int
	User        string
	Name        string
}

func executeLsofCommandForOpenFilesByDiskMountPrefix(mountPathPrefix string) (string, error) {
	if !strings.HasSuffix(mountPathPrefix, "/") {
		mountPathPrefix += "/"
	}

	// The +c 0 option tells lsof to not truncate the output
	// The -u option filters the output to only include files opened by the current user
	prompt := fmt.Sprintf("lsof -u %s +c 0 | grep \"%s\"", utilities.GetSystemUsername(), mountPathPrefix)
	cmd := exec.Command("sh", "-c", prompt)

	output, error := cmd.Output()
	if error != nil {
		return "", error
	}

	return string(output), nil
}

func executePsCommandByPid(pid int) (string, error) {
	// The -o comm= option tells ps to only output the command path of the process, without any headers or additional formatting
	// The -p option filters the output to only include the process with the specified PID
	prompt := fmt.Sprintf("ps -p %d -o comm=", pid)
	cmd := exec.Command("sh", "-c", prompt)

	output, error := cmd.Output()
	if error != nil {
		return "", error
	}

	return string(output), nil
}

func parseLineOfOutput(line string) (OpenFile, error) {
	fields := strings.Fields(line)
	if len(fields) < 9 {
		return OpenFile{}, fmt.Errorf("unexpected output format: %s", line)
	}

	pid, error := strconv.Atoi(fields[1])
	if error != nil {
		return OpenFile{}, fmt.Errorf("error parsing PID for %s: %v", fields[0], error)
	}

	commandPath, error := executePsCommandByPid(pid)
	if error != nil {
		return OpenFile{}, fmt.Errorf("error executing ps command for PID %d: %v", pid, error)
	}

	name := strings.Join(fields[8:], " ")

	return OpenFile{
		CommandName: fields[0],
		CommandPath: strings.TrimSpace(commandPath),
		PID:         pid,
		User:        fields[2],
		Name:        name,
	}, nil
}

func GetOpenFilesByDiskMountPrefix(mountPathPrefix string) ([]OpenFile, error) {
	output, error := executeLsofCommandForOpenFilesByDiskMountPrefix(mountPathPrefix)
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
