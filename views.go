package main

import (
	"errors"
	"github.com/jroimartin/gocui"
)

func checkDimensions(mY int, mX int) error {
	if mY < 10 || mX < 10 {
		return errors.New("window dimensions not in bounds")
	}
	return nil
}

func creditView(g *gocui.Gui) error {
	mX, mY := g.Size()
	if err := checkDimensions(mX, mY); err != nil {
		return nil
	}
	view, err := g.SetView("credits", mX / 5 + 1, (mY / 10) * 3 + 1, (mX / 5) * 2 , (mY / 10) * 6)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		view.Frame = true
		view.Title = "Credits"
	}
	return nil
}

func flightLogView(g *gocui.Gui) error {
	mX, mY := g.Size()
	if err := checkDimensions(mX, mY); err != nil {
		return nil
	}
	view, err := g.SetView("flightLog", mX / 5 + 1, (mY / 10) * 6 + 1, (mX / 5) * 2, mY)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		view.Wrap = true
		view.Highlight = true
		view.Frame = true
		view.Title = "Flight Log"
	}
	return nil
}

func inventoryView(g *gocui.Gui) error {
	mX, mY := g.Size()
	if err := checkDimensions(mX, mY); err != nil {
		return nil
	}
	materialXStart := (mX / 5) * 2 + 1
	invWid := (mX / 5) * 3
	materialXEnd := materialXStart + (invWid / 2)
	dataXStart := materialXEnd + 1
	dataXEnd := mX - 1
	if err := checkDimensions(mX, mY); err != nil {
		return nil
	}
	materialsView, err := g.SetView("materials", materialXStart, 1, materialXEnd, mY)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		materialsView.Wrap = true
		materialsView.Highlight = true
		materialsView.Frame = true
		materialsView.Title = "Materials"
	}
	dataView, err := g.SetView("data", dataXStart, 1, dataXEnd, mY)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		dataView.Wrap = true
		dataView.Highlight = true
		dataView.Frame = true
		dataView.Title = "Data"
	}
	return nil
}

func rankView(g *gocui.Gui) error {
	mX, mY := g.Size()
	if err := checkDimensions(mX, mY); err != nil {
		return nil
	}
	view, err := g.SetView("rank", mX / 5 + 1, 1, (mX / 5) * 2, (mY / 10) * 3)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		view.Frame = true
		view.Wrap = true
		view.Title = "Rank"
	}
	return nil
}

func sideView(g *gocui.Gui) error {
	mX, mY := g.Size()
	if err := checkDimensions(mX, mY); err != nil {
		return nil
	}
	view, err := g.SetView("side", 0, 0, mX / 5, mY)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		view.Frame = true
		view.Title = "CMDRs"
		view.Highlight = true
		view.SelBgColor = gocui.ColorGreen
		view.SelFgColor = gocui.ColorBlack

		var cmdrList string
		for k := range CmdrMap {
			cmdrList = cmdrList + k + "\n"
		}
		_, err = view.Write([]byte(cmdrList))
		if err != nil {
			return err
		}
		if _, err := g.SetCurrentView("side"); err != nil {
			return err
		}
	}
	return nil
}

func layout(g *gocui.Gui) error {
	if err := sideView(g); err != nil {
		return err
	}
	if err := rankView(g); err != nil {
		return err
	}
	if err := creditView(g); err != nil {
		return err
	}
	if err := flightLogView(g); err != nil {
		return err
	}
	if err := inventoryView(g); err != nil {
		return err
	}

	return nil
}
