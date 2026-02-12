package output

import (
	"encoding/json"
	"flag"
	"log"

	"howett.net/plist"
)

func ParseFlags() (*string, *string) {
	formatFlag := flag.String("format", "indent", "Output format: 'indent', 'json' or 'plist-openstep', 'plist-xml', 'plist-binary'")
	uiFlag := flag.String("ui", "cli", "User interface: 'cli' or 'menubar'")

	flag.Parse()

	if *formatFlag != "indent" && *formatFlag != "json" && *formatFlag != "plist-openstep" && *formatFlag != "plist-xml" && *formatFlag != "plist-binary" {
		log.Fatal("Invalid format specified. Use 'indent', 'json', 'plist-openstep', 'plist-xml' or 'plist-binary'.")
	}
	if *uiFlag != "cli" && *uiFlag != "menubar" {
		log.Fatal("Invalid UI specified. Use 'cli' or 'menubar'.")
	}

	return formatFlag, uiFlag
}

func PrintDataByFormat(formatFlag *string, data interface{}) {
	switch *formatFlag {
	case "indent":
		PrintDataInIndentedFormat(data)
	case "json":
		PrintDataInJsonFormat(data)
	case "plist-openstep":
		PrintDataInPlistFormat(data, plist.OpenStepFormat)
	case "plist-xml":
		PrintDataInPlistFormat(data, plist.XMLFormat)
	case "plist-binary":
		PrintDataInPlistFormat(data, plist.BinaryFormat)
	default:
		log.Fatal("Invalid format specified. Use 'indent', 'json' or 'plist'.")
	}
}

func PrintDataInIndentedFormat(data interface{}) {
	jsonString, error := json.MarshalIndent(data, "", "  ")
	if error != nil {
		log.Fatal("Error marshaling data to JSON:", error)
	}
	println(string(jsonString))
}

func PrintDataInJsonFormat(data interface{}) {
	jsonString, error := json.Marshal(data)
	if error != nil {
		log.Fatal("Error marshaling data to JSON:", error)
	}
	println(string(jsonString))
}

func PrintDataInPlistFormat(data interface{}, format int) {
	plistString, error := plist.Marshal(data, format)
	if error != nil {
		log.Fatal("Error marshaling data to plist:", error)
	}
	println(string(plistString))
}
