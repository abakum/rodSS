package main

import (
	"fmt"
	"strconv"
	"time"
)

func s12(slide, deb int) {
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

	tit = params[1]
	sel := fmt.Sprintf("div[aria-label=%q]", tit)
	page.Timeout(to * 3).MustElement(sel).MustClick()
	sdpt(slide, deb, page, tit)

	tit = "RF"
	sel = fmt.Sprintf("div[aria-label=%q]", tit)
	page.Timeout(to).MustElement(sel).MustClick()
	sdpt(slide, deb, page, tit)

	tit = params[2]
	sel = fmt.Sprintf("span[title=%q]", tit)
	page.Timeout(to * 2).MustElement(sel).MustClick()
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
	// if !page.MustElements(sel).Empty() {
	// 	page.Timeout(to).MustElement(sel).WaitInvisible()
	// }
	// page.Timeout(to).MustWaitStable()

	sel = "//*[@id='pvExplorationHost']/div/div/exploration/div/explore-canvas/div/div[2]/div/div[2]/div[2]/visual-container-repeat/visual-container[27]/transform/div/div[3]/div/visual-modern/div/div/div[2]/div[1]/div[3]/div/div[2]"
	vc27 := page.MustElementX(sel).MustShape().Box()

	sel = "//*[@id='pvExplorationHost']/div/div/exploration/div/explore-canvas/div/div[2]/div/div[2]/div[2]/visual-container-repeat/visual-container[22]/transform/div/div[3]/div/visual-modern/div/div/div[2]/div[1]/div[4]"
	vc22 := page.MustElementX(sel).MustShape().Box()

	sel = "div.visualContainerHost"
	vch := page.MustElement(sel).MustShape().Box()
	vch.Width = vc22.X + vc22.Width - vch.X
	vch.Height = vc27.Y - vch.Y
	sl(slide).done(page.Screenshot(true, clip(vch)))
}
