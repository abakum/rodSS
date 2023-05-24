package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

func s01(slide int) {
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
		page = browser.WithPanic(exp).MustPage().WithPanic(exp).MustSetViewport(1920, 1080, 1, false).MustNavigate(params[0])
		defer page.Close()
		we = page.Timeout(sec).MustElement("div > table.weather")
	} else {
		page, err = browser.Page(proto.TargetCreateTarget{})
		ex(slide, err)
		ex(slide, page.MustSetViewport(1920, 1080, 1, false).Navigate(params[0]))
		we, err = page.Timeout(sec).Element("div > table.weather")
		ex(slide, err)
	}
	time.Sleep(sec)
	bytes, err := we.Screenshot(proto.PageCaptureScreenshotFormatJpeg, 99)
	ex(slide, err)
	ss(bytes).write(fmt.Sprintf("%02d.jpg", slide))
	done(slide)
}
