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
	Tgkdir    string `json:"tgkdir"`
	Tgkfolder string `json:"tgkfolder"`
}
type ActiveSettings struct {
	OldDirectory string `json:"olddirectory"`
}

type helpInformation struct {
	availableFlagsWithDesc map[string]string
}

func HelpInformation() helpInformation {
	var help helpInformation
	help.availableFlagsWithDesc = map[string]string{
		"-d":  "Set default tagetik directory > fastSwapper -d <absolute path to directory>",
		"-dw": "Reset default tagetik directory to the Windows one.",
		"-tf": "Set Tagetik Addin Folder Name",
		"-o":  "Set the name of the old directory, under this name the current Addin will be saved on swap. > fastSwapper -o <name of directory you want>",
		"-h":  "Displays this help, use > fastSwapper -h <some other flag> to display only the help for a specific flag.",
		"-sw": "Swap directories..",
	}
	return help
}

func main() {
	args := os.Args[1:]
	err := parseCLIargs(args)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func parseCLIargs(args []string) error {
	var err error
	const SETTINGS_FILE_NAME string = "settings.json"
	const TGK_DIR_DEFAULT_WIN = "C:\\Tagetik\\Tagetik Excel .NET Client"
	FORBIDDEN_CHARS := [9]string{"\\", "/", ":", "*", "?", "\"", "<", ">", "|"}
	const HELP_FLAG = "-h"
	const SWAP_FLAG = "-sw"
	const SET_DEFAULT_PATH_FLAG = "-d"
	const SET_DEFAULT_WINPATH_FLAG = "-dw"
	const SET_TGK_FOLDER_FLAG = "-tf"
	const SET_OLDDIR_NAME_FLAG = "-o"
	help := HelpInformation()
	// concatenate all args after 1 (including 1) into 1
	if len(args) > 1 {
		args[1], err = CombineString(args[1:])
		args = args[:2]
	}
	if err != nil {
		return err
	}
	if len(args) == 0 {
		return err
	}
	if help.availableFlagsWithDesc[args[0]] == "" {
		err = errors.New("Flag does not exist")
		return err
	}
	if args[0] == HELP_FLAG && len(args) == 2 {
		if help.availableFlagsWithDesc[args[1]] == "" {
			err = errors.New("Flag does not exist: " + args[1])
			return err
		}
		fmt.Printf("flag: %s\t%s\n", args[1], help.availableFlagsWithDesc[args[1]])
		return err
	} else if args[0] == HELP_FLAG && len(args) == 1 {
		help := HelpInformation()
		for k, v := range help.availableFlagsWithDesc {
			fmt.Printf("flag: %s\t%s\n", k, v)
		}
		return err
	}
	if args[0] == SWAP_FLAG && len(args) == 2 {
		// add checking for correct dir names here
		err = swapDirectories(getCompleteSettings(SETTINGS_FILE_NAME), args[1], SETTINGS_FILE_NAME)
		return err
	} else if args[0] == SWAP_FLAG && len(args) != 2 {
		err = errors.New("Not the correct number of arguments supplied for -sw flag (1).")
		return err
	}
	if len(args) > 2 {
		err = errors.New("No flag supports more than 2 arguments. At most run > fastSwapper -flag <argument for flag>")
		return err
	}
	// set default path flag expects the syntax of fastSwapper -d <path to directory>
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
		setSettings(SETTINGS_FILE_NAME, "Tgkdir", candidatePath)
	}
	if ContainsString(args, SET_DEFAULT_WINPATH_FLAG) {
		if len(args) > 1 {
			err = errors.New("Flag -dw does not take any additional arguments.")
			return err
		}
		setSettings(SETTINGS_FILE_NAME, "Tgkdir", TGK_DIR_DEFAULT_WIN)
		fmt.Printf("%s set as tagetik addin directory.\n", TGK_DIR_DEFAULT_WIN)
	}
	if ContainsString(args, SET_OLDDIR_NAME_FLAG) {
		if len(args) < 2 {
			err = errors.New("No name for the old directory provided. Use fastSwapper -o <name of the old directory>.")
			return err
		}
		candidateName := args[1]
		// we should probably also check for characters not supported in directory names..
		if ContainsStringWord(FORBIDDEN_CHARS[:], candidateName) {
			err = errors.New("Supplied name must not contain forbidden character.")
			return err
		}
		setActiveSettings(SETTINGS_FILE_NAME, "OldDirectory", candidateName)
		return err
	}
	if ContainsString(args, SET_TGK_FOLDER_FLAG) {
		if len(args) < 2 {
			err = errors.New("No name for the addin directory provided. Use fastSwapper -tf <name of the old directory>.")
			return err
		}
		candidateName := args[1]
		// we should probably also check for characters not supported in directory names..
		if ContainsStringWord(FORBIDDEN_CHARS[:], candidateName) {
			err = errors.New("Supplied name must not contain forbidden character.")
			return err
		}
		setSettings(SETTINGS_FILE_NAME, "Tgkfolder", candidateName)
		return err
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

func getCompleteSettings(filename string) Settings {
	return unmarshalSettingsJson(filename)
}

func getSettings(filename string) Default {
	return unmarshalSettingsJson(filename).Settings[0]
}

func getActiveSettings(filename string) ActiveSettings {
	return unmarshalSettingsJson(filename).ActiveSettings[0]
}

func setSettings(filename string, defaultToChange string, newValue string) {
	unmarshaledJson := unmarshalSettingsJson(filename)

	err := reflections.SetField(&unmarshaledJson.Settings[0], defaultToChange, newValue)
	if err != nil {
		log.Fatal(err)
	}
	updateSettingsJson(filename, unmarshaledJson)
}

func setActiveSettings(filename string, defaultToChange string, newValue string) {
	unmarshaledJson := unmarshalSettingsJson(filename)

	err := reflections.SetField(&unmarshaledJson.ActiveSettings[0], defaultToChange, newValue)
	if err != nil {
		log.Fatal(err)
	}
	updateSettingsJson(filename, unmarshaledJson)
}

func swapDirectories(set Settings, newDirName string, settingsFileName string) error {
	var err error
	oldDirName := set.ActiveSettings[0].OldDirectory
	tgkDir := set.Settings[0].Tgkdir
	tgkfolder := set.Settings[0].Tgkfolder
	// add logic here that does the following:
	// check: does newdir exist?
	newDirPath := tgkDir + "\\" + newDirName
	if !IsDir(newDirPath) {
		err = errors.New("Folder to swap in does not exist.")
		return err
	}
	tgkDirPath := tgkDir + "\\" + tgkfolder
	oldDirPath := tgkDir + "\\" + oldDirName
	// 1. rename tgk dir to olddir
	err = os.Rename(tgkDirPath, oldDirPath)
	if err != nil {
		return err
	}
	// 2. rename newDir folder to tgk dir
	err = os.Rename(newDirPath, tgkDirPath)
	if err != nil {
		return err
	}
	// 3. update oldDir setting with newDir
	setActiveSettings(settingsFileName, "OldDirectory", newDirName)
	return err
}
