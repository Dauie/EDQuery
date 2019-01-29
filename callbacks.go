package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"io"
	"log"
)

var (
	viewArr = []string{"side", "flightLog", "materials", "data"}
	active  = 0
)


func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("side", gocui.KeyEnter, gocui.ModNone, queryCmdr); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("side", gocui.KeyCtrlS, gocui.ModNone, saveSide); err != nil {
		log.Panicln(err)
	}
	return nil
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func getLine(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	maxX, maxY := g.Size()
	if v, err := g.SetView("msg", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, l)
		if _, err := g.SetCurrentView("msg"); err != nil {
			return err
		}
	}
	return nil
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	nextIndex := (active + 1) % len(viewArr)
	name := viewArr[nextIndex]
	if _, err := g.SetCurrentView(name); err != nil {
		return err
	}
		active = nextIndex
	return nil
}

func queryCmdr(g *gocui.Gui, v *gocui.View) error {
	var err error
	var cmdr CmdrLog

	_, cy := v.Cursor()
	cmdr.Name, err = v.Line(cy)
	if err != nil {
		return err
	}
	if err := cmdrRankRequest(cmdr.Name, &cmdr.Rank, g); err != nil {
		return err
	}
	if err := cmdrCreditRequest(cmdr.Name, &cmdr.Credits, g); err != nil {
		return err
	}

	if err := cmdrFlightLogRequest(cmdr.Name, &cmdr.FlightLog, g); err != nil {
		return err
	}

	if err := cmdrInventoryRequest(cmdr.Name, &cmdr, g); err != nil {
		return err
	}
	_, _ = g.SetCurrentView("side")
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	_ = saveSide(g, v)
	return gocui.ErrQuit
}

func saveSide(g *gocui.Gui, v *gocui.View) error {
	p := make([]byte, 4096)
	v.Rewind()
	for {
		n, err := v.Read(p)
		if n > 0 {
			if _, err := CmdrListFd.Write(p[:n]); err != nil {
				return err
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}
