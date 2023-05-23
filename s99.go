package main

import (
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"
)

func s99(slide int) {
	var (
		params = conf.P[strconv.Itoa(abs(slide))]
		err    error
		we     *rod.Element
		ok     bool
	)
	stdo.Println(params)
	browser, ca := chrome()
	defer ca()
	page := browser.MustPage().MustSetViewport(1920, 1080, 1, false)
	defer page.Close()
	ex(slide,
		page.Navigate(params[0]),
	)
	ex(slide,
		page.WaitLoad(),
	)
	tit := "ar-user-name"
	sel := fmt.Sprintf("input[name=%q]", tit)
	se := "div.multiBtnInner_xbp:nth-child(1)"
	for i := 0; i < 7; i++ {
		stdo.Println(i)
		ok, we, err = page.Has(sel)
		if ok {
			tit = "Авторизация"
			break
		}
		ok, we, err = page.Has(se)
		if ok {
			tit = "Кампания"
			break
		}
		time.Sleep(sec)
	}
	ex(slide, err)
	scs(slide, page, fmt.Sprintf("%02d %s.png", slide, tit))

	if tit == "Авторизация" {
		tit = "ar-user-name"
		ex(slide, we.Input(params[1]))
		ex(slide, page.Keyboard.Press(input.Enter))
		time.Sleep(ms)
		scs(slide, page, fmt.Sprintf("%02d %s.png", slide, tit))
		tit = "ar-user-password"
		sel = fmt.Sprintf("input[name=%q]", tit)
		we, err = page.Timeout(to).Element(sel)
		ex(slide, err)
		ex(slide, we.Input(params[2]))
		ex(slide, page.Keyboard.Press(input.Enter))
		time.Sleep(ms)
		scs(slide, page, fmt.Sprintf("%02d %s.png", slide, tit))

		we, err = page.Timeout(to).Element(se)
		ex(slide, err)
	}
	ex(slide, we.Click(proto.InputMouseButtonLeft, 1))
	time.Sleep(ms)
	scs(slide, page, fmt.Sprintf("%02d %s.png", slide, tit))

	tit = "Удалить"
	sel = "button.menu-button_J9B"
	ok, we, _ = page.Has(sel)
	stdo.Println(tit, ok)
	if ok {
		ex(slide, we.Click(proto.InputMouseButtonLeft, 1))
		time.Sleep(ms)
		sel = "button.align-left_-232488494:nth-child(3)"
		we, err = page.Timeout(to).Element(sel)
		ex(slide, err)
		ex(slide, we.Click(proto.InputMouseButtonLeft, 1))
		time.Sleep(ms)
	}

	tit = "Файл"
	sel = "button.addFilesBtn_RvX"
	we, err = page.Timeout(to).Element(sel)
	ex(slide, err)
	ex(slide, we.Click(proto.InputMouseButtonLeft, 1))
	time.Sleep(ms * 2)
	sel = "button.align-left_-232488494:nth-child(3)"
	we, err = page.Timeout(to).Element(sel)
	ex(slide, err)
	ex(slide, we.Click(proto.InputMouseButtonLeft, 1))
	time.Sleep(ms)
	upload = true
	scs(slide, page, fmt.Sprintf("%02d %s.png", slide, tit))

	sel = "input[type=file]"
	files := []string{filepath.Join(root, mov)}
	stdo.Println(files)
	we, err = page.Timeout(to).Element(sel)
	ex(slide, err)
	ex(slide, we.SetFiles(files))
	time.Sleep(ms)

	tit = "Загрузка"
	sel = fmt.Sprintf("//*[contains(text(),%q)]", tit)
	_, err = page.Timeout(sec * 3).ElementX(sel)
	stdo.Println(tit, err)
	if err != nil {
		scs(slide, page, fmt.Sprintf("%02d %s.png", slide, "Загрузка НЕ началась"))
		Scanln()
		return
	}
	scs(slide, page, fmt.Sprintf("%02d %s.png", slide, tit))

	tit = "отменена"
	sel = fmt.Sprintf("//*[contains(text(),%q)]", tit)
	_, err = page.Timeout(sec * 3).ElementX(sel)
	stdo.Println(tit, err)
	if err == nil {
		scs(slide, page, fmt.Sprintf("%02d %s.png", slide, tit))
		Scanln()
		return
	}

	tit = "завершена"
	sel = fmt.Sprintf("//*[contains(text(),%q)]", tit)
	_, err = page.Timeout(to).ElementX(sel)
	stdo.Println(tit, err)
	if err != nil {
		scs(slide, page, fmt.Sprintf("%02d %s.png", slide, "Загрузка НЕ завершена"))
		Scanln()
		return
	}
	scs(slide, page, fmt.Sprintf("%02d %s.png", slide, tit))
	time.Sleep(sec * 7)

	tit = "Сохранить и закрыть"
	sel = "div.multiBtnInner_xbp:nth-child(4)"
	we, err = page.Timeout(to).Element(sel)
	ex(slide, err)
	ex(slide, we.Click(proto.InputMouseButtonLeft, 1))
	time.Sleep(sec)
	scs(deb, page, fmt.Sprintf("%02d.png", slide))
	done(slide)
}
