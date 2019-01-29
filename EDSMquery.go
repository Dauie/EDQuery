package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/jroimartin/gocui"
)

const versionNum = "1/"

const baseURL = "https://www.edsm.net/"

const cmdrEndpoint = baseURL + "api-commander-v" + versionNum

const sysEndpoint = baseURL + "api-system-v" + versionNum

const logEndpoint = baseURL + "api-logs-v" + versionNum

type RankCurrent		struct {
	Combat				int `json:"Combat"`
	Trade				int `json:"Trade"`
	Explore				int `json:"Explore"`
	CQC					int `json:"CQC"`
	Federation			int `json:"Federation"`
	Empire				int `json:"Empire"`
}

type RankProgress		struct {
	Combat				int `json:"Combat"`
	Trade				int `json:"Trade"`
	Explore				int `json:"Explore"`
	CQC					int `json:"CQC"`
	Federation			int `json:"Federation"`
	Empire				int `json:"Empire"`
}

type RankVerbose		struct {
	Combat				string `json:"Combat"`
	Trade				string `json:"Trade"`
	Explore				string `json:"Explore"`
	CQC					string `json:"CQC"`
	Federation			string `json:"Federation"`
	Empire				string `json:"Empire"`
}

type CmdrRankLog 		struct {
	Msgnum				int `json:"msgnum"`
	Msg					string `json:"msg"`
	Current				RankCurrent `json:"ranks"`
	Progress			RankProgress `json:"progress"`
	Verbose				RankVerbose `json:"ranksVerbose"`
}

type Credits 			struct {
	Balance				int `json:"balance"`
	Loan				int `json:"loan"`
	Date				string `json:"data"`
}

type CmdrCreditLog		struct {
	Msgnum				int `json:"msgnum"`
	Msg					string `json:"msg"`
	Credits				[]Credits `json:"credits"`
}

type Item				struct {
	Type				string `json:"type"`
	Name				string `json:"name"`
	Qty					int `json:"qty"`
}

type CmdrInventoryLog	struct {
	Msgnum				int `json:"msgnum"`
	Msg					string `json:"msg"`
	Items				[]Item `json:"materials"`
}

type FlightLogEntry		struct {
	ShipId				int `json:"shipId"`
	System				string `json:"system"`
	SystemId			int `json:"systemId"`
	FirstDiscover		bool `json:"firstDiscover"`
	Date				string `json:"date"`
}

type CmdrFlightLog		struct {
	Msgnum				int `json:"msgnum"`
	Msg					string `json:"msg"`
	StartDate			string `json:"startDateTime"`
	EndDate				string `json:"endDateTime"`
	Logs				[]FlightLogEntry `json:"logs"`
	LastPos				CmdrLastPosition
}

type CmdrLastPosition	struct {
	Msgnum				int `json:"msgnum"`
	Msg					string `json:"msg"`
	System				string `json:"system"`
	SystemId			int `json:"systemId"`
	FirstDiscover		bool `json:"firstDiscover"`
	Date				string `json:"date"`
	Profile				string `json:"url"`
}

type CmdrLog			struct {
	Name				string
	Rank				CmdrRankLog
	Credits				CmdrCreditLog
	Data				CmdrInventoryLog
	Materials			CmdrInventoryLog
	FlightLog			CmdrFlightLog
}

var CmdrListFd *os.File

var Client http.Client

var CmdrList string

var CmdrPath string

var EDSMErrors = map[int]string {
	201: "Missing CMDR/API key",
	203: "Commander not found",
	204: "Item type not available",
	208: "No credits stored",
	207: "No rank stored",
}

func cmdrFlightLogRequest(cmdrName string, flightLog *CmdrFlightLog, g *gocui.Gui) error {
	payload := url.Values{}
	payload.Add("commanderName", cmdrName)
	reqStr := logEndpoint + "get-logs?" + payload.Encode()
	if err := makeEDSMRequest(reqStr, flightLog); err != nil {
		return err
	}
	if err := cmdrLastPositionRequest(cmdrName, &flightLog.LastPos); err != nil {
		return err
	}
	return nil
}

func cmdrLastPositionRequest(cmdrName string, lastPos *CmdrLastPosition) error {
	payload := url.Values{}
	payload.Add("commanderName", cmdrName)
	reqStr := logEndpoint + "get-position?" + payload.Encode()
	if err := makeEDSMRequest(reqStr, lastPos); err != nil {
		return err
	}
	return nil
}

func cmdrRankRequest(cmdrName string, rankLog *CmdrRankLog, g *gocui.Gui) error {
	payload := url.Values{}
	payload.Add("commanderName", cmdrName)
	reqStr := cmdrEndpoint + "get-ranks?" + payload.Encode()
	if err := makeEDSMRequest(reqStr, rankLog); err != nil {
		return err
	}
	if err := printRank(*rankLog, g); err != nil {
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
		out = fmt.Sprintf("%s", EDSMErrors[rank.Msgnum - 1])
	}
	buf := []byte(out)
	if _, err := view.Write(buf); err != nil {
		return err
	}
	return nil
}

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
		out = fmt.Sprintf("%s", EDSMErrors[credits.Msgnum - 1])
	}
	buf := []byte(out)
	if _, err := view.Write(buf); err != nil {
		return err
	}
	return nil
}


