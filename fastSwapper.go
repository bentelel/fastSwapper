package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/oleiade/reflections"
)

type Settings struct {
	Settings       []Default        `json:"defaults"`
	ActiveSettings []ActiveSettings `json:"activesettings"`
}
type Default struct {
	Tgkdir string `json:"tgkdir"`
}
type ActiveSettings struct {
	OldDirectory string `json:"olddirectory"`
	NewDirectory string `json:"newdirectory"`
}

func main() {
	args := os.Args[1:]
	err := parseCLIargs(args)
	if err != nil {
		log.Fatal(err.Error())
	}
	set := getSettings("settings.json")

	// setSettings("settings.json", "Tgkdir", "C:\\Tagetik\\Tagetik Excel .NET Client")

	fmt.Println(GetAllInDir(set.Tgkdir))
}

func parseCLIargs(args []string) error {
	var err error
	if len(args) == 0 {
		return err
	}
	// set default path flag expects the syntax of fastSwapper -d <path to directory>
	const SET_DEFAULT_PATH_FLAG = "-d"
	if ContainsString(args, SET_DEFAULT_PATH_FLAG) {
		if len(args) < 2 {
			err = errors.New("No path provided. Use fastSwapper -d <path to default dir>.")
			return err
		}
		candidatePath := args[1]
		if !IsDir(candidatePath) {
			err = errors.New("Supplied path does not exist.")
			return err
		}
		setSettings("settings.json", "Tgkdir", candidatePath)
	}

	return err
}

func unmarshalSettingsJson(filename string) Settings {
	jsonFile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	byteResult, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	var settings Settings
	err = json.Unmarshal(byteResult, &settings)
	if err != nil {
		log.Fatal(err)
	}
	return settings
}

func updateSettingsJson(filename string, data Settings) {
	modifiedJson, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(filename, modifiedJson, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func getSettings(filename string) Default {
	return unmarshalSettingsJson(filename).Settings[0]
}

func setSettings(filename string, defaultToChange string, newValue string) {
	unmarshaledJson := unmarshalSettingsJson(filename)

	err := reflections.SetField(&unmarshaledJson.Settings[0], defaultToChange, newValue)
	if err != nil {
		log.Fatal(err)
	}
	updateSettingsJson(filename, unmarshaledJson)
}
