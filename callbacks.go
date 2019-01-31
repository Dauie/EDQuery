package main

import (
	"github.com/jroimartin/gocui"
	"log"
	"strings"
)

var (
	CmdrViewArrG = []string{"side", "flightLog", "materials", "data"}
	EditViewArrG = []string{"editCmdr", "editApi"}
	CmdrToEditG  = ""
	CmdrViewInxG = 0
	EditViewInxG = 0
)

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextCmdrView); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("side", gocui.KeyEnter, gocui.ModNone, queryCmdr); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("side", gocui.KeyCtrlE, gocui.ModNone, editCmdr); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("side", gocui.KeyCtrlS, gocui.ModNone, saveSide); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("side", gocui.KeyDelete, gocui.ModNone, deleteCmdr); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("editCmdr", gocui.KeyTab, gocui.ModNone, nextCmdrEditView); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("editCmdr", gocui.KeyCtrlE, gocui.ModNone, saveCmdrEdit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("editApi", gocui.KeyTab, gocui.ModNone, nextCmdrEditView); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("editApi", gocui.KeyCtrlE, gocui.ModNone, saveCmdrEdit); err != nil {
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

func deleteCmdr(g *gocui.Gui, v *gocui.View) error {
	cmdrName := getTrimmedLineFromCursor(v)
	delete(CmdrMapG, cmdrName)
	return nil
}

func editCmdr(g *gocui.Gui, v *gocui.View) error {
	CmdrToEditG = strings.ToUpper(getTrimmedLineFromCursor(v))
	maxX, maxY := g.Size()
	if editNameView, err := g.SetView("editCmdr", maxX/2-30, maxY/2 - 2, maxX/2+30, maxY/2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		editNameView.Title = "Edit Commander Name"
		editNameView.Editable = true
		if _, err := editNameView.Write([]byte(CmdrToEditG)); err != nil {
			log.Panicln(err)
		}
	}
	if editApiView, err := g.SetView("editApi", maxX/2-30, maxY/2 + 1, maxX/2+30, maxY/2+3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		editApiView.Title = "Edit Commander API key"
		editApiView.Editable = true
		if _, err := editApiView.Write([]byte(CmdrMapG[CmdrToEditG])); err != nil {
			log.Panicln(err)
		}
		if _, err := g.SetCurrentView("editCmdr"); err != nil {
			return err
		}
	}
	return nil
}

func saveCmdrEdit(g *gocui.Gui, v *gocui.View) error {
	editCmdrView, err := g.SetCurrentView("editCmdr")
	if err != nil {
		return err
	}
	eName := strings.ToUpper(strings.TrimSpace(editCmdrView.Buffer()))
	editApiView, err := g.SetCurrentView("editApi")
	if err != nil {
		return err
	}
	eApi := strings.TrimSpace(editApiView.Buffer())
	delete(CmdrMapG, CmdrToEditG)
	CmdrMapG[eName] = eApi
	if err := g.DeleteView("editCmdr"); err != nil {
		return err
	}
	if err := g.DeleteView("editApi"); err != nil {
		return err
	}
	if _, err := g.SetCurrentView("side"); err != nil {
		return err
	}
	if err := sideView(g); err != nil {
		return err
	}
	return nil
}

func nextCmdrEditView(g *gocui.Gui, v *gocui.View) error {
	nextIndex := (EditViewInxG + 1) % len(EditViewArrG)
	name := EditViewArrG[nextIndex]
	if _, err := g.SetCurrentView(name); err != nil {
		return err
	}
	EditViewInxG = nextIndex
	return nil
}

func nextCmdrView(g *gocui.Gui, v *gocui.View) error {
	nextIndex := (CmdrViewInxG + 1) % len(CmdrViewArrG)
	name := CmdrViewArrG[nextIndex]
	if _, err := g.SetCurrentView(name); err != nil {
		return err
	}
		CmdrViewInxG = nextIndex
	return nil
}

/*TODO: implement go routines for all api calls*/
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
	var out string

	cmdrListFd := openCmdrListFile()
	defer func() {
		err := cmdrListFd.Close()
		if err != nil {
			log.Panicln(err)
		}
	}()
	for k, v := range CmdrMapG {
		out = out + k + " " + v + "\n"
	}
	if _, err := cmdrListFd.Write([]byte(out)); err != nil {
		return err
	}
	return nil
}
