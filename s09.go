package main

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
)

func s09(slide int) {
	err := exec.Command("cmd", "/c", "start", filepath.Join(cd, "imager.xlsb")).Run()
	ex(slide, err)
	for _, v := range []int{2, 7, 8, 9} {
		_, err = os.Stat(i2p(v))
		if errors.Is(err, os.ErrNotExist) {
			ex(slide, err)
			return
		}
	}
}
