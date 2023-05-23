package main

import (
	"fmt"
	"strconv"

	"github.com/go-rod/rod/lib/proto"
)

func s01(slide int) {
	var (
		params = conf.P[strconv.Itoa(abs(slide))]
	)
	stdo.Println(params)
	browser, ca := chrome()
	defer ca()
	page := browser.MustPage().MustSetViewport(1920, 1080, 1, false)
	defer page.Close()
	ex(slide, page.Navigate(params[0]))
	we, err := page.Timeout(to).Element("div > table.weather")
	ex(slide, err)
	bytes, err := we.Screenshot(proto.PageCaptureScreenshotFormatJpeg, 99)
	ex(slide, err)
	ss(bytes).write(fmt.Sprintf("%02d.jpg", slide))
	done(slide)
}
