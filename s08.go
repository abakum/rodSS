package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func s08(slide, deb int) {
	var (
		TaskClosed = "TaskClosed.xlsx"
		params     = conf.P[strconv.Itoa(abs(slide))]
		ok         bool
		err        error
	)
	stdo.Println(params)
	exp := func(x interface{}) {
		e(slide, 14, x.(error))
	}
	br, ca := chrome()
	defer ca()
	page := br.WithPanic(exp).
		MustPage().WithPanic(exp).MustSetViewport(1920, 1080, 1, false).
		MustNavigate(params[0]).MustWaitLoad()
	defer page.Close()
	tit := page.MustInfo().Title
	scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))

	tit = "По работникам и типу задачи"
	sel := "#login_form-username"
	se := "span"
	for i := 0; i < 7; i++ {
		stdo.Println(i)
		ok, _, err = page.Has(sel)
		if ok {
			tit = "Авторизация"
			break
		}
		ok, _, err = page.HasR(se, tit)
		if ok {
			tit = "Отчеты по задачам"
			break
		}
		time.Sleep(ms)
	}
	ex(slide, err)
	if !ok {
		ex(slide, fmt.Errorf("argus error"))
	}
	scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))

	if tit == "Авторизация" {
		page.Timeout(to).MustElement(sel).MustInput(params[1])
		sel = "#login_form-password"
		page.Timeout(to).MustElement(sel).MustInput(params[2])

		tit = "Войти"
		sel = "span"
		page.Timeout(to).MustElementR(sel, tit).MustClick()
		scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))
	}

	tit = "По работникам и типу задачи"
	sel = "span"
	page.Timeout(to).MustElementR(sel, tit).MustClick()
	scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))

	tit = "месяцы"
	sel = "span"
	page.Timeout(to).MustElementR(sel, tit).MustClick()

	tit = "Обработка наряда"
	sel = "ul.ui-widget"
	page.Timeout(to).MustElement(sel).MustClick()
	sel = "//li[5]/label"
	page.Timeout(to).MustElementX(sel).MustClick()
	scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))

	sel = ".ui-tree-toggler"
	for i := 4; i < 9; i++ {
		page.Timeout(to).MustElements(sel)[i].MustClick()
		time.Sleep(ms)
	}

	tit = "Группа инсталляций"
	sel = "span"
	page.Timeout(to).MustElementR(sel, tit).MustClick()

	tit = "Группа клиентского сервиса"
	page.Timeout(to).MustElementR(sel, tit).MustClick()

	tit = "ОК"
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
