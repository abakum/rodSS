package main

import (
	"strconv"

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
	br, ca := chrome(slide)
	defer ca()
	params[0] = "http://ya.ru"
	sel := "div > table.weather"
	if wp {
		page = chromePage(br, slide)
		defer page.MustClose()
		page.Navigate(params[0])
		sl(slide).done(page.Timeout(sec*3).MustElement(sel).
			Screenshot(proto.PageCaptureScreenshotFormatJpeg, 99))
	} else {
		page, err = br.Page(proto.TargetCreateTarget{})
		ex(slide, err)
		defer page.Close()
		page.MustSetViewport(1920, 1080, 1, false)
		page.Navigate(params[0])
		we, err = page.Timeout(sec * 3).Element(sel)
		ex(slide, err)
		sl(slide).done(we.Screenshot(proto.PageCaptureScreenshotFormatJpeg, 99))
	}
}
