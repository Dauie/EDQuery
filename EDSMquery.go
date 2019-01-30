package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
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

var Client http.Client

var CmdrListPath string

var AppDataPath string

var CmdrMap = map[string]string{}

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
	AppDataPath = homePath + "/.cmdr"
	CmdrListPath = AppDataPath + "/cmdrlist"
	_, err := os.Stat(AppDataPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(AppDataPath, 0755)
		if err != nil {
			log.Fatalln(err)
		}
	}
	file, _ := os.Open(CmdrListPath)
	fscanner := bufio.NewScanner(file)
	for fscanner.Scan() {
		keyVal := strings.Split(fscanner.Text(), " ")
		if len(keyVal) == 1 {
			CmdrMap[keyVal[0]] = "No API key"
		} else {
			CmdrMap[keyVal[0]] = keyVal[1]
		}
	}
	if err := fscanner.Err(); err != nil {
		log.Panicln(err)
	}
	if err != nil {
		log.Panicln(err)
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
