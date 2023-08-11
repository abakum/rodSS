package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-rod/rod/lib/proto"
)

func s13(slide, deb int) {
	var (
		params = conf.P[strconv.Itoa(abs(slide))]
	)
	ltf.Println(params, sc, rf)
	br, ca := chrome(slide)
	defer ca()
	page := chromePage(br, slide)
	defer page.MustClose()
	page.Navigate(params[0])
	time.Sleep(sec)
	tit := page.MustWaitLoad().MustInfo().Title
	sdpt(slide, deb, page, tit)

	tit = "mrf"
	sel := fmt.Sprintf("div[aria-label=%q]", tit)
	page.Timeout(to * 2).MustElement(sel).MustClick()

	tit = params[1]
	sel = fmt.Sprintf("span[title=%q]", tit)
	page.Timeout(to * 2).MustElement(sel).MustClick()
	sdpt(slide, deb, page, tit)

	tit = "РФ"
	sel = fmt.Sprintf("div[aria-label=%q]", tit)
	page.Timeout(to).MustElement(sel).MustClick()
	sdpt(slide, deb, page, tit)

	tit = rf
	sel = fmt.Sprintf("span[title=%q]", tit)
	page.Timeout(to * 2).MustElement(sel).MustClick()
	sdpt(slide, deb, page, tit)

	tit = "СЦ"
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

	tit = "доля вне ЗО БТИ"
	sel = fmt.Sprintf("text[title=%q]", tit)
	page.Timeout(to).MustElement(sel)

	tit = "Кавр.Динамика по:"
	sel = fmt.Sprintf("text[title=%q]", tit)
	page.Timeout(to).MustElement(sel)

	tit = "Пред.Неделя"
	sel = "div"
	page.Timeout(to*2).MustElementR(sel, tit)
	sdpt(slide, deb, page, tit)

	sel = "div.circle"
	WaitElementsLessThan(page.Timeout(to), sel, 1)
	// if !page.MustElements(sel).Empty() {
	// 	page.Timeout(to).MustElement(sel).WaitInvisible()
	// }
	// page.Timeout(to).MustWaitStable()

	sel = "div.visualContainerHost"
	sl(slide).done(page.Timeout(to).MustElement(sel).Screenshot(proto.PageCaptureScreenshotFormatJpeg, 99))
}
