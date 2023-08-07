package main

import (
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
)

func s08(slide, deb int) {
	var (
		TaskClosed = "TaskClosed.xlsx"
		params     = conf.P[strconv.Itoa(abs(slide))]
	)
	ltf.Println(params)
	br, ca := chrome(slide)
	defer ca()
	page := chromePage(br, slide).
		MustNavigate(params[0]).MustWaitLoad()
	defer page.MustClose()
	tit := page.MustInfo().Title
	sdpt(slide, deb, page, tit)

	tit = "По работникам и типу задачи"
	sel := "input#login_form-username"
	page.Timeout(to).Race().Element(sel).MustHandle(func(e *rod.Element) {
		tit := "Авторизация"
		e.MustInput(params[1])
		sel = "input#login_form-password"
		e.Page().MustElement(sel).MustInput(params[2]).MustType(input.Enter)
		sdpt(slide, deb, e.Page(), tit)
		if !GetCookies(page, []string{}, slide) {
			GetAllCookies(br, slide)
		}
	}).Search(tit).MustHandle(func(e *rod.Element) {
	}).MustDo()

	page.Timeout(to).MustSearch(tit).MustClick()
	sdpt(slide, deb, page, tit)

	tit = "месяцы"
	page.Timeout(to).MustSearch(tit).MustClick()

	tit = "Обработка наряда"
	sel = "ul.ui-widget"
	page.Timeout(to).MustElement(sel).MustClick()
	// sel = "//li[5]/label"
	// stdo.Println(page.Timeout(to).MustElementX(sel).MustHTML())
	sel = "li > label"
	page.Timeout(to).MustElementR(sel, tit).MustClick()
	sdpt(slide, deb, page, tit)

	sel = "span.ui-tree-toggler"
	for i := 4; i < 9; i++ {
		page.Timeout(to).MustElements(sel)[i].MustClick()
		time.Sleep(ms * 2)
	}

	tit = "Группа инсталляций"
	page.Timeout(to).MustSearch(tit).MustClick()

	tit = "Группа клиентского сервиса"
	page.Timeout(to).MustSearch(tit).MustClick()

	tit = "ОК"
	sel = "span"
	page.Timeout(to).MustElementR(sel, tit).MustClick()
	sdpt(slide, deb, page, tit)

	sel = "span.ui-chkbox-label"
	page.Timeout(to * 2).MustElement(sel).MustClick()
	sel = "button#report_actions_form-export_report_data > span"
	page.Timeout(to).MustElement(sel).MustClick()

	excel := filepath.Join(root, doc, TaskClosed)
	os.Remove(excel)
	wait := br.MustWaitDownload()

	tit = "EXCEL"
	sel = "span"
	page.Timeout(to).MustElementR(sel, tit).MustClick()

	ex(slide, os.WriteFile(excel, wait(), 0o644))
	sp(slide, page)

	s09(9)
}
