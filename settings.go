/*
Copyright © 2019 Sad Pencil
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Account struct {
	UID      string
	CardID   string
	UserName string
}

type Log struct {
	Filename string
}

type Control struct {
	Interval int32
}

type SduElec struct {
	ServerHostPort string
	ServerScheme   string
	AppAid         string
	AppName        string
	Campus         string
	CampusName     string
	BuildingName   string
	BuildingID     string
	RoomID         string
	RoomName       string
}

type Settings struct {
	Account Account
	//Log     Log
	//Control Control
	SduElec SduElec
}

func NewSettings() Settings {
	return Settings{
		SduElec: SduElec{ServerScheme: DEFAULT_SERVER_SCHEME, ServerHostPort: DEFAULT_SERVER_HOST_PORT},
		//Control: Control{Interval: DEFAULT_INTERVAL},
	}
}

// LoadSettings -- Load settings from config file
func LoadSettings(configPath string) (settings Settings, err error) {
	// LoadSettings from config file
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error occurs while reading config file.")
		return NewSettings(), err
	}

	err = json.Unmarshal(file, &settings)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error occurs while unmarshalling config file.")
		return NewSettings(), err
	}

	return settings, nil
}
func SaveSettings(configPath string, setting Settings) (err error) {
	f, err := os.Create(configPath)
	defer f.Close()
	if err != nil {
		return err
	}
	jsonBytes, err := json.Marshal(setting)
	if err != nil {
		return err
	}
	_, err = f.WriteString(string(jsonBytes))
	if err != nil {
		return err
	}

	return nil
}
