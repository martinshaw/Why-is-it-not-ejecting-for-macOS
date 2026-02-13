package darwinkit

import (
	"runtime"
	"strconv"
	"sync"

	emersion "github.com/emersion/go-autostart"
	"github.com/progrium/darwinkit/dispatch"
	"github.com/progrium/darwinkit/macos/appkit"
	"github.com/progrium/darwinkit/macos/foundation"
	"github.com/progrium/darwinkit/objc"
	"github.com/sanity-io/litter"
	"martinshaw.co/ejecting/autostart"
	"martinshaw.co/ejecting/data"
	"martinshaw.co/ejecting/diskutil"
	"martinshaw.co/ejecting/ps"
)

const appName = "Why is it not ejecting? for macOS"
const shortAppName = "Why is it not ejecting?"

var refreshMutex sync.Mutex
var statusItem appkit.StatusItem

func StartMenubarUi(autostartApp *emersion.App) {
	runtime.LockOSThread()

	app := appkit.Application_SharedApplication()
	app.SetActivationPolicy(appkit.ApplicationActivationPolicyAccessory)

	delegate := &appkit.ApplicationDelegate{}
	delegate.SetApplicationDidFinishLaunching(func(foundation.Notification) {
		setSystemBar(&app, autostartApp)
	})
	delegate.SetApplicationShouldTerminateAfterLastWindowClosed(func(appkit.Application) bool {
		return false
	})

	delegate.SetApplicationDidBecomeActive(func(foundation.Notification) {
		litter.Dump("SetApplicationDidBecomeActive")
	})

	delegate.SetApplicationWillTerminate(func(foundation.Notification) {
		litter.Dump("SetApplicationWillTerminate")
	})

	app.SetDelegate(delegate)
	app.Run()
}

func refreshMenuWithData(menu *appkit.Menu, app *appkit.Application, autostartApp *emersion.App) {
	if !refreshMutex.TryLock() {
		return
	}

	go func() {
		defer refreshMutex.Unlock()

		latestData := *data.DetermineData()

		dispatch.MainQueue().DispatchAsync(func() {
			menu.RemoveAllItems()

			if len(latestData) == 0 {
				menu.AddItem(appkit.NewMenuItemWithAction("No issues detected. Ejection should work fine.", "", func(sender objc.Object) {}))
			}

			for _, diskWithOpenFiles := range latestData {
				diskInfo := "Disk: " + diskWithOpenFiles.Disk.DeviceIdentifier + " (" + diskWithOpenFiles.Disk.VolumeName + " at " + diskWithOpenFiles.Disk.MountPoint + ") - Click to kill all and eject"
				menu.AddItem(appkit.NewMenuItemWithAction(diskInfo, "", func(sender objc.Object) {
					for _, openFile := range diskWithOpenFiles.OpenFiles {
						ps.KillProcessByPid(openFile.PID)
					}

					diskutil.EjectDiskByIdentifier(&diskWithOpenFiles.Disk)
					refreshMenuWithData(menu, app, autostartApp)
				}))

				for _, openFile := range diskWithOpenFiles.OpenFiles {
					openFileInfo := " - " + openFile.Name + " (PID " + strconv.Itoa(openFile.PID) + ", Process: " + openFile.CommandName + ") - Click to kill process"
					menu.AddItem(appkit.NewMenuItemWithAction(openFileInfo, "", func(sender objc.Object) {
						ps.KillProcessByPid(openFile.PID)
						refreshMenuWithData(menu, app, autostartApp)
					}))
				}
			}

			menu.AddItem(appkit.MenuItem_SeparatorItem())

			menu.AddItem(appkit.NewMenuItemWithAction("Refresh", "r", func(sender objc.Object) {
				refreshMenuWithData(menu, app, autostartApp)
			}))
			openOnLoginMenuItem := appkit.NewMenuItemWithAction("Open on Login", "l", func(sender objc.Object) {
				autostart.ToggleAutostartApp(autostartApp)
				refreshMenuWithData(menu, app, autostartApp)
			})
			var openOnLoginMenuItemState appkit.ControlStateValue
			if autostartApp.IsEnabled() {
				openOnLoginMenuItemState = appkit.ControlStateValueOn
			} else {
				openOnLoginMenuItemState = appkit.ControlStateValueOff
			}
			openOnLoginMenuItem.SetState(openOnLoginMenuItemState)
			menu.AddItem(openOnLoginMenuItem)
			menu.AddItem(appkit.NewMenuItemWithAction("Quit", "q", func(sender objc.Object) { app.Terminate(nil) }))
		})
	}()
}

func setSystemBar(app *appkit.Application, autostart *emersion.App) {
	statusItem = appkit.StatusBar_SystemStatusBar().StatusItemWithLength(appkit.VariableStatusItemLength)
	objc.Retain(&statusItem)

	img := appkit.Image_ImageWithSystemSymbolNameAccessibilityDescription("eject.circle.fill", "An eject symbol inside a filled circle.")
	statusItem.Button().SetImage(img)
	statusItem.Button().SetToolTip(appName)

	menu := appkit.NewMenuWithTitle(shortAppName)
	statusItem.SetMenu(menu)

	refreshMenuWithData(&menu, app, autostart)
}
