package main

import (
	"os"
	"os/exec"
	"path/filepath"
)

func s97(slide int) {
	cmd := exec.Command("cmd", "/c", filepath.Join(root, bat))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	ex(slide, cmd.Run())
	done(slide)
}
