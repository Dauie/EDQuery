package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
)

const EdsmApiVersionNum = "1/"

const EdsmBaseURL = "https://www.edsm.net/"

const CmdrEndpointG = EdsmBaseURL + "api-commander-v" + EdsmApiVersionNum

const SysEndpointG = EdsmBaseURL + "api-system-v" + EdsmApiVersionNum

const LogEndpointG = EdsmBaseURL + "api-logs-v" + EdsmApiVersionNum

var ClientG http.Client

var CmdrListPathG string

var AppDataPathG string

var CmdrMapG = map[string]string{}

/*
** common error responses from edsm.net
*/
var EDSMErrorsG = map[int]string {
	201: "missing cmdr/api key",
	203: "commander not found",
	204: "item type not available",
	208: "no credits stored",
	207: "no rank stored",
}

/*
** openCmdrListFile
** @params: none
** @function: opens ~/.cmdr/cmdrlist for saving
*/
func openCmdrListFile() (CmdrListFd *os.File) {
	CmdrListFd, err := os.OpenFile(CmdrListPathG, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err == io.EOF {
		return nil
	} else if err != nil {
		log.Panicln(err)
	}
	return CmdrListFd
}

/*
** getTrimmedLineFromCursor
** @params: v - the view you want the line from
** @function: retrieves the string the cursor is currently on and removes whitespace.
*/
func getTrimmedLineFromCursor(v *gocui.View) string {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}
	return strings.TrimSpace(l)
}

/* initAppdata
** @param: none
** @function: opens a file with saved cmdr/api pairs, and reads them into a map
*/
func initAppdata() error {
	homePath, ok := os.LookupEnv("HOME")
	if !ok {
		return errors.New("HOME env variable not set")
	}
	AppDataPathG = homePath + "/.cmdr"
	CmdrListPathG = AppDataPathG + "/cmdrlist"
	_, err := os.Stat(AppDataPathG)
	if os.IsNotExist(err) {
		err = os.MkdirAll(AppDataPathG, 0755)
		if err != nil {
			log.Fatalln(err)
		}
	}
	file, _ := os.Open(CmdrListPathG)
	fscanner := bufio.NewScanner(file)
	for fscanner.Scan() {
		keyVal := strings.Split(fscanner.Text(), " ")
		if len(keyVal) == 1 {
			CmdrMapG[keyVal[0]] = "API_KEY_MISSING"
		} else {
			CmdrMapG[keyVal[0]] = keyVal[1]
		}
	}
	return nil
}

/*
** initGui
** @param: none
** @function: initialize our gocui Gui and set global attributes
*/
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

	ClientG = http.Client {
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
