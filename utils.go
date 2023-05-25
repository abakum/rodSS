package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/xlab/closer"
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
	err := os.WriteFile(fullName, []byte(i), 0644)
	if err != nil {
		return
	}
	// exec.Command("rundll32", "url.dll,FileProtocolHandler", fullName).Run()
	exec.Command("cmd", "/c", "start", "chrome", fullName).Run()
	// stdo.Println(src(8), fullName)
}

func ex(slide int, err error) {
	if err != nil {
		exit = slide
		stdo.Println(src(8), err.Error())
		Scanln()
		closer.Close()
	}
}
func e(slide int, level int, err error) {
	if err != nil {
		exit = slide
		stdo.Println(src(level), err.Error())
		Scanln()
		closer.Close()
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

func getClientRect(el *rod.Element, clip *proto.PageViewport) error {
	if clip == nil {
		return fmt.Errorf("clip cannot be nil")
	}
	if el == nil {
		return fmt.Errorf("el cannot be nil")
	}
	// so that it won't clip the css-transformed element
	shape, err := el.Shape()
	if err != nil {
		return err
	}

	box := shape.Box()
	var Viewport proto.PageViewport
	Viewport.X = box.X
	Viewport.Y = box.Y
	Viewport.Width = box.Width
	Viewport.Height = box.Height
	Viewport.Scale = 1

	*clip = Viewport
	return nil
}

func clip(X, Y, Width, Height float64) *proto.PageViewport {
	clip := proto.PageViewport{
		X:      X,
		Y:      Y,
		Width:  Width,
		Height: Height,
		Scale:  1,
	}
	return &clip
}

func scs(slide, deb int, page *rod.Page, fn string) {
	stdo.Println(src(8), fn)
	if deb != slide {
		return
	}
	bytes, err := page.Screenshot(false, &proto.PageCaptureScreenshot{Format: proto.PageCaptureScreenshotFormatPng})
	if err == nil {
		ss(bytes).write(fn)
	}
}

func chrome() (b *rod.Browser, f func() error) {
	if multiBrowser {
		b = rod.New().ControlURL(launch().
			MustLaunch(),
		).MustConnect().
			Context(ctRoot)
		f = b.Close
	} else {
		b = bro
		f = func() error { return nil }
	}
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
func start(fu func(slide, deb int), slide, deb int, wg *sync.WaitGroup, started chan bool) {
	switch deb {
	case 0, slide, -slide:
	default:
		return
	}
	if wg != nil {
		wg.Add(1)
		if started != nil {
			started <- true
		}
		defer wg.Done()
	}
	fu(slide, deb)
	stdo.Printf("%s %02d Done\n", src(8), slide)

}

func abs(i int) int {
	if i > 0 {
		return i
	}
	return -i
}
func autoStart(started chan bool, d time.Duration) *time.Timer {
	for len(started) > 0 {
		<-started
	}
	return time.AfterFunc(d, func() {
		stdo.Println("auto started")
		started <- true
	})
}

func launch() (l *launcher.Launcher) {
	l = launcher.New().
		Leakless(false). //panic: open C:\Users\user\AppData\Local\Temp\leakless-0c3354cd58f0813bb5b34ddf3a7c16ed\leakless.exe: Access is denied.
		Bin(chromeBin).
		Delete("enable-automation").
		Set("start-maximized")
	if headLess {
		l = l.
			Set("headless", "new")
	} else {
		l = l.
			Headless(false).
			Logger(stdo.Writer())
	}
	time.Sleep(sec)
	return
}
func WaitElementsLessThan(p *rod.Page, selector string, num int) error {
	return p.Wait(rod.Eval(`(s, n) => document.querySelectorAll(s).length < n`, selector, num))
}
