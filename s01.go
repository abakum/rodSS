package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

func s01(slide, deb int) {
	var (
		params = conf.P[strconv.Itoa(abs(slide))]
		wp     = true
		we     *rod.Element
		page   *rod.Page
		err    error
	)
	stdo.Println(params)
	exp := func(x interface{}) {
		e(slide, 14, x.(error))
	}
	browser, ca := chrome()
	defer ca()
	if wp {
		page = browser.WithPanic(exp).MustPage().MustSetViewport(1920, 1080, 1, false)
		defer page.Close()
		page.Navigate(params[0])
		we = page.WithPanic(exp).Timeout(to).MustElement("div > table.weather").WithPanic(exp)
	} else {
		page, err = browser.Page(proto.TargetCreateTarget{})
		ex(slide, err)
		defer page.Close()
		page = page.MustSetViewport(1920, 1080, 1, false)
		defer page.Close()
		page.Navigate(params[0])
		we, err = page.Timeout(to).Element("div > table.weather")
		ex(slide, err)
	}
	time.Sleep(ms)
	bytes, err := we.Screenshot(proto.PageCaptureScreenshotFormatJpeg, 99)
	ex(slide, err)
	ss(bytes).write(fmt.Sprintf("%02d.jpg", slide))
}
