package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/xlab/closer"
)

const (
	doc       = "doc"
	bat       = "abaku.bat"
	mov       = "abaku.mp4"
	chDataDir = `Google\Chrome\User Data\Default`
	yaBrowser = `Yandex\YandexBrowser`
	yaBin     = `Application\browser.exe`
	yaDataDir = `User Data\Default`
	to        = time.Minute * 2
	ms        = time.Millisecond * 200
	sec       = time.Second
)

var (
	stdo         *log.Logger
	cd           string // s:\bin
	root         string // s:
	exit         int    = 2
	sc           string
	rf           string
	ctRoot       context.Context
	caRoot       context.CancelFunc
	multiBrowser = true
	headLess     = true
	sequentially = false
	userMode     = false
	yandex       = false
	bro          *rod.Browser
	bin          string
	userDataDir  = filepath.Join(os.Getenv("LOCALAPPDATA"), chDataDir)
	yowser       = filepath.Join(os.Getenv("LOCALAPPDATA"), yaBrowser)
)

func main() {
	var (
		wg  sync.WaitGroup
		err error
		ok  bool
	)
	defer func() {
		exit = 0
		closer.Close()
	}()
	stdo = log.New(os.Stdout, "", log.Lshortfile|log.Ltime)
	cd, err = os.Getwd()
	ex(2, err)
	stdo.Println("Getwd:", cd)
	root = filepath.Dir(cd)
	slides := []int{}

	for _, s := range os.Args[1:] {
		i, err := strconv.Atoi(s)
		if err != nil {
			continue
		}
		if i > 0 {
			headLess = false
		}
		if i < 0 {
			multiBrowser = true
			userMode = false
		}
		switch i {
		case 0:
			multiBrowser = false
		case 2, -2:
		case 3, -3:
			sequentially = true
		case 6, -6:
			userMode = true
			multiBrowser = false
		case 7, -7:
			yandex = true
		case 14:
			slides = []int{1, 4, 5, 8, 12, 13}
		case -14:
			slides = []int{-1, -4, -5, -8, -12, -13}
		case 100:
			slides = []int{97, 98, 99}
		case -100:
			slides = []int{-97, -98, -99}
		default:
			slides = append(slides, i)
		}
	}
	if len(slides) == 0 {
		slides = append(slides, 0)
	}
	ctRoot, caRoot = context.WithCancel(context.Background())
	bin, ok = launcher.LookPath()
	if !ok {
		yandex = true
		// ex(2, fmt.Errorf("not found chromeBin"))
	}
	if yandex {
		bin = filepath.Join(yowser, yaBin)
		userDataDir = filepath.Join(yowser, yaDataDir)
	}
	exeN, exeF, err := exeFN()
	ex(2, err)
	conf, err = loader(filepath.Join(cd, exeF+".json"))
	if err != nil {
		conf.P = map[string][]string{}
		conf.Ids = []int{}
		conf.saver()
		ex(2, err)
		return
	}
	sc = conf.P["4"][1]
	rf = conf.P["12"][2]

	if !multiBrowser {
		bro, _ = chrome(2)
	}
	closer.Bind(func() {
		caRoot()
		if !multiBrowser {
			bro.MustClose()
		}
		stdo.Println("main done", exit)
		switch {
		case exit == 0:
		case exit < 0:
			exit = -exit
			fallthrough
		default:
			time.Sleep(sec)
			taskKill("/F", "/IM", exeN, "/T")
		}
		os.Exit(exit)
	})
	started := make(chan int, 10)
	st := autoStart(started, sec)
	stdo.Println("multiBrowser:", multiBrowser, "headLess:", headLess, "sequentially:", sequentially, "userMode", userMode)
	for _, de := range slides {
		if abs(de) > 13 {
			break
		}
		stdo.Println(de)
		go start(s01, 1, de, &wg, started)
		go start(s04, 4, de, &wg, started)
		go start(s05, 5, de, &wg, started)
		go start(s08, 8, de, &wg, started)
		go start(s12, 12, de, &wg, started)
		go start(s13, 13, de, &wg, started)
		go start(s10, 10, de, &wg, started)
		if sequentially {
			wait(st, &wg, started)
			st = autoStart(started, sec)
		}
	}
	if !sequentially {
		wait(st, &wg, started)
	}
	for _, de := range slides {
		if de != 0 && abs(de) < 97 {
			continue
		}
		stdo.Println(de)
		start(s97, 97, de, nil, nil) //bat jpgs to mov
		st = autoStart(started, sec)
		go start(s98, 98, de, &wg, started) //telegram
		go start(s99, 99, de, &wg, started) //ss
		wait(st, &wg, started)
	}
}
