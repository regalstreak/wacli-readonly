package main

import (
	"os"
	"strings"

	"go.mau.fi/whatsmeow/proto/waCompanionReg"
	"go.mau.fi/whatsmeow/store"
	"google.golang.org/protobuf/proto"
)

func main() {
	applyDeviceLabel()
	if err := execute(os.Args[1:]); err != nil {
		os.Exit(1)
	}
}

func applyDeviceLabel() {
	label := strings.TrimSpace(os.Getenv("WACLI_DEVICE_LABEL"))
	platformRaw := strings.TrimSpace(os.Getenv("WACLI_DEVICE_PLATFORM"))

	// Always set platform type - defaults to Chrome if not specified
	platform := parsePlatformType(platformRaw)
	store.DeviceProps.PlatformType = platform.Enum()

	// Set OS info to look like Chrome browser
	if label == "" {
		label = "Chrome"
	}
	store.SetOSInfo(label, [3]uint32{10, 0, 0})
	store.DeviceProps.Os = proto.String(label)
}

func parsePlatformType(raw string) waCompanionReg.DeviceProps_PlatformType {
	value := strings.TrimSpace(raw)
	if value == "" {
		return waCompanionReg.DeviceProps_CHROME
	}
	value = strings.ToUpper(value)
	if enumValue, ok := waCompanionReg.DeviceProps_PlatformType_value[value]; ok {
		return waCompanionReg.DeviceProps_PlatformType(enumValue)
	}
	return waCompanionReg.DeviceProps_CHROME
}
