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
	br, ca := chrome(slide)
	defer ca()
	page := chromePage(br, slide)
	defer page.MustClose()
	page.Navigate(params[0])
	time.Sleep(sec)
	tit := page.MustWaitLoad().MustInfo().Title
	sdpt(slide, deb, page, tit)

	tit = params[1]
	sel := fmt.Sprintf("div[aria-label=%q]", tit)
	page.Timeout(to * 2).MustElement(sel).MustClick()
	sdpt(slide, deb, page, tit)

	tit = "RF"
	sel = fmt.Sprintf("div[aria-label=%q]", tit)
	page.Timeout(to).MustElement(sel).MustClick()
	sdpt(slide, deb, page, tit)

	tit = params[2]
	sel = fmt.Sprintf("span[title=%q]", tit)
	page.Timeout(to).MustElement(sel).MustClick()
	sdpt(slide, deb, page, tit)

	tit = "SC_NAME"
	sel = fmt.Sprintf("div[aria-label=%q]", tit)
	page.Timeout(to).MustElement(sel).MustClick()
	sdpt(slide, deb, page, tit)

	tit = sc
	{
		sel := fmt.Sprintf("span[title=%q]", tit)
		page.Timeout(to * 2).MustElement(sel).MustClick()
	}
	page.Timeout(to).MustElement(sel).MustClick()
	sdpt(slide, deb, page, tit)

	sel = "div.circle"
	WaitElementsLessThan(page.Timeout(to), sel, 1)

	sel = "//*[@id='pvExplorationHost']/div/div/exploration/div/explore-canvas/div/div[2]/div/div[2]/div[2]/visual-container-repeat/visual-container[27]/transform/div/div[3]/div/visual-modern/div/div/div[2]/div[1]/div[3]/div/div[2]"
	ex(slide, getClientRect(page.MustElementX(sel), &vc27))
	sel = "//*[@id='pvExplorationHost']/div/div/exploration/div/explore-canvas/div/div[2]/div/div[2]/div[2]/visual-container-repeat/visual-container[22]/transform/div/div[3]/div/visual-modern/div/div/div[2]/div[1]/div[4]"
	ex(slide, getClientRect(page.MustElementX(sel), &vc22))
	sel = "div.visualContainerHost"
	ex(slide, getClientRect(page.MustElement(sel), &vcHost))

	sl(slide).done(page.Screenshot(false, &proto.PageCaptureScreenshot{
		Format: proto.PageCaptureScreenshotFormatJpeg,
		Clip: clip(
			vcHost.X,
			vcHost.Y,
			vc22.X+vc22.Width-vcHost.X,
			vc27.Y-vcHost.Y,
		),
	}))
}
