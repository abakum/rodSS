package main

import (
	"fmt"
	"strconv"
	"time"
)

func s05(slide, deb int) {
	var (
		params = conf.P[strconv.Itoa(abs(slide))]
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

	tit = "Статистика по сотрудникам"
	sel := fmt.Sprintf("div[title=%q]", tit)
	page.Timeout(to).MustElement(sel).MustClick()
	sdpt(slide, deb, page, tit)

	cb(slide, deb, page, "СЦ/ЦЭ")

	tit = "Ср. длительность работ сотрудника за день, часы"
	sel = fmt.Sprintf("div[title=%q]", tit)
	t := page.MustElement(sel).MustShape().Box()

	sel = "div.innerContainer"
	ic := page.MustElement(sel).MustShape().Box()

	sel = "div.visualContainerHost"
	vch := page.MustElement(sel).MustShape().Box()
	vch.Width = t.X - vch.X
	vch.Height = ic.Y + ic.Height - 2*vch.Y
	sl(slide).done(page.Screenshot(true, clip(vch)))
}