func cmdrCreditRequest(cmdrName string, creditLog *CmdrCreditLog, g *gocui.Gui) error {
	payload := url.Values{}
	payload.Add("commanderName", cmdrName)
	reqStr := cmdrEndpoint + "get-credits?" + payload.Encode()
	if err := makeEDSMRequest(reqStr, creditLog); err != nil {
		return err
	}
	if err := printCredits(*creditLog, g); err != nil {
		return err
	}
	return nil
}

func makeEDSMRequest(request string, store interface{}) error {
	req, err := http.NewRequest("GET", request, nil)
	if err != nil {
		log.Panicln("Error creating request", err)
	}
	resp, err := Client.Do(req)
	if err != nil {
		log.Panicln("Error making commander info request", err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panicln("ioutil.ReadAll()", err)
	}
	err = json.Unmarshal(data, store)
	if err != nil {
		log.Panicln("json.Unmarshal()", err)
	}
	return nil
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	if v == nil || v.Name() == "side" {
		_, err := g.SetCurrentView("main")
		return err
	}
	_, err := g.SetCurrentView("side")
	return err
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
	_, _ = g.SetCurrentView("side")
	return nil
}


func quit(g *gocui.Gui, v *gocui.View) error {
	_ = saveSide(g, v)
	return gocui.ErrQuit
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", gocui.KeyEnter, gocui.ModNone, queryCmdr); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", gocui.KeyCtrlS, gocui.ModNone, saveSide); err != nil {
		return err
	}
	return nil
}

func checkDeminsions(mY int, mX int) error {
	if mY < 10 || mX < 10 {
		return errors.New("not in bounds")
	}
	return nil
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

func sideView(g *gocui.Gui) error {
	mX, mY := g.Size()
	if err := checkDeminsions(mX, mY); err != nil {
		return nil
	}
	view, err := g.SetView("side", 0, 0, mX / 5, mY)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		view.Frame = true
		view.Title = "CMDRs"
		view.Autoscroll = true
		view.Editable = true
		view.SelBgColor = gocui.ColorGreen
		view.SelFgColor = gocui.ColorBlack
		buff, err := ioutil.ReadFile(CmdrList)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		_, err = view.Write(buff)
		if err != nil {
			return err
		}
	}
	return nil
}

func rankView(g *gocui.Gui) error {
	mX, mY := g.Size()
	if err := checkDeminsions(mX, mY); err != nil {
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

func creditView(g *gocui.Gui) error {
	mX, mY := g.Size()
	if err := checkDeminsions(mX, mY); err != nil {
		return nil
	}
	view, err := g.SetView("credits", mX / 5 + 1, (mY / 10) * 3 + 1, (mX / 5) * 2 , (mY / 10) * 6)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		view.Frame = true
		view.Wrap = true
		view.Title = "Credits"
	}
	return nil
}

func flightLogView(g *gocui.Gui) error {
	mX, mY := g.Size()
	if err := checkDeminsions(mX, mY); err != nil {
		return nil
	}
	view, err := g.SetView("flightLog", mX / 5 + 1, (mY / 10) * 6 + 1, (mX / 5) * 2, mY)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		view.Frame = true
		view.Title = "Flight Log"
	}
	return nil
}

func inventoryView(g *gocui.Gui) error {
	mX, mY := g.Size()
	if err := checkDeminsions(mX, mY); err != nil {
		return nil
	}
	materialXStart := (mX / 5) * 2 + 1
	invWid := (mX / 5) * 3
	materialXEnd := materialXStart + (invWid / 2)
	dataXStart := materialXEnd + 1
	dataXEnd := mX - 1
	if err := checkDeminsions(mX, mY); err != nil {
		return nil
	}
	materialsView, err := g.SetView("materials", materialXStart, 1, materialXEnd, mY)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		materialsView.Frame = true
		materialsView.Title = "Materials"
	}
	dataView, err := g.SetView("data", dataXStart, 1, dataXEnd, mY)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		dataView.Frame = true
		dataView.Title = "Data"
	}
	return nil
}

func mainView(g *gocui.Gui) error {
	mX, mY := g.Size()
	if err := checkDeminsions(mX, mY); err != nil {
		return nil
	}
	view, err := g.SetView("main", mX / 5, 0, mX, mY)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		view.Frame = true
		view.Title = "Info"
	}
	return nil
}

func layout(g *gocui.Gui) error {
	if err := sideView(g); err != nil {
		return err
	}
	if err := mainView(g); err != nil {
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
	if _, err := g.SetCurrentView("side"); err != nil {
		return err
	}
	return nil
}

func initAppdata() error {
	homePath, ok := os.LookupEnv("HOME")
	if !ok {
		return errors.New("HOME env variable not set")
	}
	CmdrPath = homePath + "/.cmdr"
	CmdrList = CmdrPath + "/cmdrlist"
	_, err := os.Stat(CmdrPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(CmdrPath, 0755)
		if err != nil {
			return err
		}
	}
	CmdrListFd, err = os.OpenFile(CmdrList, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	return nil
}

func initGui() (*gocui.Gui, error) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return nil, err
	}
	g.Cursor = true
	g.Mouse = true
	return g, nil
}

func main() {

	Client = http.Client {
		Timeout: time.Second * 10,
	}

	if err := initAppdata(); err != nil {
		log.Fatalln("Error initializing application data.", err)
	}
	defer func() {
		err := CmdrListFd.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()
	gui, err := initGui()
	defer gui.Close()
	if err != nil {
		log.Fatalln(err)
	}
	gui.SetManagerFunc(layout)
	if err := keybindings(gui); err != nil {
		log.Fatalln(err)
	}
	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Fatalln(err)
	}
}

