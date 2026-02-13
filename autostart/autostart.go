package autostart

import (
	"log"

	"github.com/emersion/go-autostart"
	"github.com/sanity-io/litter"
	"martinshaw.co/ejecting/utilities"
)

func GetAutostartApp() *autostart.App {
	binaryPath, err := utilities.GetBinaryPath()
	if err != nil {
		litter.Dump("Error getting binary path", err)
	}

	command := binaryPath + " -ui menubar"

	log.Println("Getting status of any autostart entries for \"" + command + "\"")

	return &autostart.App{
		Name:        "Why is it not ejecting? for macOS",
		DisplayName: "Why is it not ejecting? for macOS",
		Exec:        []string{binaryPath},
	}
}

func EnableAutostartApp(app *autostart.App) {
	enabled := app.IsEnabled()
	if !enabled {
		log.Println("Enabling autostart for menubar app...")
		err := app.Enable()
		if err != nil {
			log.Fatal("Error enabling autostart:", err)
		}
	} else {
		log.Println("Autostart for menubar app is already enabled.")
	}
}

func DisableAutostartApp(app *autostart.App) {
	enabled := app.IsEnabled()
	if enabled {
		log.Println("Disabling autostart for menubar app...")
		err := app.Disable()
		if err != nil {
			log.Fatal("Error disabling autostart:", err)
		}
	} else {
		log.Println("Autostart for menubar app is already disabled.")
	}
}

func ToggleAutostartApp(app *autostart.App) {
	enabled := app.IsEnabled()
	if enabled {
		DisableAutostartApp(app)
	} else {
		EnableAutostartApp(app)
	}
}
