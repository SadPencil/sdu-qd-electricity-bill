/*
Copyright © 2019 Sad Pencil
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func cartman() {
	newSetting := NewSettings()
	reader := bufio.NewReader(os.Stdin)
	var err error

	{
		fmt.Print("输入学工号: ")
		uidStr, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		uidStr = strings.TrimSpace(uidStr)

		cardinfo, err := synjones_onecard_query(DEFAULT_SERVER_SCHEME, DEFAULT_SERVER_HOST_PORT,
			"/web/Common/Tsm.html",
			"synjones.onecard.query.card",
			map[string]interface{}{
				"query_card": query_card{IdType: "sno", Id: uidStr},
			}, &http.Client{})
		if err != nil {
			panic(err)
		}

		newSetting.Account.UID = cardinfo.(map[string]interface{})["query_card"].(map[string]interface{})["card"].([]interface{})[0].(map[string]interface{})["sno"].(string)
		newSetting.Account.UserName = cardinfo.(map[string]interface{})["query_card"].(map[string]interface{})["card"].([]interface{})[0].(map[string]interface{})["name"].(string)
		newSetting.Account.CardID = cardinfo.(map[string]interface{})["query_card"].(map[string]interface{})["card"].([]interface{})[0].(map[string]interface{})["account"].(string)

		fmt.Println("学工号: ", newSetting.Account.UID)
		fmt.Println("姓名: ", newSetting.Account.UserName)
		fmt.Println("卡号: ", newSetting.Account.CardID)
	}
	{
		applist, err := synjones_onecard_query(DEFAULT_SERVER_SCHEME, DEFAULT_SERVER_HOST_PORT,
			"/web/NetWork/AppList.html",
			"synjones.onecard.query.applist",
			map[string]interface{}{
				"query_applist": query_applist{AppType: "elec"},
			}, &http.Client{})
		if err != nil {
			panic(err)
		}

		appSelect := make([]struct {
			name string
			aid  string
		}, 0)

		for _, v := range applist.(map[string]interface{})["query_applist"].(map[string]interface{})["applist"].([]interface{}) {
			appSelect = append(appSelect, struct {
				name string
				aid  string
			}{
				name: v.(map[string]interface{})["name"].(string),
				aid:  v.(map[string]interface{})["aid"].(string),
			})
		}

		for k, v := range appSelect {
			fmt.Println("[" + fmt.Sprint(k) + "]\t" + v.name)
		}
		fmt.Print("请选择: ")
		choice, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		choice = strings.TrimSpace(choice)

		choiceID, err := strconv.Atoi(choice)
		if err != nil {
			panic(err)
		}
		newSetting.SduElec.AppAid = appSelect[choiceID].aid
		newSetting.SduElec.AppName = appSelect[choiceID].name
	}
	{
		campus, err := synjones_onecard_query(DEFAULT_SERVER_SCHEME, DEFAULT_SERVER_HOST_PORT,
			"/web/Common/Tsm.html",
			"synjones.onecard.query.elec.area",
			map[string]interface{}{
				"query_elec_area": query_elec_area{Aid: newSetting.SduElec.AppAid, CardID: newSetting.Account.CardID},
			}, &http.Client{})
		if err != nil {
			panic(err)
		}

		campusSelect := make([]struct {
			area     string
			areaname string
		}, 0)

		for _, v := range campus.(map[string]interface{})["query_elec_area"].(map[string]interface{})["areatab"].([]interface{}) {
			campusSelect = append(campusSelect, struct {
				area     string
				areaname string
			}{area: v.(map[string]interface{})["area"].(string), areaname: v.(map[string]interface{})["areaname"].(string)})
		}

		for k, v := range campusSelect {
			fmt.Println("[" + fmt.Sprint(k) + "]\t" + v.areaname)
		}
		fmt.Print("请选择: ")
		choice, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		choice = strings.TrimSpace(choice)

		choiceID, err := strconv.Atoi(choice)
		if err != nil {
			panic(err)
		}

		newSetting.SduElec.Campus = campusSelect[choiceID].area
		newSetting.SduElec.CampusName = campusSelect[choiceID].areaname
	}
	{
		building, err := synjones_onecard_query(DEFAULT_SERVER_SCHEME, DEFAULT_SERVER_HOST_PORT,
			"/web/Common/Tsm.html",
			"synjones.onecard.query.elec.building",
			map[string]interface{}{
				"query_elec_building": query_elec_building{
					Aid:    newSetting.SduElec.AppAid,
					CardID: newSetting.Account.CardID,
					Campus: Campus{
						Campus:     newSetting.SduElec.Campus,
						CampusName: newSetting.SduElec.CampusName,
					},
				},
			},
			&http.Client{})
		if err != nil {
			panic(err)
		}

		buildingSelect := make([]struct {
			building   string
			buildingid string
		}, 0)
		for _, v := range building.(map[string]interface{})["query_elec_building"].(map[string]interface{})["buildingtab"].([]interface{}) {
			buildingSelect = append(buildingSelect, struct {
				building   string
				buildingid string
			}{building: v.(map[string]interface{})["building"].(string), buildingid: v.(map[string]interface{})["buildingid"].(string)})
		}
		for k, v := range buildingSelect {
			fmt.Println("[" + fmt.Sprint(k) + "]\t" + v.building)
		}
		fmt.Print("请选择: ")
		choice, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		choice = strings.TrimSpace(choice)

		choiceID, err := strconv.Atoi(choice)
		if err != nil {
			panic(err)
		}

		newSetting.SduElec.BuildingID = buildingSelect[choiceID].buildingid
		newSetting.SduElec.BuildingName = buildingSelect[choiceID].building
	}
	{
		fmt.Print("输入房间号（例：A101）: ")
		inputstr, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		inputstr = strings.TrimSpace(inputstr)
		newSetting.SduElec.RoomID = inputstr
		newSetting.SduElec.RoomName = inputstr
	}
	{
		ret, err := get_electricity_bill(newSetting)
		if err != nil {
			panic(err)
		}
		fmt.Println(ret)
	}

	err = SaveSettings(DEFAULT_CONFIG_FILENAME, newSetting)
	if err != nil {
		panic(err)
	}
}
