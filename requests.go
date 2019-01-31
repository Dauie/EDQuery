package main

import (
	"encoding/json"
	"errors"
	"github.com/jroimartin/gocui"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

/*
** makeAPIRequest
** @params:
**		request - the full request url
**		store - a structure mocking the json payload you expect from the response.
** @function: pass any api request, and the appropriate structure to hold the
**            the response and the deed will be done.
*/
func makeAPIRequest(request string, store interface{}) error {
	req, err := http.NewRequest("GET", request, nil)
	if err != nil {
		return errors.New("error creating request" + err.Error())
	}
	resp, err := ClientG.Do(req)
	if err != nil {
		return errors.New("Error making API request\n" + request + "\n" + err.Error())
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("ioutil.ReadAll()" + err.Error())
	}
	err = json.Unmarshal(data, store)
	if err != nil {
		return errors.New("json.Unmarshal()" + err.Error())
	}
	return nil
}

func cmdrCreditRequest(cmdrName string, creditLog *CmdrCreditLog, g *gocui.Gui) error {
	payload := url.Values{}
	payload.Add("commanderName", cmdrName)
	apiKey, _  := CmdrMapG[cmdrName]
	payload.Add("apiKey", apiKey)
	reqStr := CmdrEndpointG + "get-credits?" + payload.Encode()
	if err := makeAPIRequest(reqStr, creditLog); err != nil {
		log.Panicln(err)
	}
	if err := printCredits(*creditLog, g); err != nil {
		log.Panicln(err)
	}
	return nil
}

func cmdrFlightLogRequest(cmdrName string, flightLog *CmdrFlightLog, g *gocui.Gui) error {
	payload := url.Values{}
	payload.Add("commanderName", cmdrName)
	apiKey, _  := CmdrMapG[cmdrName]
	payload.Add("apiKey", apiKey)
	reqStr := LogEndpointG + "get-logs?" + payload.Encode()
	if err := makeAPIRequest(reqStr, flightLog); err != nil {
		log.Panicln(err)
	}
	if err := cmdrLastPositionRequest(cmdrName, &flightLog.LastPos); err != nil {
		log.Panicln(err)
	}
	if err := printFlightLog(flightLog, g); err != nil {
		log.Panicln(err)
	}
	return nil
}

func cmdrInventoryRequest(cmdrName string, cmdrLog *CmdrLog, g *gocui.Gui) error {
	if err := cmdrInvMatRequest(cmdrName, &cmdrLog.Materials); err != nil {
		log.Panicln(err)
	}
	if err := cmdrInvDataMatReqest(cmdrName, &cmdrLog.Data); err != nil {
		log.Panicln(err)
	}
	if err := printInventory(cmdrLog, g); err != nil {
		log.Panicln(err)
	}
	return nil
}

func cmdrInvDataMatReqest(cmdrName string, dataMatLog *CmdrInventoryLog) error {
	payload := url.Values{}
	payload.Add("commanderName", cmdrName)
	apiKey, _  := CmdrMapG[cmdrName]
	payload.Add("apiKey", apiKey)
	payload.Add("type", "data")
	reqStr := CmdrEndpointG + "get-materials?" + payload.Encode()
	if err := makeAPIRequest(reqStr, dataMatLog); err != nil {
		return err
	}
	return nil
}

func cmdrInvMatRequest(cmdrName string, matLog *CmdrInventoryLog) error {
	payload := url.Values{}
	payload.Add("commanderName", cmdrName)
	apiKey, _  := CmdrMapG[cmdrName]
	payload.Add("apiKey", apiKey)
	payload.Add("type", "materials")
	reqStr := CmdrEndpointG + "get-materials?" + payload.Encode()
	if err := makeAPIRequest(reqStr, matLog); err != nil {
		log.Panicln(err)
	}
	return nil
}

func cmdrLastPositionRequest(cmdrName string, lastPos *CmdrLastPosition) error {
	payload := url.Values{}
	payload.Add("commanderName", cmdrName)
	apiKey, _  := CmdrMapG[cmdrName]
	payload.Add("apiKey", apiKey)
	reqStr := LogEndpointG + "get-position?" + payload.Encode()
	if err := makeAPIRequest(reqStr, lastPos); err != nil {
		log.Panicln(err)
	}
	return nil
}

func cmdrRankRequest(cmdrName string, rankLog *CmdrRankLog, g *gocui.Gui) error {
	payload := url.Values{}
	payload.Add("commanderName", cmdrName)
	apiKey, _  := CmdrMapG[cmdrName]
	payload.Add("apiKey", apiKey)
	reqStr := CmdrEndpointG + "get-ranks?" + payload.Encode()
	if err := makeAPIRequest(reqStr, rankLog); err != nil {
		log.Panicln(err)
	}
	if err := printRank(*rankLog, g); err != nil {
		log.Panicln(err)
	}
	return nil
}
