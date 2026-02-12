package structs

type OpenFile struct {
	CommandName string
	CommandPath string
	PID         int
	User        string
	Name        string
}

type DiskPartition struct {
	DeviceIdentifier string `plist:"DeviceIdentifier"`
	MountPoint       string `plist:"MountPoint"`
	Size             int    `plist:"Size"`
	VolumeName       string `plist:"VolumeName"`
}

type DiskWithOpenFiles struct {
	Disk      DiskPartition
	OpenFiles []OpenFile
}

type DisksWithOpenFiles []DiskWithOpenFiles
