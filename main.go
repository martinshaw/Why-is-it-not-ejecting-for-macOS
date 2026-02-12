package main

import (
	"log"
	"os"

	"martinshaw.co/ejecting/darwinkit"
	"martinshaw.co/ejecting/data"
	"martinshaw.co/ejecting/output"
	"martinshaw.co/ejecting/utilities"
)

func main() {
	formatFlag, uiFlag := output.ParseFlags()

	if !utilities.IsMacOs() {
		log.Fatal("This application is only supported on macOS.")
	}

	if uiFlag != nil && *uiFlag == "cli" {
		data := data.DetermineData()
		output.PrintDataByFormat(formatFlag, data)
		os.Exit(0)
	}

	if uiFlag != nil && *uiFlag == "menubar" {
		darwinkit.StartMenubarUi()
	}
}
