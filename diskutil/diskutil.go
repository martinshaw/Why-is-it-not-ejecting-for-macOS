package diskutil

import (
	"os/exec"
	"strings"

	"howett.net/plist"
	"martinshaw.co/ejecting/structs"
)

func executeCommandForDiskInfo() (string, error) {
	cmd := exec.Command("diskutil", "list", "-plist")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func decodePlistOutputToDisks(plistOutput string) ([]structs.DiskPartition, error) {
	decoder := plist.NewDecoder(strings.NewReader(plistOutput))

	var data struct {
		AllDisksAndPartitions []struct {
			Partitions []structs.DiskPartition `plist:"Partitions"`
		} `plist:"AllDisksAndPartitions"`
	}

	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}

	disks := make([]structs.DiskPartition, 0)
	for _, diskInfo := range data.AllDisksAndPartitions {
		for _, partition := range diskInfo.Partitions {
			if partition.MountPoint == "" ||
				partition.MountPoint == "/" ||
				partition.DeviceIdentifier == "" ||
				partition.Size == 0 {
				continue
			}
			disks = append(disks, partition)
		}
	}

	return disks, nil
}

func GetDisks() ([]structs.DiskPartition, error) {
	output, error := executeCommandForDiskInfo()
	if error != nil {
		return nil, error
	}

	disks, error := decodePlistOutputToDisks(output)
	if error != nil {
		return nil, error
	}

	return disks, nil
}
