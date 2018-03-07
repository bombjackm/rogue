package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"github.com/mattn/go-runewidth"
)

var row = 0
var style = tcell.StyleDefault

func putln(s tcell.Screen, str string) {

	puts(s, style, 1, row, str)
	row++
}

func puts(s tcell.Screen, style tcell.Style, x, y int, str string) {
	i := 0
	var deferred []rune
	dwidth := 0
	for _, r := range str {
		switch runewidth.RuneWidth(r) {
		case 0:
			if len(deferred) == 0 {
				deferred = append(deferred, ' ')
				dwidth = 1
			}
		case 1:
			if len(deferred) != 0 {
				s.SetContent(x+i, y, deferred[0], deferred[1:], style)
				i += dwidth
			}
			deferred = nil
			dwidth = 1
		case 2:
			if len(deferred) != 0 {
				s.SetContent(x+i, y, deferred[0], deferred[1:], style)
				i += dwidth
			}
			deferred = nil
			dwidth = 2
		}
		deferred = append(deferred, r)
	}
	if len(deferred) != 0 {
		s.SetContent(x+i, y, deferred[0], deferred[1:], style)
		i += dwidth
	}
}

func main() {

	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	encoding.Register()

	if e = s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	//plain := tcell.StyleDefault
	//bold := style.Bold(true)
	black := tcell.NewHexColor(0x171817)
	white := tcell.NewHexColor(0xdad6c6)

	s.SetStyle(tcell.StyleDefault.
		Foreground(white).
		Background(black))
	s.Clear()

	w, h := s.Size()
	//putln(s, strconv.Itoa(w)+"x"+strconv.Itoa(h))
	x := w / 2
	y := h / 2

	//quit := make(chan struct{})

	s.Show()
	quit := false
	//go func() {
	for !quit {
		s.SetContent(x, y, '@', nil, tcell.StyleDefault)
		s.SetContent(1, 1, 'x', nil, tcell.StyleDefault)
		ev := s.PollEvent()
		s.SetContent(1, 1, 'o', nil, tcell.StyleDefault)
		s.Show()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape:
				//close(quit)
				//return
				quit = true
			case tcell.KeyCtrlL:
				s.Sync()
			case tcell.KeyLeft:
				putln(s, "LEFTD")
				x--
			case tcell.KeyRight:
				x++
			case tcell.KeyUp:
				y--
			case tcell.KeyDown:
				y++
			}
		case *tcell.EventResize:
			s.Sync()
		}
	}
	//}()

	//<-quit

	s.Fini()
}
