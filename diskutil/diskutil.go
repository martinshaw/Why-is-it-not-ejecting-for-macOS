package diskutil

import (
	"os/exec"
	"strings"

	"howett.net/plist"
)

type DiskPartition struct {
	DeviceIdentifier string `plist:"DeviceIdentifier"`
	MountPoint       string `plist:"MountPoint"`
	Size             int    `plist:"Size"`
	VolumeName       string `plist:"VolumeName"`
}

func executeCommandForDiskInfo() (string, error) {
	cmd := exec.Command("diskutil", "list", "-plist")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func decodePlistOutputToDisks(plistOutput string) ([]DiskPartition, error) {
	decoder := plist.NewDecoder(strings.NewReader(plistOutput))

	var data struct {
		AllDisksAndPartitions []struct {
			Partitions []DiskPartition `plist:"Partitions"`
		} `plist:"AllDisksAndPartitions"`
	}

	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}

	disks := make([]DiskPartition, 0)
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

func GetDisks() ([]DiskPartition, error) {
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
