package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
	"sort"
	"strconv"
	"strings"
)

func printCredits(credits CmdrCreditLog, g *gocui.Gui) error {
	var out string

	view, err := g.SetCurrentView("credits")
	if err != nil {
		return err
	}
	view.Clear()
	if credits.Msgnum == 100 {
		recentLog := credits.Credits[0]
		out = fmt.Sprintf("%d © ", recentLog.Balance)
	} else {
		out = fmt.Sprintf("%s", EDSMErrorsG[credits.Msgnum - 1])
	}
	buf := []byte(out)
	if _, err := view.Write(buf); err != nil {
		return err
	}
	return nil
}

func printFlightLog(log *CmdrFlightLog, g *gocui.Gui) error {
	var out string

	view, err := g.SetCurrentView("flightLog")
	if err != nil {
		return err
	}
	view.Clear()
	if log.LastPos.Msgnum == 100 {
		out = "Last seen: " + log.LastPos.System + "\n"
	} else {
		out = "Position unknown\n"
	}
	if log.Msgnum == 100 {
		for _, v := range log.Logs {
			dateTime := strings.Split(v.Date, " ")
			out = out + v.System + " - " + dateTime[0] + "\n"
		}
	} else {
		out = out + fmt.Sprintf("%s", EDSMErrorsG[log.Msgnum - 1])
	}
	buf := []byte(out)
	if _, err := view.Write(buf); err != nil {
		return err
	}
	return nil
}

func printInventory(log *CmdrLog, g *gocui.Gui) error {
	var materials string
	var data string

	matView, err := g.SetCurrentView("materials")
	if err != nil {
		return err
	}
	matView.Clear()
	if log.Materials.Msgnum == 100 {
		if len(log.Materials.Items) == 0 {

		} else {
			for _, v := range log.Materials.Items {
				materials = materials + v.Name + "\t x " + strconv.Itoa(v.Qty) + "\n"
			}
		}
	} else {
		materials = EDSMErrorsG[log.Materials.Msgnum] + "\n"
	}
	matBuf := []byte(materials)
	if _, err := matView.Write(matBuf); err != nil {
		return err
	}
	dataView, err := g.SetCurrentView("data")
	if err != nil {
		return err
	}
	dataView.Clear()
	if log.Data.Msgnum == 100 {
		if len(log.Data.Items) == 0 {
			data = "No data materials on record"
		} else {
			for _, v := range log.Data.Items {
				data = data + v.Name + "\t x " + strconv.Itoa(v.Qty) + "\n"}
		}
	} else {
		data = EDSMErrorsG[log.Data.Msgnum] + "\n"
	}
	if _, err := dataView.Write([]byte(data)); err != nil {
		return err
	}
	return nil
}


func printRank(rank CmdrRankLog, g *gocui.Gui) error {
	var out string
	view, err := g.SetCurrentView("rank")
	if err != nil {
		return err
	}
	view.Clear()
	if rank.Msgnum == 100 {
		out = fmt.Sprintf("Trade Rank: %s(%%%d)\n" +
			"Explorer Rank: %s(%%%d)\n" +
			"Combat Rank: %s(%%%d)\n\n",
			rank.Verbose.Trade, rank.Progress.Trade,
			rank.Verbose.Explore, rank.Progress.Explore,
			rank.Verbose.Combat, rank.Progress.Combat)
	} else {
		out = fmt.Sprintf("%s", EDSMErrorsG[rank.Msgnum - 1])
	}
	if _, err := view.Write([]byte(out)); err != nil {
		return err
	}
	return nil
}

func printCmdrList(view *gocui.View) {
	view.Clear()
	var cmdrArr []string
	for k := range CmdrMapG {
		cmdrArr = append(cmdrArr, k)
	}
	sort.Strings(cmdrArr)
	var cmdrStr string
	for e := range cmdrArr {
		cmdrStr = cmdrStr + cmdrArr[e] + "\n"
	}
	_, err := view.Write([]byte(cmdrStr))
	if err != nil {
		log.Panicln(err)
	}
	_, err = view.Write([]byte("NEW_CMDR"))
	if err != nil {
		log.Panicln(err)
	}
}
