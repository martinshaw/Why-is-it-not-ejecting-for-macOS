package darwinkit

import (
	"strconv"

	"github.com/progrium/darwinkit/macos/appkit"
	"github.com/progrium/darwinkit/macos/foundation"
	"github.com/progrium/darwinkit/objc"
	"martinshaw.co/ejecting/data"
	"martinshaw.co/ejecting/diskutil"
	"martinshaw.co/ejecting/ps"
	"martinshaw.co/ejecting/structs"
)

const appName = "Why is it not ejecting? for macOS"
const shortAppName = "Why is it not ejecting?"

func StartMenubarUi() {
	app := appkit.Application_SharedApplication()
	app.SetActivationPolicy(appkit.ApplicationActivationPolicyAccessory)

	delegate := &appkit.ApplicationDelegate{}
	delegate.SetApplicationDidFinishLaunching(func(foundation.Notification) {
		setSystemBar(app)
	})
	delegate.SetApplicationShouldTerminateAfterLastWindowClosed(func(appkit.Application) bool {
		return false
	})
	app.SetDelegate(delegate)
	app.Run()
}

func refreshMenuWithData(menu *appkit.Menu, latestData *structs.DisksWithOpenFiles, app *appkit.Application) {
	latestData = data.DetermineData()

	menu.RemoveAllItems()

	if len(*latestData) == 0 {
		menu.AddItem(appkit.NewMenuItemWithAction("No issues detected. Ejection should work fine.", "", func(sender objc.Object) {}))
	}

	for _, diskWithOpenFiles := range *latestData {
		diskInfo := "Disk: " + diskWithOpenFiles.Disk.DeviceIdentifier + " (" + diskWithOpenFiles.Disk.VolumeName + " at " + diskWithOpenFiles.Disk.MountPoint + ") - Click to kill all and eject"
		menu.AddItem(appkit.NewMenuItemWithAction(diskInfo, "", func(sender objc.Object) {
			for _, openFile := range diskWithOpenFiles.OpenFiles {
				ps.KillProcessByPid(openFile.PID)
			}

			diskutil.EjectDiskByIdentifier(&diskWithOpenFiles.Disk)
			refreshMenuWithData(menu, latestData, app)
		}))

		for _, openFile := range diskWithOpenFiles.OpenFiles {
			openFileInfo := " - " + openFile.Name + " (PID " + strconv.Itoa(openFile.PID) + ", Process: " + openFile.CommandName + ") - Click to kill process"
			menu.AddItem(appkit.NewMenuItemWithAction(openFileInfo, "", func(sender objc.Object) {
				ps.KillProcessByPid(openFile.PID)
				refreshMenuWithData(menu, latestData, app)
			}))
		}
	}

	menu.AddItem(appkit.MenuItem_SeparatorItem())

	menu.AddItem(appkit.NewMenuItemWithAction("Refresh", "r", func(sender objc.Object) {
		refreshMenuWithData(menu, latestData, app)
	}))
	menu.AddItem(appkit.NewMenuItemWithAction("Quit", "q", func(sender objc.Object) { app.Terminate(nil) }))
}

func setSystemBar(app appkit.Application) {
	item := appkit.StatusBar_SystemStatusBar().StatusItemWithLength(appkit.VariableStatusItemLength)
	objc.Retain(&item)

	img := appkit.Image_ImageWithSystemSymbolNameAccessibilityDescription("eject.circle.fill", "An eject symbol inside a filled circle.")
	item.Button().SetImage(img)
	item.Button().SetToolTip(appName)

	menu := appkit.NewMenuWithTitle(shortAppName)
	item.SetMenu(menu)

	var latestData structs.DisksWithOpenFiles
	refreshMenuWithData(&menu, &latestData, &app)
}
