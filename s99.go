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
	ltf.Println(params)
	br, ca := chrome(slide)
	defer ca()
	page := chromePage(br, slide).
		MustNavigate(params[0]).MustWaitLoad()
	defer page.MustClose()
	tit := page.MustInfo().Title
	sdpt(slide, deb, page, tit)

	tit = "ar-user-name"
	sel := fmt.Sprintf("input[name=%q]", tit)
	tit = "Редактировать"
	page.Timeout(to).Race().Element(sel).MustHandle(func(e *rod.Element) {
		e.MustInput(params[1]).Page().Keyboard.MustType(input.Enter)
		{
			tit := "ar-user-password"
			sel = fmt.Sprintf("input[name=%q]", tit)
			page.Timeout(to).MustElement(sel).MustInput(params[2]).Page().Keyboard.MustType(input.Enter)
			sdpt(slide, deb, page, tit)
			if !GetCookies(page, []string{}, slide) {
				if !GetCookies(page, []string{".wapi.ds.rt.ru"}, slide) {
					GetAllCookies(br, slide)
				}
			}
		}
	}).Search(tit).MustHandle(func(e *rod.Element) {
	}).MustDo()
	// Scanln()
	page.Timeout(to).MustSearch(tit).MustClick()
	time.Sleep(ms)
	sdpt(slide, deb, page, tit)

	tit = "Удалить"
	sel = "button.menu-button_J9B"
	if page.MustHas(sel) {
		page.Timeout(to).MustElement(sel).MustClick()
		sel = "button.button_-794993099"
		page.Timeout(to).MustElementR(sel, tit).MustClick()
	}
	tit = "Файл"
	sel = "button.addFilesBtn_RvX"
	page.Timeout(to).MustElement(sel).MustClick()
	time.Sleep(ms)
	sel = "button.button_-794993099"
	page.Timeout(to).MustElementR(sel, tit).MustClick()

	sel = "input[type=file]"
	page.Timeout(to).MustElement(sel).MustSetFiles(filepath.Join(root, mov))

	tit = "Загрузка контента"
	ti := "Загрузка отменена"
	sel = "div.title_-1510807907"
	time.Sleep(ms)
	page.Timeout(sec * 3).WithPanic(func(x interface{}) {
		tit = "загрузка не началась за 3 секунды"
		sdpt(slide, deb, page, tit)
		ex(slide, fmt.Errorf(tit))
	}).Race().Search(ti).MustHandle(func(e *rod.Element) {
		sdpf(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, ti))
		ex(slide, fmt.Errorf(ti))
	}).Search(tit).MustHandle(func(e *rod.Element) {
	}).MustDo()
	sdpt(slide, deb, page, tit)

	tit = "Загрузка завершена"
	page.Timeout(sec * 33).WithPanic(func(x interface{}) {
		tit = "загрузка не завершилась за 33 секунд"
		sdpt(slide, deb, page, tit)
		ex(slide, fmt.Errorf(tit))
	}).Race().Search(ti).MustHandle(func(e *rod.Element) {
		ltf.Println(page.MustSearch(ti).MustHTML())
		sdpf(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, ti))
		ex(slide, fmt.Errorf(ti))
	}).Search(tit).MustHandle(func(e *rod.Element) {
	}).MustDo()
	sdpt(slide, deb, page, tit)

	tit = "Обрабатывается"
	// for i := 0; i < 7; i++ {
	// 	stdo.Println(i, page.MustSearch(tit).MustHTML())
	// 	time.Sleep(sec)
	// }
	sel = fmt.Sprintf("div[title=%q]", tit)
	WaitElementsLessThan(page.Timeout(sec*7), sel, 1)
	// if !page.MustElements(sel).Empty() {
	// 	page.Timeout(sec * 7).MustElement(sel).WaitInvisible()
	// }
	// page.Timeout(to).MustWaitStable()

	tit = "Обработка завершена"
	sdpt(slide, deb, page, tit)

	tit = "Сохранить и закрыть"
	page.Timeout(to).MustSearch(tit).MustClick()
	time.Sleep(sec)
	sp(slide, page)
}
