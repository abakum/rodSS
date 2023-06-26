package main

import (
	"strconv"
	"time"

	tg "github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/xlab/closer"
)

func s11(slide, deb int) {
	if deb == 0 {
		return
	}
	var (
		params = conf.P[strconv.Itoa(abs(slide))]
	)
	ltf.Println(params)

	bot, err := tg.NewBot(params[0], tg.WithDefaultDebugLogger())
	ex(slide, err)
	i, err := strconv.ParseInt(params[1], 10, 64)
	ex(slide, err)
	ChatID := tu.ID(i)
	MessageID, err := strconv.Atoi(params[2])
	if err != nil {
		MessageID = 0
	}

	ticker := time.NewTicker(time.Hour)
	ecs := []tu.MessageEntityCollection{
		tu.Entity("/Nǐ chīfàn le ma?").BotCommand(),
	}
	closer.Bind(func() {
		if bot != nil {
			ecs = []tu.MessageEntityCollection{
				tu.Entity("Chī le"),
			}
			ltf.Println(ecs, MessageID, ChatID)
			MessageID, params[2] = delSend(bot, ChatID, MessageID, ecs...)
			bot.Close()
			conf.saver()
		}
		ticker.Stop()
	})
	for {
		MessageID, params[2] = delSend(bot, ChatID, MessageID, ecs...)
		t := <-ticker.C
		ltf.Println("Tick at", t)
	}
}
