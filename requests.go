package main

import (
	"github.com/jroimartin/gocui"
	"net/url"
)

func cmdrCreditRequest(cmdrName string, creditLog *CmdrCreditLog, g *gocui.Gui) error {
	payload := url.Values{}
	payload.Add("commanderName", cmdrName)
	apiKey, _  := CmdrMapG[cmdrName]
	payload.Add("apiKey", apiKey)
	reqStr := CmdrEndpointG + "get-credits?" + payload.Encode()
	if err := makeAPIRequest(reqStr, creditLog); err != nil {
		return err
	}
	if err := printCredits(*creditLog, g); err != nil {
		return err
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
		return err
	}
	if err := cmdrLastPositionRequest(cmdrName, &flightLog.LastPos); err != nil {
		return err
	}
	if err := printFlightLog(flightLog, g); err != nil {
		return err
	}
	return nil
}

func cmdrInventoryRequest(cmdrName string, cmdrLog *CmdrLog, g *gocui.Gui) error {
	if err := cmdrInvMatRequest(cmdrName, &cmdrLog.Materials); err != nil {
		return err
	}
	if err := cmdrInvDataMatReqest(cmdrName, &cmdrLog.Data); err != nil {
		return err
	}
	if err := printInventory(cmdrLog, g); err != nil {
		return err
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
		return err
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
		return err
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
		return err
	}
	if err := printRank(*rankLog, g); err != nil {
		return err
	}
	return nil
}
