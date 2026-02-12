package main

import (
	"log"

	"martinshaw.co/ejecting/diskutil"
	"martinshaw.co/ejecting/lsof"
	"martinshaw.co/ejecting/output"
	"martinshaw.co/ejecting/utilities"
)

// Ignore fucking darwinkit

func main() {
	if !utilities.IsMacOs() {
		log.Fatal("This application is only supported on macOS.")
	}

	disks, error := diskutil.GetDisks()
	if error != nil {
		log.Fatal("Error retrieving disks:", error)
	}

	for _, disk := range disks {
		output.PrintDiskInfo(disk)

		openFiles, error := lsof.GetOpenFilesByDiskMountPrefix(disk.MountPoint)
		if error != nil {
			log.Printf("Warning: Error retrieving open files for disk %s: %v", disk.MountPoint, error)
			continue
		}

		for _, openFile := range openFiles {
			output.PrintOpenFileInfo(openFile)
		}
	}
}
