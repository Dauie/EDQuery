package main

import (
	"github.com/jroimartin/gocui"
	"io"
	"log"
	"os"
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

func openCmdrListFile() (CmdrListFd *os.File) {
	CmdrListFd, err := os.OpenFile(CmdrListPath, os.O_CREATE|os.O_WRONLY, 0755)
	if err == io.EOF {
		return nil
	} else if err != nil {
		log.Panicln(err)
	}
	return CmdrListFd
}

func getLine(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}
	maxX, maxY := g.Size()
	if editNameView, err := g.SetView("editCmdr", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		editNameView.Title = "Edit Commander Name"
		editNameView.Editable = true
		if _, err := editNameView.Write([]byte(l)); err != nil {
			log.Panicln(err)
		}
	}
	if editApiView, err := g.SetView("editApi", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		editApiView.Title = "Edit Commander API key"
		editApiView.Editable = true
		if _, err := editApiView.Write([]byte(CmdrMap[l])); err != nil {
			log.Panicln(err)
		}
		if _, err := g.SetCurrentView("editCmdr"); err != nil {
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
		log.Panicln(err)
	}
	if err := cmdrRankRequest(cmdr.Name, &cmdr.Rank, g); err != nil {
		log.Panicln(err)
	}
	if err := cmdrCreditRequest(cmdr.Name, &cmdr.Credits, g); err != nil {
		log.Panicln(err)
	}

	if err := cmdrFlightLogRequest(cmdr.Name, &cmdr.FlightLog, g); err != nil {
		log.Panicln(err)
	}

	if err := cmdrInventoryRequest(cmdr.Name, &cmdr, g); err != nil {
		log.Panicln(err)
	}
	_, _ = g.SetCurrentView("side")
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	_ = saveSide(g, v)
	return gocui.ErrQuit
}

func saveSide(g *gocui.Gui, v *gocui.View) error {
	var out string

	cmdrListFd := openCmdrListFile()
	defer func() {
		err := cmdrListFd.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()
	for k, v := range CmdrMap {
		out = out + k + " " + v + "\n"
	}
	if _, err := cmdrListFd.Write([]byte(out)); err != nil {
		log.Panicln(err)
	}
	return nil
}
