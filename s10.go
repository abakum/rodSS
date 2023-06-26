package main

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	tg "github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/xlab/closer"
)

func s10(slide, deb int) {
	if deb == 0 {
		return
	}
	var (
		params = conf.P[strconv.Itoa(abs(slide))]
		timeline,
		text,
		hintbox,
		href string
		telegs = conf.P["98"]
		parts  []string
		eol    = tu.Entity("\n")
		space  = tu.Entity(" ")
		mecs   = []tu.MessageEntityCollection{}
	)
	ltf.Println(params)

	MessageID, err := strconv.Atoi(params[5])
	if err != nil {
		MessageID = 0
	}

	bot, err := tg.NewBot(telegs[0], tg.WithDefaultDebugLogger())
	ex(slide, err)
	i, err := strconv.ParseInt(telegs[1], 10, 64)
	ex(slide, err)
	ChatID := tu.ID(i)

	base, err := url.Parse(params[0])
	ex(slide, err)
	base.Path = ""
	base.RawQuery = ""

	br, ca := chrome(slide)
	page := chromePage(br, slide).
		MustNavigate(params[0]).MustWaitLoad()
	tit := page.MustInfo().Title

	sel := "button#login"
	se := "tbody > tr"
	page.Timeout(to).Race().Element(sel).MustHandle(func(e *rod.Element) {
		e.MustClick()
		sel = "input#name"
		page.Timeout(to).MustElement(sel).MustInput(params[1])
		sel = "input#password"
		page.Timeout(to).MustElement(sel).MustInput(params[2]).MustType(input.Enter)
		sdpt(slide, deb, page, tit)
		if !GetCookies(page, []string{}, slide) {
			GetAllCookies(br, slide)
		}
	}).Element(se).MustHandle(func(e *rod.Element) {
	}).MustDo()

	page.Timeout(to).MustElement(se).MustWaitStable()
	tit = page.MustInfo().Title
	sdpt(slide, deb, page, tit)

	suffix := []tu.MessageEntityCollection{
		tu.Entity("мониторить").TextLink(params[0]),
		tu.Entity(" ЕГЭ"),
	}
	ecs := []tu.MessageEntityCollection{
		tu.Entity("Начал "),
	}
	ecs = append(ecs, suffix...)
	for _, v := range params[6:] {
		ecs = append(ecs, tu.Entityf("\n%s", v))
	}
	MessageID, params[5] = delSend(bot, ChatID, MessageID, ecs...)
	closer.Bind(func() {
		if page != nil {
			page.MustClose()
		}
		if ca != nil {
			ca()
		}
		if bot != nil {
			ecs = []tu.MessageEntityCollection{
				tu.Entity("Прекратил "),
			}
			ecs = append(ecs, suffix...)
			MessageID, params[5] = delSend(bot, ChatID, MessageID, ecs...)
			bot.Close()
			conf.saver()
		}
	})
	for {
		ecs = []tu.MessageEntityCollection{}
		sel = "tr.nothing-to-show"
		if !page.MustHas(sel) {
			for _, tr := range page.Timeout(sec * 7).MustElements(se) {
				text, err = tr.Text()
				if err != nil {
					continue
				}
				if strings.Contains(text, "РЕШЕНО") {
					continue
				}
				for _, v := range params[6:] {
					if !strings.Contains(text, v) {
						continue
					}

					sel = "td > a.link-action[data-hintbox='1']"
					el, err := tr.Timeout(sec).Element(sel)
					if err != nil {
						continue
					}
					hintbox = sErr(el.Text())

					sel = "td.timeline-date > a"
					el, err = tr.Timeout(sec).Element(sel)
					timeline = ""
					href = ""
					if err == nil {
						timeline = sErr(el.Text())
						href = nErr(el.Attribute("href"))
					}

					sel = "td > a.link-action[aria-haspopup='true']"
					el, err = tr.Timeout(sec).Element(sel)
					parts = []string{"", "", "", ""}
					if err == nil {
						text, err = el.Text()
						if err == nil {
							parts = strings.Split(text, " | ")
						}
					}

					if len(parts) > 3 {
						ecs = append(ecs, tu.Entity(timeline).TextLink(base.String()+"/"+href), space)
						ecs = append(ecs, tu.Entity(parts[2]).Code(), space)
						ecs = append(ecs, tu.Entity(parts[1]).TextLink(params[3]+parts[1]), eol)
						ecs = append(ecs, tu.Entity(strings.TrimPrefix(parts[3], params[4])), eol)
						ecs = append(ecs, tu.Entity(hintbox), eol, eol)
					}
				}
			}
		}
		mec, _ := tu.MessageEntities(mecs...)
		ec, _ := tu.MessageEntities(ecs...)
		if mec != ec {
			mecs = ecs[:]
			if len(ecs) == 0 {
				ecs = []tu.MessageEntityCollection{
					tu.Entity("Продолжаю "),
				}
				ecs = append(ecs, suffix...)
			}
			MessageID, params[5] = delSend(bot, ChatID, MessageID, ecs...)
		}
		time.Sleep(sec * 30)
	}
}
