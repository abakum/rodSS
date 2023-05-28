package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"
)

func s04(slide, deb int) {
	var (
		params = conf.P[strconv.Itoa(abs(slide))]
		imageBackground,
		visualContainerHost proto.PageViewport
	)
	stdo.Println(params)
	br, ca := chrome(slide)
	defer ca()
	page := chromePage(br, slide)
	defer page.MustClose()
	page.Navigate(params[0])
	time.Sleep(sec)
	tit := page.MustWaitLoad().MustInfo().Title
	sdpt(slide, deb, page, tit)

	cb(slide, deb, page, "СЦ")

	sel := "div.imageBackground"
	ex(slide, getClientRect(page.MustElement(sel), &imageBackground))

	sel = "div.visualContainerHost"
	ex(slide, getClientRect(page.MustElement(sel), &visualContainerHost))

	sl(slide).done(page.Screenshot(false, &proto.PageCaptureScreenshot{
		Format: proto.PageCaptureScreenshotFormatJpeg,
		Clip: clip(
			visualContainerHost.X,
			visualContainerHost.Y,
			imageBackground.X-visualContainerHost.X,
			visualContainerHost.Height,
		),
	}))
}

func cb(slide, deb int, page *rod.Page, key string) {
	tit := "СЦ"
	se := fmt.Sprintf("div[aria-label=%q] > i", key)
	page.Timeout(to).MustElement(se).MustClick()
	sdpt(slide, deb, page, tit)

	tit = "Поиск"
	sel := "div.searchHeader.show > input"
	page.Timeout(to * 2).MustElement(sel).Input(sc)
	page.Keyboard.Press(input.Enter)
	sdpt(slide, deb, page, tit)

	tit = sc
	sel = "span"
	page.Timeout(to*2).MustElementR(sel, tit).MustClick()
	// page.Keyboard.MustType(input.Tab)
	page.Timeout(to).MustElement(se).MustClick()

	sel = "div.circle"
	WaitElementsLessThan(page.Timeout(to*3), sel, 1)
	// time.Sleep(ms)
	sdpt(slide, deb, page, tit)
}
