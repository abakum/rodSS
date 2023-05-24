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
		we     *rod.Element
		ok     bool
	)
	stdo.Println(params)
	exp := func(x interface{}) {
		e(slide, 14, x.(error))
	}
	browser, ca := chrome()
	defer ca()
	page := browser.WithPanic(exp).
		MustPage().WithPanic(exp).MustSetViewport(1920, 1080, 1, false).
		MustNavigate(params[0]).MustWaitLoad()
	defer page.Close()
	tit := page.MustInfo().Title
	scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))

	tit = "ar-user-name"
	sel := fmt.Sprintf("input[name=%q]", tit)
	se := "div.multiBtnInner_xbp:nth-child(1)"
	for i := 0; i < 7; i++ {
		stdo.Println(i)
		ok, _, err = page.Has(sel)
		if ok {
			tit = "Авторизация"
			break
		}
		ok, _, err = page.Has(se)
		if ok {
			tit = "Кампания"
			break
		}
		time.Sleep(sec)
	}
	ex(slide, err)
	scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))

	if tit == "Авторизация" {
		tit = "ar-user-name"
		page.Timeout(to).MustElement(sel).MustInput(params[1])
		ex(slide, page.Keyboard.Press(input.Enter))
		scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))

		tit = "ar-user-password"
		sel = fmt.Sprintf("input[name=%q]", tit)
		page.Timeout(to).MustElement(sel).MustInput(params[2])
		ex(slide, page.Keyboard.Press(input.Enter))
		scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))
	}
	tit = "Кампания"
	page.Timeout(to).MustElement(se).MustClick()
	time.Sleep(sec)
	scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))

	tit = "Удалить"
	sel = "button.menu-button_J9B"
	ok, _, _ = page.Has(sel)
	stdo.Println(tit, ok)
	if ok {
		page.Timeout(to).MustElement(sel).MustClick()
		sel = "button.align-left_-232488494:nth-child(3)"
		page.Timeout(to).MustElement(sel).MustClick()
	}

	tit = "Файл"
	sel = "button.addFilesBtn_RvX"
	page.Timeout(to).MustElement(sel).MustClick()
	time.Sleep(ms * 2)
	sel = "button.align-left_-232488494:nth-child(3)"
	page.Timeout(to).MustElement(sel).MustClick()
	// upload = true

	sel = "input[type=file]"
	page.Timeout(to).MustElement(sel).MustSetFiles(filepath.Join(root, mov))

	tit = "Загрузка контента"
	sel = "div"
	_, err = page.Timeout(sec*3).ElementR(sel, tit)
	if err != nil {
		tit = "загрузка не началась"
		scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))
		ex(slide, fmt.Errorf(tit))
	}
	scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))
	// stdo.Println(we.MustHTML())

	tit = "отменена"
	we, err = page.Timeout(sec*3).ElementR(sel, tit)
	if err == nil {
		stdo.Println(we.MustHTML())
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
	// stdo.Println(we.MustHTML())
	scs(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))
	time.Sleep(sec * 7) //обработка

	tit = "Сохранить и закрыть"
	sel = "div.multiBtnInner_xbp:nth-child(4)"
	page.Timeout(to).MustElement(sel).MustClick()
	time.Sleep(sec)
	scs(slide, slide, page, fmt.Sprintf("%02d.png", slide))
}
