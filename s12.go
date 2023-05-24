package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-rod/rod/lib/proto"
)

func s12(slide, deb int) {
	var (
		params = conf.P[strconv.Itoa(abs(slide))]
		vc27,
		vc22,
		vcHost proto.PageViewport
	)
	stdo.Println(params, sc)
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
	tit := page.MustInfo().Title
	scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))

	tit = params[1]
	sel := fmt.Sprintf("div[aria-label=%q]", tit)
	page.Timeout(to * 2).MustElement(sel).MustClick()
	scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))

	tit = "RF"
	sel = fmt.Sprintf("div[aria-label=%q]", tit)
	page.Timeout(to).MustElement(sel).MustClick()
	scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))

	tit = params[2]
	sel = fmt.Sprintf("span[title=%q]", tit)
	page.Timeout(to).MustElement(sel).MustClick()
	scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))

	tit = "SC_NAME"
	sel = fmt.Sprintf("div[aria-label=%q]", tit)
	page.Timeout(to).MustElement(sel).MustClick()
	scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))

	tit = sc
	se := fmt.Sprintf("span[title=%q]", tit)
	page.Timeout(to * 2).MustElement(se).MustClick()

	page.Timeout(to).MustElement(sel).MustClick()
	scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))

	tit = "circle"
	sel = "div.circle"
	WaitElementsLessThan(page.Timeout(to), sel, 1)
	scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))

	sel = "//*[@id='pvExplorationHost']/div/div/exploration/div/explore-canvas/div/div[2]/div/div[2]/div[2]/visual-container-repeat/visual-container[27]/transform/div/div[3]/div/visual-modern/div/div/div[2]/div[1]/div[3]/div/div[2]"
	ex(slide, getClientRect(page.MustElementX(sel), &vc27))
	sel = "//*[@id='pvExplorationHost']/div/div/exploration/div/explore-canvas/div/div[2]/div/div[2]/div[2]/visual-container-repeat/visual-container[22]/transform/div/div[3]/div/visual-modern/div/div/div[2]/div[1]/div[4]"
	ex(slide, getClientRect(page.MustElementX(sel), &vc22))
	sel = "div.visualContainerHost"
	ex(slide, getClientRect(page.MustElement(sel), &vcHost))

	bytes, err := page.Screenshot(false, &proto.PageCaptureScreenshot{
		Format: proto.PageCaptureScreenshotFormatJpeg,
		Clip: clip(
			vcHost.X,
			vcHost.Y,
			vc22.X+vc22.Width-vcHost.X,
			vc27.Y-vcHost.Y,
		),
	})
	ex(slide, err)

	ss(bytes).write(fmt.Sprintf("%02d.jpg", slide))
}
