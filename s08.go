package main

import (
	"fmt"
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
	stdo.Println(params)
	exp := func(x interface{}) {
		e(slide, 14, x.(error))
	}
	br, ca := chrome(slide)
	defer ca()
	page := br.MustPage().WithPanic(exp).MustSetViewport(1920, 1080, 1, false).
		MustNavigate(params[0]).MustWaitLoad()
	defer page.MustClose()
	tit := page.MustInfo().Title
	scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))

	tit = "По работникам и типу задачи"
	sel := "input#login_form-username"
	page.Timeout(to).Race().Element(sel).MustHandle(func(e *rod.Element) {
		tit := "Авторизация"
		e.MustInput(params[1])
		sel = "input#login_form-password"
		e.Page().MustElement(sel).MustInput(params[2]).MustType(input.Enter)
		scs(slide, deb, e.Page(), fmt.Sprintf("%02d %s.png", slide, tit))
	}).Search(tit).MustHandle(func(e *rod.Element) {
	}).MustDo()

	page.Timeout(to).MustSearch(tit).MustClick()
	scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))

	tit = "месяцы"
	page.Timeout(to).MustSearch(tit).MustClick()

	tit = "Обработка наряда"
	sel = "ul.ui-widget"
	page.Timeout(to).MustElement(sel).MustClick()
	// sel = "//li[5]/label"
	// stdo.Println(page.Timeout(to).MustElementX(sel).MustHTML())
	sel = "li > label"
	page.Timeout(to).MustElementR(sel, tit).MustClick()
	scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))

	sel = "span.ui-tree-toggler"
	for i := 4; i < 9; i++ {
		page.Timeout(to).MustElements(sel)[i].MustClick()
		time.Sleep(ms)
	}

	tit = "Группа инсталляций"
	page.Timeout(to).MustSearch(tit).MustClick()

	tit = "Группа клиентского сервиса"
	page.Timeout(to).MustSearch(tit).MustClick()

	tit = "ОК"
	sel = "span"
	page.Timeout(to).MustElementR(sel, tit).MustClick()
	scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))

	sel = "span.ui-chkbox-label"
	page.Timeout(to).MustElement(sel).MustClick()
	sel = "button#report_actions_form-export_report_data > span"
	page.Timeout(to).MustElement(sel).MustClick()

	excel := filepath.Join(root, doc, TaskClosed)
	os.Remove(excel)
	wait := br.MustWaitDownload()

	tit = "EXCEL"
	sel = "span"
	page.Timeout(to).MustElementR(sel, tit).MustClick()

	ex(slide, os.WriteFile(excel, wait(), 0o644))
	scs(slide, slide, page, fmt.Sprintf("%02d.png", slide))

	s09(9)
}
