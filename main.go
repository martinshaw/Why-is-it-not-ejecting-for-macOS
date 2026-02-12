package main

import (
	"log"

	"martinshaw.co/ejecting/diskutil"
	"martinshaw.co/ejecting/lsof"
	"martinshaw.co/ejecting/output"
	"martinshaw.co/ejecting/structs"
	"martinshaw.co/ejecting/utilities"
)

// Ignore fucking darwinkit

func main() {
	formatFlag := output.ParseFlags()

	if !utilities.IsMacOs() {
		log.Fatal("This application is only supported on macOS.")
	}

	disks, error := diskutil.GetDisks()
	data := make(structs.DisksWithOpenFiles, 0, len(disks))
	if error != nil {
		log.Fatal("Error retrieving disks:", error)
	}

	for _, disk := range disks {
		openFiles, error := lsof.GetOpenFilesByDiskMountPrefix(disk.MountPoint)
		if error != nil {
			log.Printf("Warning: Error retrieving open files for disk %s: %v", disk.MountPoint, error)
			continue
		}

		data = append(data, structs.DiskWithOpenFiles{
			Disk:      disk,
			OpenFiles: openFiles,
		})
	}

	output.PrintDataByFormat(formatFlag, data)
}
