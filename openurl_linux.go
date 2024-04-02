//go:build linux
// +build linux

package main

import (
	"os/exec"
)

func openUrl(url string) {
	cmd := exec.Command("xdg-open", url)
	cmd.Run()
}
