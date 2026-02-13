package data

import (
	"log"

	"martinshaw.co/ejecting/diskutil"
	"martinshaw.co/ejecting/lsof"
	"martinshaw.co/ejecting/structs"
)

func DetermineData() *structs.DisksWithOpenFiles {
	disks, error := diskutil.GetDisks()
	data := make(structs.DisksWithOpenFiles, 0, len(*disks))
	if error != nil {
		log.Fatal("Error retrieving disks:", error)
	}

	for _, disk := range *disks {
		openFiles, error := lsof.GetOpenFilesByDiskMountPrefix(disk.MountPoint)
		if error != nil {
			log.Printf("Warning: Error retrieving open files for disk %s: %v (usually means no open files)", disk.MountPoint, error)
			continue
		}

		data = append(data, structs.DiskWithOpenFiles{
			Disk:      disk,
			OpenFiles: *openFiles,
		})
	}

	return &data
}
