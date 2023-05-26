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
	br, ca := chrome(slide)
	defer ca()
	// params[0] = "http://ya.ru"
	if wp {
		page = br.MustPage().MustSetViewport(1920, 1080, 1, false)
		defer page.MustClose()
		page.Navigate(params[0])
		we = page.WithPanic(exp).Timeout(sec * 3).MustElement("div > table.weather")
	} else {
		page, err = br.Page(proto.TargetCreateTarget{})
		ex(slide, err)
		defer page.Close()
		page = page.MustSetViewport(1920, 1080, 1, false)
		page.Navigate(params[0])
		we, err = page.Timeout(sec * 3).Element("div > table.weather")
		ex(slide, err)
	}
	time.Sleep(ms)
	bytes, err := we.Screenshot(proto.PageCaptureScreenshotFormatJpeg, 99)
	ex(slide, err)
	ss(bytes).write(fmt.Sprintf("%02d.jpg", slide))
}
