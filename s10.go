package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/mymmrac/telego"
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
		text string
		telegs = conf.P["98"]
	)
	stdo.Println(params)
	MessageID, _ := strconv.Atoi(params[5])
	bot, err := telego.NewBot(telegs[0], telego.WithDefaultDebugLogger())
	ex(slide, err)
	defer bot.Close()
	i, err := strconv.ParseInt(telegs[1], 10, 64)
	ex(slide, err)
	chat := tu.ID(i)

	br, ca := chrome(slide)
	defer ca()
	page := chromePage(br, slide).
		MustNavigate(params[0]).MustWaitLoad()
	defer page.MustClose()
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

	mecs := []tu.MessageEntityCollection{}
	ecs := []tu.MessageEntityCollection{
		tu.Entity("Начал "),
		tu.Entity("мониторить").TextLink(params[0]),
		tu.Entity(" ЕГЭ:")}
	for _, v := range params[6:] {
		ecs = append(ecs, tu.Entityf("\n%s", v))
	}
	MessageID, params[5] = delSend(bot, chat, MessageID, ecs...)
	closer.Bind(func() {
		ecs = []tu.MessageEntityCollection{
			tu.Entity("Завершил "),
			tu.Entity("мониторить").TextLink(params[0]),
			tu.Entity(" ЕГЭ")}
		MessageID, params[5] = delSend(bot, chat, MessageID, ecs...)
		conf.saver()
	})
	for {
		ecs = []tu.MessageEntityCollection{}
		sel = "tr.nothing-to-show"
		if !page.MustHas(sel) {
			for _, el := range page.Timeout(sec).MustElements(se) {
				sel = "td.timeline-date > a"
				timeline = el.MustElement(sel).MustText()
				sel = "td > a.link-action"
				text = el.MustElement(sel).MustText()
				stdo.Println(timeline, text, el)
				for _, v := range params[6:] {
					if strings.Contains(text, v) {
						parts := strings.Split(text, " | ")
						if len(parts) == 4 {
							ecs = append(ecs, tu.Entity(timeline).Code(), tu.Entity("│"))
							ecs = append(ecs, tu.Entity(parts[1]).TextLink(params[3]+parts[1]), tu.Entity("│"))
							ecs = append(ecs, tu.Entity(parts[2]).Code(), tu.Entity("│"))
							ecs = append(ecs, tu.Entity(strings.TrimPrefix(parts[3], params[4])).Code(), tu.Entity("\n"))
						}
					}

				}
			}
		}
		mec, _ := tu.MessageEntities(mecs...)
		ec, _ := tu.MessageEntities(ecs...)
		if mec != ec {
			mecs = ecs[:]
			MessageID, params[5] = delSend(bot, chat, MessageID, mecs...)
		}
		time.Sleep(sec * 30)
	}
}

func delSend(bot *telego.Bot, chat telego.ChatID, MessageID int, mecs ...tu.MessageEntityCollection) (int, string) {
	bot.DeleteMessage(DeleteMessage(chat, MessageID))
	if len(mecs) > 0 {
		stdo.Println(tu.MessageEntities(mecs...))
		tm, err := bot.SendMessage(tu.MessageWithEntities(chat, mecs...))
		ex(10, err)
		MessageID = tm.MessageID
	}
	return MessageID, strconv.Itoa(MessageID)
}
