package main

import (
	"os"
	"os/exec"
	"path/filepath"
)

func s97(slide, deb int) {
	cmd := exec.Command("cmd", "/c", filepath.Join(root, bat))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	ex(slide, cmd.Run())
}
