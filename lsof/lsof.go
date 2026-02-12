package lsof

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"martinshaw.co/ejecting/ps"
	"martinshaw.co/ejecting/structs"
	"martinshaw.co/ejecting/utilities"
)

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
		// When there is an error response, it is usually simply that there are no open files for the specified disk, so there is nothing to grep. This isn't really an error and can be ignored as expected default behavior
		return "", error
	}

	return string(output), nil
}

func parseLineOfOutput(line string) (structs.OpenFile, error) {
	fields := strings.Fields(line)
	if len(fields) < 9 {
		return structs.OpenFile{}, fmt.Errorf("unexpected output format: %s", line)
	}

	pid, error := strconv.Atoi(fields[1])
	if error != nil {
		return structs.OpenFile{}, fmt.Errorf("error parsing PID for %s: %v", fields[0], error)
	}

	commandPath, error := ps.ExecutePsCommandByPid(pid)
	if error != nil {
		return structs.OpenFile{}, fmt.Errorf("error executing ps command for PID %d: %v", pid, error)
	}

	name := strings.Join(fields[8:], " ")

	return structs.OpenFile{
		CommandName: fields[0],
		CommandPath: strings.TrimSpace(commandPath),
		PID:         pid,
		User:        fields[2],
		Name:        name,
	}, nil
}

func GetOpenFilesByDiskMountPrefix(mountPathPrefix string) ([]structs.OpenFile, error) {
	output, error := executeLsofCommandForOpenFilesByDiskMountPrefix(mountPathPrefix)
	if error != nil {
		return nil, error
	}

	openFilesCount := utilities.GetLineCountOfOutput(output)
	openFiles := make([]structs.OpenFile, 0, openFilesCount)

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		if openFile, error := parseLineOfOutput(scanner.Text()); error == nil {
			openFiles = append(openFiles, openFile)
		} else {
			fmt.Printf("Warning: Error parsing line of lsof output: %v\n", error)
			continue
		}
	}

	return openFiles, nil
}
