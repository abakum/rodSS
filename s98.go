package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	tg "github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func s98(slide, deb int) {
	var (
		bot      *tg.Bot
		file     *os.File
		medias   []tg.InputMedia
		messages []tg.Message
		err      error
		inds     = []int{1, 4, 5, 8, 12, 13, 97}
		params   = conf.P[strconv.Itoa(slide)]
	)
	ltf.Println(params)
	i, err := strconv.ParseInt(params[1], 10, 64)
	ex(slide, err)
	chat := tu.ID(i)
	for _, v := range inds {
		if _, err = os.Stat(i2p(v)); errors.Is(err, os.ErrNotExist) {
			ex(slide, err)
			return
		}
	}
	bot, err = tg.NewBot(params[0], tg.WithLogger(tg.Logger(Logger{}))) //  tg.WithDefaultDebugLogger()
	ex(slide, err)
	defer bot.Close()
	medias = []tg.InputMedia{}
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
			// bot.DeleteMessage(&telego.DeleteMessageParams{ChatID: chat, MessageID: v.MessageID})
			bot.DeleteMessage(DeleteMessage(chat, v.MessageID))
		}
		ltf.Println()
		return
	}
	for _, v := range conf.Ids {
		if v == 0 {
			continue
		}
		// bot.DeleteMessage(&telego.DeleteMessageParams{ChatID: chat, MessageID: v})
		bot.DeleteMessage(DeleteMessage(chat, v))
	}
	conf.Ids = []int{}
	for _, v := range messages {
		conf.Ids = append(conf.Ids, v.MessageID)
	}
	ex(slide, conf.saver())
}

// Custom loger type
type Logger struct{}

// Hide bot token
func woToken(format string, args ...any) (s string) {
	s = src(10) + " " + fmt.Sprintf(format, args...)
	btStart := strings.Index(s, "/bot") + 4
	if btStart > 4-1 {
		btLen := strings.Index(s[btStart:], "/")
		if btLen > 0 {
			s = s[:btStart] + s[btStart+btLen:]
		}
	}
	return
}

// Custom logger method for debug
func (Logger) Debugf(format string, args ...any) {
	lt.Print(woToken(format, args...))
}

// Custom logger method for error
func (Logger) Errorf(format string, args ...any) {
	let.Print(woToken(format, args...))
}
