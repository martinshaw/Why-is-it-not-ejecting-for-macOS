package df

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/sanity-io/litter"
	"martinshaw.co/ejecting/utilities"
)

type Disk struct {
	Filesystem string
	Used       uint64
	Available  uint64
	Capacity   string
	MountPoint string
	IsExternal bool
}

func executeCommandForDiskInfo() (string, error) {
	cmd := exec.Command("df")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func parseLineOfOutput(line string) (Disk, error) {
	// fields := strings.FieldsFunc(line, func(r rune) bool {
	// 	return r == '\t'
	// })
	fields := strings.Fields(line)
	litter.Dump("yyy", fields)
	os.Exit(0)

	if len(fields) < 9 {
		return Disk{}, fmt.Errorf("unexpected output format: %s", line)
	}

	used, error := strconv.Atoi(fields[2])
	if error != nil {
		return Disk{}, fmt.Errorf("error parsing used space for %s: %v", fields[0], error)
	}

	available, error := strconv.Atoi(fields[3])
	if error != nil {
		return Disk{}, fmt.Errorf("error parsing available space for %s: %v", fields[0], error)
	}

	// isExternal := strings.HasPrefix(fields[0], "/dev/disk") && !strings.Contains(fields[0], "s1")
	isExternal := strings.HasPrefix(fields[8], "/Volumes/") && !strings.Contains(fields[8], "Macintosh HD")

	litter.Dump("xxx", fields)
	return Disk{
		Filesystem: fields[0],
		Used:       uint64(used),
		Available:  uint64(available),
		Capacity:   fields[4],
		MountPoint: fields[8],
		IsExternal: isExternal,
	}, nil
}

func GetDisks() ([]Disk, error) {
	output, error := executeCommandForDiskInfo()
	if error != nil {
		return nil, error
	}

	diskCount := utilities.GetLineCountOfOutput(output) - 1 // Subtract 1 for the header line
	disks := make([]Disk, 0, diskCount)
	isScanningHeaderRow := true

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		if isScanningHeaderRow {
			isScanningHeaderRow = false
			continue
		}

		if disk, err := parseLineOfOutput(scanner.Text()); err == nil {
			disks = append(disks, disk)
		} else {
			fmt.Printf("Warning: Error parsing line of df output: %v\n", err)
			continue
		}
	}

	return disks, nil
}
