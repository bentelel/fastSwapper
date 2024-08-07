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
	NewDirectory string `json:"newdirectory"`
}

type helpInformation struct {
	availableFlagsWithDesc map[string]string
}

func HelpInformation() helpInformation {
	var help helpInformation
	help.availableFlagsWithDesc = map[string]string{
		"-d":  "Set default tagetik directory > fastSwapper -d <absolute path to directory>",
		"-dw": "Reset default tagetik directory to the Windows one.",
		"-o":  "Set the name of the old directory, under this name the current Addin will be saved on swap. > fastSwapper -o <name of directory you want>",
		"-n":  "Set the name of the new directory, this will be used to remember the chosen name of the current Addin version, \n\t\tbecause we set that to the default directory and don't want the user to retype it every time. \n\t\t> fastSwapper -n <name of directory you want>",
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
	if len(args) == 0 {
		return err
	}
	const HELP_FLAG = "-h"
	if args[0] == HELP_FLAG && len(args) == 2 {
		help := HelpInformation()
		if help.availableFlagsWithDesc[args[1]] == "" {
			err = errors.New("Flag does not exist: " + args[1])
			return err
		}
		fmt.Printf("flag: %s\t%s\n", args[1], help.availableFlagsWithDesc[args[1]])
		return err
	}
	if args[0] == HELP_FLAG {
		help := HelpInformation()
		for k, v := range help.availableFlagsWithDesc {
			fmt.Printf("flag: %s\t%s\n", k, v)
			return err
		}
	}
	const SWAP_FLAG = "-sw"
	if args[0] == SWAP_FLAG && len(args) == 1 {
		// add checking for correct dir names here
		swapDirectories(getCompleteSettingis(SETTINGS_FILE_NAME))
		return err
	} else if args[0] == SWAP_FLAG && len(args) > 1 {
		err = errors.New("Flag -sw does not take any arguments.")
		return err
	}
	if len(args) > 2 {
		err = errors.New("No flag supports more than 2 arguments. At most run > fastSwapper -flag <argument for flag>")
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
		setSettings(SETTINGS_FILE_NAME, "Tgkdir", candidatePath)
	}
	const SET_DEFAULT_WINPATH_FLAG = "-dw"
	const TGK_DIR_DEFAULT_WIN = "C:\\Tagetik\\Tagetik Excel .NET Client"
	if ContainsString(args, SET_DEFAULT_WINPATH_FLAG) {
		if len(args) > 1 {
			err = errors.New("Flag -dw does not take any additional arguments.")
			return err
		}
		setSettings(SETTINGS_FILE_NAME, "Tgkdir", TGK_DIR_DEFAULT_WIN)
		fmt.Printf("%s set as tagetik addin directory.\n", TGK_DIR_DEFAULT_WIN)
	}
	// set old directory name flag expects the syntax of fastSwapper -o <name directory>
	const SET_OLDDIR_NAME_FLAG = "-o"
	FORBIDDEN_CHARS := [9]string{"\\", "/", ":", "*", "?", "\"", "<", ">", "|"}
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
	// set new directory name flag expects the syntax of fastSwapper -n <name directory>
	const SET_NEWDIR_NAME_FLAG = "-n"
	if ContainsString(args, SET_NEWDIR_NAME_FLAG) {
		if len(args) < 2 {
			err = errors.New("No name for the new directory provided. Use fastSwapper -n <name of the new directory>.")
			return err
		}
		candidateName := args[1]
		// we should probably also check for characters not supported in directory names..
		if ContainsStringWord(FORBIDDEN_CHARS[:], candidateName) {
			err = errors.New("Supplied name must not contain forbidden character.")
			return err
		}
		setActiveSettings(SETTINGS_FILE_NAME, "NewDirectory", candidateName)
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

func swapDirectories(set Settings) error {
	var err error
	oldDir := set.ActiveSettings[0].OldDirectory
	newDir := set.ActiveSettings[0].NewDirectory
	tgkDir := set.Settings[0].Tgkdir
	tgkfolder := set.Settings[0].Tgkfolder

	return err
}
