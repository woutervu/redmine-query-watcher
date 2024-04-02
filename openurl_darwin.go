//go:build darwin
// +build darwin

package main

import (
	"os/exec"
)

func openUrl(url string) {
	cmd := exec.Command("open", url)
	cmd.Run()
}
