package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/xlab/closer"
	"github.com/ysmood/gson"
)

func i2p(v int) (fn string) {
	fn = fmt.Sprintf("%02d.jpg", v)
	if v == 97 {
		fn = mov
	}
	fn = filepath.Join(root, fn)
	return
}

type ss []byte

func (i ss) write(fileName string) {
	fullName := filepath.Join(root, doc, fileName)
	jpg := strings.HasSuffix(fileName, ".jpg")
	if jpg {
		fullName = filepath.Join(root, fileName)
	}
	err := os.WriteFile(fullName, []byte(i), 0o644)
	if err != nil {
		return
	}
	if userMode {
		exec.Command("rundll32", "url.dll,FileProtocolHandler", fullName).Run()
		// exec.Command("cmd", "/c", "start", "browser", fullName).Run() //for yandex browser
	} else {
		exec.Command("cmd", "/c", "start", "chrome", fullName).Run()
	}
	// exec.Command(chromeBin, fullName).Run() //not closed

}

func ex(slide int, err error) {
	if err != nil {
		exit = slide
		stdo.Println(src(8), err.Error())
		Scanln()
		closer.Close()
		runtime.Goexit()
	}
}
func e(slide int, level int, err error) {
	if err != nil {
		exit = slide
		stdo.Println(src(level), err.Error())
		Scanln()
		closer.Close()
		runtime.Goexit()
	}
}

func src(deep int) (s string) {
	s = string(debug.Stack())
	// for k, v := range strings.Split(s, "\n") {
	// 	stdo.Println(k, v)
	// }
	s = strings.Split(s, "\n")[deep]
	s = strings.Split(s, " +0x")[0]
	_, s = path.Split(s)
	return
}

func sdpt(slide, deb int, page *rod.Page, tit string) {
	sdpf(slide, deb, page, fmt.Sprintf("%02d %s.png", slide, tit))
}

func sp(slide int, page *rod.Page) {
	sdpf(slide, slide, page, fmt.Sprintf("%02d.png", slide))
}

func sdpf(slide, deb int, page *rod.Page, fn string) {
	stdo.Println(src(10), fn)
	if deb != slide {
		return
	}
	bytes, err := page.Screenshot(false, &proto.PageCaptureScreenshot{Format: proto.PageCaptureScreenshotFormatPng})
	if err == nil {
		ss(bytes).write(strings.ReplaceAll(fn, ":", ""))
	}
}
func launch() (l *launcher.Launcher) {
	if userMode {
		_, exeN := filepath.Split(bin)
		taskKill("/f", "/im", exeN)
		time.Sleep(sec)
		l = launcher.NewUserMode()
	} else {
		l = launcher.New().
			Leakless(false). //panic: open C:\Users\user\AppData\Local\Temp\leakless-0c3354cd58f0813bb5b34ddf3a7c16ed\leakless.exe: Access is denied.
			Bin(bin)
	}
	l = l.
		Delete("enable-automation").
		Set("start-maximized")
	if !multiBrowser && !userMode {
		l = l.
			UserDataDir(userDataDir)
	}
	if headLess {
		l = l.
			Set("headless", "new")
	} else {
		l = l.
			Headless(false).
			Logger(stdo.Writer())
	}
	return
}

func chrome(slide int) (b *rod.Browser, f func()) {
	if multiBrowser || slide == 2 {
		b = rod.New().
			WithPanic(func(x interface{}) {
				e(slide, 14, x.(error))
			}).
			ControlURL(launch().
				MustLaunch(),
			).MustConnect().
			// MustIgnoreCertErrors(true).
			Context(ctRoot)
		f = b.MustClose
	} else {
		b = bro
		f = func() {}
	}
	if !headLess {
		b = b.SlowMotion(sec).Trace(true)
	}
	stdo.Println(b.MustVersion())
	return
}

func exeFN() (string, string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", "", err
	}
	_, exeN := filepath.Split(exe)
	return exeN, strings.TrimSuffix(exeN, filepath.Ext(exeN)), err
}

