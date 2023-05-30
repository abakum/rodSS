package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
)

func s04(slide, deb int) {
	var (
		params = conf.P[strconv.Itoa(abs(slide))]
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
	ib := page.MustElement(sel).MustShape().Box()

	sel = "div.visualContainerHost"
	vch := page.MustElement(sel).MustShape().Box()
	vch.Width = ib.X - vch.X
	sl(slide).done(page.Screenshot(true, clip(vch)))
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
	// if !page.MustElements(sel).Empty() {
	// 	page.Timeout(to).MustElement(sel).WaitInvisible()
	// }
	// WaitElementsLessThan(page.Timeout(to), sel, 1)
	// stdo.Println(len(page.MustElements(sel)))
	// page.Timeout(to).MustWaitStable()
	stdo.Println(len(page.MustElements(sel)))

	sdpt(slide, deb, page, tit)
}
