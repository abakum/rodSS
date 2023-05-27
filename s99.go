package main

import (
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
)

func s99(slide, deb int) {
	var (
		params = conf.P[strconv.Itoa(abs(slide))]
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

	tit = "ar-user-name"
	sel := fmt.Sprintf("input[name=%q]", tit)
	se := "div.multiBtnInner_xbp" //:nth-child(1)"
	page.Timeout(to).Race().Element(sel).MustHandle(func(e *rod.Element) {
		e.MustInput(params[1]).Page().Keyboard.MustType(input.Enter)

		tit = "ar-user-password"
		sel = fmt.Sprintf("input[name=%q]", tit)
		page.Timeout(to).MustElement(sel).MustInput(params[2]).Page().Keyboard.MustType(input.Enter)
		scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))
	}).Element(se).MustHandle(func(e *rod.Element) {
	}).MustDo()

	// page.Timeout(to).MustElement(se).MustClick()
	tit = "Редактировать"
	// stdo.Println(page.MustSearch(tit).MustHTML())
	page.Timeout(to).MustSearch(tit).MustClick()
	time.Sleep(ms)
	scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))

	tit = "Удалить"
	sel = "button.menu-button_J9B"
	if page.MustHas(sel) {
		page.Timeout(to).MustElement(sel).MustClick()
		sel = "button.align-left_-232488494:nth-child(3)"
		page.Timeout(to).MustElement(sel).MustClick()
	}
	tit = "Файл"
	sel = "button.addFilesBtn_RvX"
	page.Timeout(to).MustElement(sel).MustClick()
	time.Sleep(ms)
	sel = "button.align-left_-232488494:nth-child(3)"
	page.Timeout(to).MustElement(sel).MustClick()

	sel = "input[type=file]"
	page.Timeout(to).MustElement(sel).MustSetFiles(filepath.Join(root, mov))

	tit = "Загрузка контента"
	ti := "Загрузка отменена"
	sel = "div.title_-1510807907"
	time.Sleep(ms)
	page.Timeout(sec*3).WithPanic(func(x interface{}) {
		tit = "загрузка не началась за 3 секунды"
		scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))
		ex(slide, fmt.Errorf(tit))
	}).Race().ElementR(sel, ti).MustHandle(func(e *rod.Element) {
		scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, ti))
		ex(slide, fmt.Errorf(ti))
	}).ElementR(sel, tit).MustHandle(func(e *rod.Element) {
	}).MustDo()
	scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))

	tit = "Загрузка завершена"
	page.Timeout(sec*13).WithPanic(func(x interface{}) {
		tit = "загрузка не завершилась за 13 секунды"
		scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))
		ex(slide, fmt.Errorf(tit))
	}).Race().ElementR(sel, ti).MustHandle(func(e *rod.Element) {
		scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, ti))
		ex(slide, fmt.Errorf(ti))
	}).ElementR(sel, tit).MustHandle(func(e *rod.Element) {
	}).MustDo()
	scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))

	// for i := 0; i < 7; i++ {
	// 	stdo.Println(i, page.MustSearch("обр").MustHTML())
	// 	time.Sleep(sec)
	// }
	tit = "Обрабатывается"
	sel = fmt.Sprintf("div[title=%q]", tit)
	page.Timeout(sec * 3).Element(sel)
	scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))

	tit = "Обработка завершена"
	WaitElementsLessThan(page.Timeout(sec*7), sel, 1)
	scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))

	tit = "Сохранить и закрыть"
	// sel = "div.multiBtnInner_xbp:nth-child(4)"
	// page.Timeout(to).MustElement(sel).MustClick()
	// stdo.Println(page.MustSearch(tit).MustHTML())
	page.Timeout(to).MustSearch(tit).MustClick()
	time.Sleep(sec)
	scs(slide, slide, page, fmt.Sprintf("%02d.png", slide))
}