func taskKill(arg ...string) {
	cmd := exec.Command("taskKill.exe", arg...)
	stdo.Println(src(8), cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func Scanln() {
	if headLess {
		return
	}
	stdo.Print(src(8), "\nPress Enter>")
	fmt.Scanln()
}

func start(fu func(slide, deb int), slide, deb int, wg *sync.WaitGroup, started chan int) {
	switch deb {
	case 0, slide, -slide:
	default:
		return
	}
	if wg != nil {
		wg.Add(1)
		if started != nil {
			started <- slide
		}
		defer wg.Done()
	}
	fu(slide, deb)
	stdo.Printf("%s %02d done\n", src(8), slide)
}

func abs(i int) int {
	if i > 0 {
		return i
	}
	return -i
}

func autoStart(started chan int, d time.Duration) *time.Timer {
	for len(started) > 0 {
		<-started
	}
	return time.AfterFunc(d, func() {
		started <- 0
	})
}

func wait(st *time.Timer, wg *sync.WaitGroup, started chan int) {
	i := <-started
	st.Stop()
	if i == 0 {
		stdo.Println("auto started")
	} else {
		stdo.Printf("s%02d %s", i, "started")
	}
	wg.Wait()
	stdo.Println("all done")
}

type sl int

func (slide sl) done(bytes []byte, err error) {
	ex(int(slide), err)
	ss(bytes).write(fmt.Sprintf("%02d.jpg", slide))
}

func chromePage(br *rod.Browser, slide int) (page *rod.Page) {
	page = br.MustPage().WithPanic(func(x interface{}) {
		e(slide, 14, x.(error))
	}).MustSetViewport(1920, 1080, 1, false)
	SetCookiesP(page, slide)
	if headLess {
		return
	}
	return page.MustWindowFullscreen()
}

func clip(r *proto.DOMRect) (clip *proto.PageCaptureScreenshot) {
	clip = &proto.PageCaptureScreenshot{
		Format:  proto.PageCaptureScreenshotFormatJpeg,
		Quality: gson.Int(99),
		Clip: &proto.PageViewport{
			X:      r.X,
			Y:      r.Y,
			Width:  r.Width,
			Height: r.Height,
			Scale:  1,
		},
		FromSurface: true,
	}
	return
}

func WaitElementsLessThan(p *rod.Page, selector string, num int) error {
	return p.Wait(rod.Eval(`(s, n) => document.querySelectorAll(s).length < n`, selector, num))
}
func GetCookiesB(br *rod.Browser, slide int) {
	cookies, err := br.GetCookies()
	if err == nil {
		bytes, err := json.Marshal(&cookies)
		if err == nil {
			os.WriteFile(filepath.Join(cd, fmt.Sprintf("%02d.json", slide)), bytes, 0o644)
		}
	}
}
func SetCookiesB(br *rod.Browser, slide int) {
	bytes, err := os.ReadFile(filepath.Join(cd, fmt.Sprintf("%02d.json", slide)))
	if err == nil {
		cookies := []*proto.NetworkCookie{}
		if json.Unmarshal(bytes, &cookies) == nil {
			br.SetCookies(proto.CookiesToParams(cookies))
		}
	}
}
func GetCookiesP(page *rod.Page, slide int) {
	cookies, err := page.Cookies([]string{})
	if err != nil {
		bytes, err := json.Marshal(&cookies)
		if err == nil {
			os.WriteFile(filepath.Join(cd, fmt.Sprintf("%02d.json", slide)), bytes, 0o644)
		}
	}
}
func SetCookiesP(page *rod.Page, slide int) {
	bytes, err := os.ReadFile(filepath.Join(cd, fmt.Sprintf("%02d.json", slide)))
	if err == nil {
		cookies := []*proto.NetworkCookie{}
		if json.Unmarshal(bytes, &cookies) == nil {
			page.SetCookies(proto.CookiesToParams(cookies))
		}
	}
}
