package main

import (
	"log"

	"github.com/sanity-io/litter"
	"martinshaw.co/ejecting/df"
	"martinshaw.co/ejecting/lsof"
	"martinshaw.co/ejecting/output"
	"martinshaw.co/ejecting/utilities"
)

// Ignore fucking darwinkit

func main() {
	if !utilities.IsMacOs() {
		log.Fatal("This application is only supported on macOS.")
	}

	disks, error := df.GetDisks()
	if error != nil {
		log.Fatal("Error retrieving disks:", error)
	}

	for _, disk := range disks {
		if !disk.IsExternal {
			continue
		}
		output.PrintDiskInfo(disk)

		litter.Dump("zzz", disk)

		openFiles, error := lsof.GetOpenFilesByDiskMountPrefix(disk.MountPoint)
		if error != nil {
			log.Printf("Warning: Error retrieving open files for disk %s: %v", disk.MountPoint, error)
			continue
		}

		litter.Dump(openFiles)
	}
}
