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
	stdo.Println(params, sc, rf)
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

	tit = "mrf"
	sel := fmt.Sprintf("div[aria-label=%q]", tit)
	page.Timeout(to).MustElement(sel).MustClick()

	tit = params[1]
	sel = fmt.Sprintf("span[title=%q]", tit)
	page.Timeout(to * 2).MustElement(sel).MustClick()
	scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))

	tit = "РФ"
	sel = fmt.Sprintf("div[aria-label=%q]", tit)
	page.Timeout(to).MustElement(sel).MustClick()
	scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))

	tit = rf
	sel = fmt.Sprintf("span[title=%q]", tit)
	page.Timeout(to * 2).MustElement(sel).MustClick()
	scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))

	tit = "СЦ"
	sel = fmt.Sprintf("div[aria-label=%q]", tit)
	page.Timeout(to).MustElement(sel).MustClick()
	scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))

	tit = sc
	se := fmt.Sprintf("span[title=%q]", tit)
	page.Timeout(to * 2).MustElement(se).MustClick()

	page.Timeout(to).MustElement(sel).MustClick()
	scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))

	tit = "доля вне ЗО БТИ"
	sel = fmt.Sprintf("text[title=%q]", tit)
	page.Timeout(to).MustElement(sel)

	tit = "Кавр.Динамика по:"
	sel = fmt.Sprintf("text[title=%q]", tit)
	page.Timeout(to).MustElement(sel)

	tit = "Пред.Неделя"
	sel = "div"
	page.Timeout(to).MustElementR(sel, tit)
	scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))

	sel = "div.visualContainerHost"
	bytes, err := page.Timeout(to).MustElement(sel).Screenshot(proto.PageCaptureScreenshotFormatJpeg, 99)
	ex(slide, err)
	ss(bytes).write(fmt.Sprintf("%02d.jpg", slide))
}
