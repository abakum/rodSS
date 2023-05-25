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
	br, ca := chrome()
	defer ca()
	page := br.WithPanic(exp).MustPage().MustSetViewport(1920, 1080, 1, false)
	defer page.Close()
	page.Navigate(params[0])
	time.Sleep(sec)
	page = page.WithPanic(exp).MustWaitLoad()
	tit := page.MustInfo().Title
	scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))

	tit = "Статистика по сотрудникам"
	sel := fmt.Sprintf("div[title=%q]", tit)
	page.Timeout(to).MustElement(sel).MustClick()
	scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))

	cb(slide, deb, page, "СЦ/ЦЭ")

	tit = "Ср. длительность работ сотрудника за день, часы"
	sel = fmt.Sprintf("div[title=%q]", tit)
	ex(slide, getClientRect(page.MustElement(sel), &title))

	sel = "div.innerContainer"
	ex(slide, getClientRect(page.MustElement(sel), &innerContainer))

	sel = "div.visualContainerHost"
	ex(slide, getClientRect(page.MustElement(sel), &visualContainerHost))

	bytes, err := page.Screenshot(false, &proto.PageCaptureScreenshot{
		Format: proto.PageCaptureScreenshotFormatJpeg,
		Clip: clip(
			visualContainerHost.X,
			visualContainerHost.Y,
			title.X-visualContainerHost.X,
			innerContainer.Y+innerContainer.Height-2*visualContainerHost.Y,
		),
	})
	ex(slide, err)
	ss(bytes).write(fmt.Sprintf("%02d.jpg", slide))
}
