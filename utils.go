package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/runtime"
	dp "github.com/chromedp/chromedp"
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

func Screenshot(sel interface{}, picbuf *[]byte, quality int, opts ...dp.QueryOption) dp.QueryAction {
	if picbuf == nil {
		panic("picbuf cannot be nil")
	}

	return dp.QueryAfter(sel, func(ctx context.Context, execCtx runtime.ExecutionContextID, nodes ...*cdp.Node) error {
		if len(nodes) < 1 {
			return fmt.Errorf("selector %q did not return any nodes", sel)
		}

		// get box model
		var clip page.Viewport
		if err := callFunctionOnNode(ctx, nodes[0], getClientRectJS, &clip); err != nil {
			return err
		}
		// The "Capture node screenshot" command does not handle fractional dimensions properly.
		// Let's align with puppeteer:
		// https://github.com/puppeteer/puppeteer/blob/bba3f41286908ced8f03faf98242d4c3359a5efc/src/common/Page.ts#L2002-L2011
		x, y := math.Round(clip.X), math.Round(clip.Y)
		clip.Width, clip.Height = math.Round(clip.Width+clip.X-x), math.Round(clip.Height+clip.Y-y)
		clip.X, clip.Y = x, y

		// The next comment is copied from the original code.
		// This seems to be necessary? Seems to do the right thing regardless of DPI.
		clip.Scale = 1

		format := page.CaptureScreenshotFormatPng
		if quality != 100 {
			format = page.CaptureScreenshotFormatJpeg
		}
		// take screenshot of the box
		buf, err := page.CaptureScreenshot().
			WithFormat(format).
			WithQuality(int64(quality)).
			WithCaptureBeyondViewport(true).
			WithFromSurface(true).
			WithClip(&clip).
			Do(ctx)
		if err != nil {
			return err
		}

		*picbuf = buf
		return nil
	}, append(opts, dp.NodeVisible)...)
}
func FullScreenshot(res *[]byte, quality int, clip *page.Viewport) dp.EmulateAction {
	if res == nil {
		panic("res cannot be nil")
	}
	if clip == nil {
		panic("clip cannot be nil")
	}
	return dp.ActionFunc(func(ctx context.Context) error {
		format := page.CaptureScreenshotFormatPng
		if quality != 100 {
			format = page.CaptureScreenshotFormatJpeg
		}

		// capture screenshot
		var err error
		*res, err = page.CaptureScreenshot().
			WithCaptureBeyondViewport(true).
			WithFromSurface(true).
			WithFormat(format).
			WithQuality(int64(quality)).
			WithClip(clip).
			Do(ctx)
		if err != nil {
			return err
		}
		return nil
	})
}

func getClientRect(sel interface{}, clip *page.Viewport, opts ...dp.QueryOption) dp.QueryAction {
	if clip == nil {
		panic("clip cannot be nil")
	}

	return dp.QueryAfter(sel, func(ctx context.Context, execCtx runtime.ExecutionContextID, nodes ...*cdp.Node) error {
		if len(nodes) < 1 {
			return fmt.Errorf("selector %q did not return any nodes", sel)
		}

		// get box model
		var Viewport page.Viewport
		if err := callFunctionOnNode(ctx, nodes[0], getClientRectJS, &Viewport); err != nil {
			return err
		}
		// The "Capture node screenshot" command does not handle fractional dimensions properly.
		// Let's align with puppeteer:
		// https://github.com/puppeteer/puppeteer/blob/bba3f41286908ced8f03faf98242d4c3359a5efc/src/common/Page.ts#L2002-L2011
		x, y := math.Round(Viewport.X), math.Round(Viewport.Y)
		Viewport.Width, Viewport.Height = math.Round(Viewport.Width+Viewport.X-x), math.Round(Viewport.Height+Viewport.Y-y)
		Viewport.X, Viewport.Y = x, y

		// The next comment is copied from the original code.
		// This seems to be necessary? Seems to do the right thing regardless of DPI.
		Viewport.Scale = 1

		*clip = Viewport
		return nil
	}, append(opts, dp.NodeVisible)...)
}

func callFunctionOnNode(ctx context.Context, node *cdp.Node, function string, res interface{}, args ...interface{}) error {
	r, err := dom.ResolveNode().WithNodeID(node.NodeID).Do(ctx)
	if err != nil {
		return err
	}
	err = dp.CallFunctionOn(function, res,
		func(p *runtime.CallFunctionOnParams) *runtime.CallFunctionOnParams {
			return p.WithObjectID(r.ObjectID)
		},
		args...,
	).Do(ctx)

	if err != nil {
		return err
	}

	// Try to release the remote object.
	// It will fail if the page is navigated or closed,
	// and it's okay to ignore the error in this case.
	_ = runtime.ReleaseObject(r.ObjectID).Do(ctx)

	return nil
}

var getClientRectJS = `function getClientRect() {
	const e = this.getBoundingClientRect(),
	  t = this.ownerDocument.documentElement.getBoundingClientRect();
	return {
	  x: e.left - t.left,
	  y: e.top - t.top,
	  width: e.width,
	  height: e.height,
	};
}`

func clip(X, Y, Width, Height float64) *page.Viewport {
	clip := page.Viewport{
		X:      X,
		Y:      Y,
		Width:  Width,
		Height: Height,
		Scale:  1,
	}
	return &clip
}

