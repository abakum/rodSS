package main

import (
	"errors"
	"os"
	"strconv"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func s98(slide int) {
	var (
		bot      *telego.Bot
		file     *os.File
		medias   []telego.InputMedia
		messages []telego.Message
		err      error
		inds     = []int{1, 4, 5, 8, 12, 13, 97}
		params   = conf.P[strconv.Itoa(slide)]
	)
	stdo.Println(params)
	i, err := strconv.ParseInt(params[1], 10, 64)
	ex(slide, err)
	chat := tu.ID(i)
	for _, v := range inds {
		if _, err = os.Stat(i2p(v)); errors.Is(err, os.ErrNotExist) {
			ex(slide, err)
			return
		}
	}
	bot, err = telego.NewBot(params[0], telego.WithDefaultDebugLogger())
	ex(slide, err)
	defer bot.Close()
	medias = []telego.InputMedia{}
	for _, v := range inds {
		file, err = os.Open(i2p(v))
		if err != nil {
			ex(slide, err)
			return
		}
		defer file.Close()
		switch v {
		case 1:
			medias = append(medias, tu.MediaPhoto(tu.File(file)).WithCaption("⚡#умныеЭкраны"))
		case 97:
			medias = append(medias, tu.MediaVideo(tu.File(file)))
		default:
			medias = append(medias, tu.MediaPhoto(tu.File(file)))
		}
	}
	messages, _ = bot.SendMediaGroup(tu.MediaGroup(chat).WithMedia(medias...))
	if len(messages) != len(medias) {
		for _, v := range messages {
			bot.DeleteMessage(&telego.DeleteMessageParams{ChatID: chat, MessageID: v.MessageID})
		}
		stdo.Println()
		return
	}
	for _, v := range conf.Ids {
		if v == 0 {
			continue
		}
		bot.DeleteMessage(&telego.DeleteMessageParams{ChatID: chat, MessageID: v})
	}
	conf.Ids = []int{}
	for _, v := range messages {
		conf.Ids = append(conf.Ids, v.MessageID)
	}
	err = conf.saver()
	ex(slide, err)
	done(slide)
}
