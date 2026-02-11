package output

import "martinshaw.co/ejecting/df"

func PrintDiskInfo(disk df.Disk) {
	println("Disk:", disk.Filesystem, "Mounted on:", disk.MountPoint, "External:", disk.IsExternal)
}
