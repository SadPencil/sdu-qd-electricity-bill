/*
Copyright © 2019 Sad Pencil
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

type query_elec_roominfo struct {
	Aid      string   `json:"aid"`
	CardID   string   `json:"account"`
	Room     Room     `json:"room"`
	Floor    Floor    `json:"floor"`
	Campus   Campus   `json:"area"`
	Building Building `json:"building"`
}

type query_applist struct {
	AppType string `json:"apptype"`
}
type query_card struct {
	IdType string `json:"idtype"`
	Id     string `json:"id"`
}
type query_appinfo struct {
	Aid    string `json:"aid"`
	CardID string `json:"account"`
}
type query_elec_area struct {
	Aid    string `json:"aid"`
	CardID string `json:"account"`
}
type query_elec_building struct {
	Aid    string `json:"aid"`
	CardID string `json:"account"`
	Campus Campus `json:"area"`
}
type Floor struct {
	FloorID string `json:"floorid"`
	Floor   string `json:"floor"`
}
type Campus struct {
	Campus     string `json:"area"`
	CampusName string `json:"areaname"`
}
type Room struct {
	RoomID   string `json:"roomid"`
	RoomName string `json:"room"`
}
type Building struct {
	BuildingID   string `json:"buildingid"`
	BuildingName string `json:"building"`
}

func synjones_onecard_query(scheme, server, relativeURL, function_name string, jsonStruct map[string]interface{}, client *http.Client) (ret interface{}, err error) {
	json_byte, err := json.Marshal(jsonStruct)
	if err != nil {
		return nil, err
	}
	//	json_urlencoded := url.QueryEscape(string(json_byte))

	req, err := http.NewRequest("POST", scheme+"://"+server+relativeURL, nil)
	q := req.URL.Query()
	q.Add("jsondata", string(json_byte))
	q.Add("funname", function_name)
	q.Add("json", "true")
	bodyReader := strings.NewReader(q.Encode())

	//fmt.Println(q.Encode())

	req, err = http.NewRequest("POST", scheme+"://"+server+relativeURL, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	//fmt.Println(string(respBody))

	err = json.Unmarshal(respBody, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
func get_electricity_bill(settings Settings) (ret string, err error) {
	info, err := synjones_onecard_query(DEFAULT_SERVER_SCHEME, DEFAULT_SERVER_HOST_PORT,
		"/web/Common/Tsm.html",
		"synjones.onecard.query.elec.roominfo",
		map[string]interface{}{
			"query_elec_roominfo": query_elec_roominfo{
				Aid:    settings.SduElec.AppAid,
				CardID: settings.Account.CardID,
				Campus: Campus{
					Campus:     settings.SduElec.Campus,
					CampusName: settings.SduElec.CampusName,
				},
				Building: Building{
					BuildingName: settings.SduElec.BuildingName,
					BuildingID:   settings.SduElec.BuildingID,
				},
				Floor: Floor{},
				Room: Room{
					RoomName: settings.SduElec.RoomName,
					RoomID:   settings.SduElec.RoomID,
				},
			},
		},
		&http.Client{})
	if err != nil {
		return "", err
	}

	ret = info.(map[string]interface{})["query_elec_roominfo"].(map[string]interface{})["errmsg"].(string)
	return ret, nil
}
