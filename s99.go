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
		err    error
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
	se := "div.multiBtnInner_xbp:nth-child(1)"
	ok := false
	page.Timeout(to).Race().Element(sel).MustHandle(func(e *rod.Element) {
		e.MustInput(params[1]).Page().Keyboard.MustType(input.Enter)

		tit = "ar-user-password"
		sel = fmt.Sprintf("input[name=%q]", tit)
		page.Timeout(to).MustElement(sel).MustInput(params[2]).Page().Keyboard.MustType(input.Enter)
		scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))

		ok = true
	}).Element(se).MustHandle(func(e *rod.Element) {
		ok = true
	}).MustDo()
	if !ok {
		ex(slide, fmt.Errorf("ss error"))
	}

	tit = "Кампания"
	page.Timeout(to).MustElement(se).MustClick()
	time.Sleep(ms)
	scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))

	tit = "Удалить"
	sel = "button.menu-button_J9B"
	if page.MustHas(sel) {
		page.Timeout(to).MustElement(sel).MustClick()
		sel = "button.align-left_-232488494:nth-child(3)"
		page.Timeout(to).MustElement(sel).MustClick()
	}
	// Scanln()
	tit = "Файл"
	sel = "button.addFilesBtn_RvX"
	page.Timeout(to).MustElement(sel).MustClick()
	time.Sleep(ms)
	sel = "button.align-left_-232488494:nth-child(3)"
	page.Timeout(to).MustElement(sel).MustClick()

	sel = "input[type=file]"
	page.Timeout(to).MustElement(sel).MustSetFiles(filepath.Join(root, mov))

	tit = "Загрузка контента"
	sel = "div.title_-1510807907"
	_, err = page.Timeout(sec*3).ElementR(sel, tit)
	if err != nil {
		tit = "загрузка не началась"
		scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))
		ex(slide, fmt.Errorf(tit))
	}
	scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))

	tit = "отменена"
	_, err = page.Timeout(sec*3).ElementR(sel, tit)
	if err == nil {
		tit = "загрузка отменена"
		scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))
		ex(slide, fmt.Errorf(tit))
	}

	tit = "Загрузка завершена"
	_, err = page.Timeout(to).ElementR(sel, tit)
	if err != nil {
		tit = "загрузка не завершена"
		scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))
		ex(slide, fmt.Errorf(tit))
	}
	scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))
	time.Sleep(sec * 7) //обработка

	tit = "Сохранить и закрыть"
	sel = "div.multiBtnInner_xbp:nth-child(4)"
	page.Timeout(to).MustElement(sel).MustClick()
	time.Sleep(sec)
	scs(slide, slide, page, fmt.Sprintf("%02d.png", slide))
}
