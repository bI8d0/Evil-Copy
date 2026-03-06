//go:build windows
// +build windows

package main

import (
	"golang.org/x/sys/windows"
)

func isRemovableDrive(drive string) bool {
	driveType := windows.GetDriveType(windows.StringToUTF16Ptr(drive))
	return driveType == windows.DRIVE_REMOVABLE
}