func scs(slide int, page *rod.Page, fn string) {
	stdo.Println(src(8), fn)
	if deb != slide {
		return
	}
	bytes, err := page.Screenshot(true, &proto.PageCaptureScreenshot{Format: proto.PageCaptureScreenshotFormatPng})
	if err == nil {
		ss(bytes).write(fn)
	}
}

func chrom() (b *rod.Browser, f func()) {
	if multiBrowser {
		ctTab, caTab := context.WithCancel(ctRoot)
		b = rod.New().ControlURL(launch().
			MustLaunch(),
		).MustConnect().
			Context(ctTab)
		f = func() {
			caTab()
			b.Close()
		}
	} else {
		b = browser
		f = func() {}
	}
	return
}
func chrome() (b *rod.Browser, f func() error) {
	if multiBrowser {
		b = rod.New().ControlURL(launch().
			MustLaunch(),
		).MustConnect().
			Context(ctRoot)
		f = b.Close
	} else {
		b = browser
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

// MouseMoveXY is an action that sends a mouse move event to the X, Y location.
func MouseMoveXY(x, y float64, opts ...dp.MouseOption) dp.MouseAction {
	return dp.ActionFunc(func(ctx context.Context) error {
		p := &input.DispatchMouseEventParams{
			Type: input.MouseMoved,
			X:    x,
			Y:    y,
		}

		// apply opts
		for _, o := range opts {
			p = o(p)
		}

		if err := p.Do(ctx); err != nil {
			return err
		}

		return p.Do(ctx)
	})
}

// MouseMoveNode is an action that dispatches a mouse move event
// at the center of a specified node.
//
// Note that the window will be scrolled if the node is not within the window's
// viewport.
func MouseMoveNode(n *cdp.Node, opts ...dp.MouseOption) dp.MouseAction {
	return dp.ActionFunc(func(ctx context.Context) error {
		t := cdp.ExecutorFromContext(ctx).(*dp.Target)
		if t == nil {
			return dp.ErrInvalidTarget
		}

		if err := dom.ScrollIntoViewIfNeeded().WithNodeID(n.NodeID).Do(ctx); err != nil {
			return err
		}

		boxes, err := dom.GetContentQuads().WithNodeID(n.NodeID).Do(ctx)
		if err != nil {
			return err
		}

		if len(boxes) == 0 {
			return dp.ErrInvalidDimensions
		}

		content := boxes[0]

		c := len(content)
		if c%2 != 0 || c < 1 {
			return dp.ErrInvalidDimensions
		}

		var x, y float64
		for i := 0; i < c; i += 2 {
			x += content[i]
			y += content[i+1]
		}
		x /= float64(c / 2)
		y /= float64(c / 2)

		return MouseMoveXY(x, y, opts...).Do(ctx)
	})
}

// MouseMove is an element query action that sends a mouse move event to the first element
// node matching the selector.
func MouseMove(sel interface{}, opts ...dp.QueryOption) dp.QueryAction {
	return dp.QueryAfter(sel, func(ctx context.Context, execCtx runtime.ExecutionContextID, nodes ...*cdp.Node) error {
		if len(nodes) < 1 {
			return fmt.Errorf("selector %q did not return any nodes", sel)
		}

		return MouseMoveNode(nodes[0]).Do(ctx)
	}, append(opts, dp.NodeVisible)...)
}

func taskKill(arg ...string) {
	cmd := exec.Command("taskKill.exe", arg...)
	stdo.Println(src(8), cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

//	func RunTO(ctx context.Context, to time.Duration, actions ...dp.Action) (cty context.Context, cay context.CancelFunc, err error) {
//		cty, cay = context.WithTimeout(ctx, to)
//		err = dp.Run(cty, actions...)
//		return
//	}
func Run(ctx context.Context, to time.Duration, actions ...dp.Action) error {
	cty, cay := context.WithTimeout(ctx, to)
	defer cay()
	return dp.Run(cty, actions...)
}
func Scanln() {
	if headLess {
		return
	}
	stdo.Print(src(8), "\nPress Enter>")
	fmt.Scanln()
}
func start(fu func(slide int), slide int, wg *sync.WaitGroup, started chan bool) {
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
	fu(slide)
}
func EmulateViewport(width, height int64, opts ...dp.EmulateViewportOption) dp.EmulateAction {
	if headLess {
		return dp.EmulateViewport(width, height, opts...)
	}
	return dp.ResetViewport()
}
func done(slide int) {
	stdo.Printf("%02d Done\n", slide)
}
func iframe(slide int, ct context.Context, url string) {
	ex(slide, dp.Run(ct,
		dp.Navigate(url+"?rs:Embed=true"),
	))
	var (
		src string
		ok  bool
	)
	ex(slide, dp.Run(ct,
		dp.WaitReady("//iframe"),
		dp.Sleep(sec),
		dp.AttributeValue("//iframe", "src", &src, &ok, dp.NodeReady),
	))
	if !ok {
		ex(slide, fmt.Errorf("no src of iframe"))
	}
	src = strings.Split(src, "&")[0]
	stdo.Println(src)
	ex(slide, dp.Run(ct,
		dp.Navigate(src),
		dp.Sleep(sec),
	))
}

func abs(i int) int {
	if i > 0 {
		return i
	}
	return -i
}
func autoStart(started chan bool, d time.Duration) {
	for len(started) > 0 {
		<-started
	}
	time.AfterFunc(d, func() {
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
	return
}
