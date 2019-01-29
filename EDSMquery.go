package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jroimartin/gocui"
)

const versionNum = "1/"

/* I literally dont care. */
const apiKey = "15a77f8df4b47582347aa83bb29f4245335a334a"

const baseURL = "https://www.edsm.net/"

const cmdrEndpoint = baseURL + "api-commander-v" + versionNum

const sysEndpoint = baseURL + "api-system-v" + versionNum

const logEndpoint = baseURL + "api-logs-v" + versionNum

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
	//if _, err := gui.SetCurrentView("side"); err != nil {
	//	if err != gocui.ErrUnknownView {
	//		log.Fatalln(err)
	//	}
	//}
	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Fatalln(err)
	}
}
