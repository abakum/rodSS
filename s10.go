package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
)

func s10(slide, deb int) {
	if deb == 0 {
		return
	}
	var (
		params = conf.P[strconv.Itoa(abs(slide))]
	)
	stdo.Println(params)
	br, ca := chrome(slide)
	defer ca()
	// bytes, err := os.ReadFile(cf)
	// if err == nil {
	// 	cookies := []*proto.NetworkCookie{}
	// 	if json.Unmarshal(bytes, &cookies) == nil {
	// 		br.SetCookies(proto.CookiesToParams(cookies))
	// 	}
	// }
	page := chromePage(br, slide).
		MustNavigate(params[0]).MustWaitLoad()
	defer page.MustClose()
	tit := page.MustInfo().Title
	sdpt(slide, deb, page, tit)

	sel := "button#login"
	se := "table.list-table"
	page.Timeout(to).Race().Element(sel).MustHandle(func(e *rod.Element) {
		tit := "Авторизация"
		e.MustClick()
		sel = "input#name"
		page.Timeout(to).MustElement(sel).MustInput(params[1])
		sel = "input#password"
		page.Timeout(to).MustElement(sel).MustInput(params[2]).MustType(input.Enter)
		sdpt(slide, deb, page, tit)
		GetCookiesP(page, slide)
	}).Element(se).MustHandle(func(e *rod.Element) {
	}).MustDo()

	tit = page.MustInfo().Title
	sdpt(slide, deb, page, tit)
	mess := ""
	for {
		messN := ""
		sel = "tr.nothing-to-show"
		if !page.MustHas(sel) {
			sel = "tbody > tr"
			for _, el := range page.Timeout(sec).MustElements(sel) {
				sel = "td.timeline-date > a"
				timeline := el.MustElement(sel).MustText()
				sel = "td > a.link-action"
				text := el.MustElement(sel).MustText()
				for _, v := range []string{
					"станица Боковская",
					"станица Казанская",
					"пос Тарасовский",
					"пос Чертково",
					"станица Вешенская",
					"г. Миллерово",
				} {
					if strings.Contains(text, v) {
						messN = fmt.Sprintf("%s %s\n%s", timeline, text, messN)
						stdo.Println(timeline)
						stdo.Println(text)
					}

				}
			}
		}
		if mess != messN {
			mess = messN
			if mess == "" {

			} else {
				stdo.Println(mess)
			}
		}
		time.Sleep(sec * 60)
		stdo.Println(time.Now())
	}
}
