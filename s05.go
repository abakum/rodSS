package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-rod/rod/lib/proto"
)

func s05(slide, deb int) {
	var (
		params = conf.P[strconv.Itoa(abs(slide))]
		title,
		innerContainer,
		visualContainerHost proto.PageViewport
	)
	stdo.Println(params, sc)
	exp := func(x interface{}) {
		e(slide, 14, x.(error))
	}
	br, ca := chrome(slide)
	defer ca()
	page := br.MustPage().WithPanic(exp).MustSetViewport(1920, 1080, 1, false)
	defer page.MustClose()
	page.Navigate(params[0])
	time.Sleep(sec)
	tit := page.MustWaitLoad().MustInfo().Title
	sdpt(slide, deb, page, tit)

	tit = "Статистика по сотрудникам"
	sel := fmt.Sprintf("div[title=%q]", tit)
	page.Timeout(to).MustElement(sel).MustClick()
	sdpt(slide, deb, page, tit)

	cb(slide, deb, page, "СЦ/ЦЭ")

	tit = "Ср. длительность работ сотрудника за день, часы"
	sel = fmt.Sprintf("div[title=%q]", tit)
	ex(slide, getClientRect(page.MustElement(sel), &title))

	sel = "div.innerContainer"
	ex(slide, getClientRect(page.MustElement(sel), &innerContainer))

	sel = "div.visualContainerHost"
	ex(slide, getClientRect(page.MustElement(sel), &visualContainerHost))

	sl(slide).done(page.Screenshot(false, &proto.PageCaptureScreenshot{
		Format: proto.PageCaptureScreenshotFormatJpeg,
		Clip: clip(
			visualContainerHost.X,
			visualContainerHost.Y,
			title.X-visualContainerHost.X,
			innerContainer.Y+innerContainer.Height-2*visualContainerHost.Y,
		),
	}))
}
