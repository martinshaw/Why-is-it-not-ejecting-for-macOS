package output

import (
	"github.com/dustin/go-humanize"
	"martinshaw.co/ejecting/diskutil"
	"martinshaw.co/ejecting/lsof"
)

func PrintDiskInfo(disk diskutil.DiskPartition) {
	println(
		"Disk:", disk.DeviceIdentifier, "\n",
		"Mounted on:", disk.MountPoint, "\n",
		"Size:", humanize.Bytes(uint64(disk.Size)), "\n",
		"Size (Bytes):", disk.Size, "\n",
		"Volume Name:", disk.VolumeName, "\n",
	)
}

func PrintOpenFileInfo(openFile lsof.OpenFile) {
	println(
		"	Command Name:", openFile.CommandName, "\n",
		"	Command Path:", openFile.CommandPath, "\n",
		"	PID:", openFile.PID, "\n",
		"	User:", openFile.User, "\n",
		"	File Name:", openFile.Name, "\n",
	)
}
