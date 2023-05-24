package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"
)

func s04(slide int) {
	var (
		params = conf.P[strconv.Itoa(abs(slide))]
		imageBackground,
		visualContainerHost proto.PageViewport
		tit string
	)
	stdo.Println(params)
	exp := func(x interface{}) {
		e(slide, 14, x.(error))
	}
	browser, ca := chrome()
	defer ca()
	page := browser.WithPanic(exp).MustPage().MustSetViewport(1920, 1080, 1, false)
	defer page.Close()
	page.Navigate(params[0])
	time.Sleep(sec)
	page = page.WithPanic(exp).MustWaitLoad()
	tit = page.MustInfo().Title
	scs(slide, page, fmt.Sprintf("%02d %s.png", slide, tit))

	cb(slide, page, "СЦ")

	ex(slide, getClientRect(page.MustElement("div.imageBackground"), &imageBackground))
	ex(slide, getClientRect(page.MustElement("div.visualContainerHost"), &visualContainerHost))

	bytes, err := page.Screenshot(false, &proto.PageCaptureScreenshot{
		Format: proto.PageCaptureScreenshotFormatJpeg,
		Clip: clip(
			visualContainerHost.X,
			visualContainerHost.Y,
			imageBackground.X-visualContainerHost.X,
			visualContainerHost.Height,
		),
	})
	ex(slide, err)
	ss(bytes).write(fmt.Sprintf("%02d.jpg", slide))
	done(slide)
}

func cb(slide int, page *rod.Page, key string) {
	tit := "СЦ"
	se := fmt.Sprintf("div[aria-label=%q] > i", key)
	page.Timeout(to).MustElement(se).MustClick()
	scs(slide, page, fmt.Sprintf("%02d %s.png", slide, tit))

	tit = "Поиск"
	sel := "div.searchHeader.show > input"
	page.Timeout(to * 2).MustElement(sel).Input(sc)
	page.Keyboard.Press(input.Enter)
	scs(slide, page, fmt.Sprintf("%02d %s.png", slide, tit))

	tit = sc
	// sel = fmt.Sprintf("//span[.=%q]", tit)
	// page.Timeout(to * 2).MustElementX(sel).MustClick()
	page.Timeout(to*2).MustElementR("span", tit).MustClick()
	page.Timeout(to).MustElement(se).MustClick()
	scs(slide, page, fmt.Sprintf("%02d %s.png", slide, tit))

	sel = "div.circle"
	WaitElementsLessThan(page.Timeout(to*3), sel, 1)
	time.Sleep(sec)
}
