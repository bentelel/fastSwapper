package main

import (
        "fmt"
        "os" 
        "errors"
        "path/filepath"
        "log"
        "encoding/json"
        "io/ioutil"
        "reflect"
        "github.com/oleiade/reflections"
        )

type Defaults struct {
  Defaults []Default `json:"defaults"` 
}
type Default struct {
  Tgkdir string `json:"tgkdir"`
}
type ActiveSettings struct {
  OldDirectory string
  NewDirectory string
}


func main() {
  args := os.Args[1:]
  err := parseCLIargs(args)
  if err != nil {
    log.Fatal(err.Error())
  }
  set := getDefaults("settings.json")

  // setDefaults("settings.json", "Tgkdir", "C:\\Tagetik\\Tagetik Excel .NET Client")
  
  fmt.Println(getAllInDir(set.Tgkdir))
}


func typeof(v interface{}) string {
  return reflect.TypeOf(v).String()
}

func containsString(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

func parseCLIargs (args []string) error {
  var err error
  if len(args)==0 {
    return err
  }
  //set default path flag expects the syntax of fastSwapper -d <path to directory>
  const SET_DEFAULT_PATH_FLAG = "-d"
  if containsString(args, SET_DEFAULT_PATH_FLAG){
    if len(args)<2 {
      err = errors.New("No path provided. Use fastSwapper -d <path to default dir>.")
      return err
    }
    candidatePath := args[1]
    if !isDir(candidatePath) {
      err = errors.New("Supplied path does not exist.")
      return err
    }
    setDefaults("settings.json", "Tgkdir", candidatePath)
  }
  
  return err
}


func isDir(path string) bool {
  if _, err := os.Stat(path); os.IsNotExist(err) {
    return false    
  }
  return true
} 

func unmarshalSettingsJson(filename string) Defaults {
  jsonFile, err := os.Open(filename)
  if err != nil{
    log.Fatal(err)
  }
  defer jsonFile.Close()  
  byteResult, err := ioutil.ReadAll(jsonFile)
  if err != nil{
    log.Fatal(err)
  }
  var settings Defaults
  err = json.Unmarshal(byteResult, &settings)
  if err != nil{
    log.Fatal(err)
  }
  return settings
}

func updateSettingsJson(filename string, data Defaults) {
  modifiedJson, err := json.MarshalIndent(data, "", "    ")
  if err != nil{
    log.Fatal(err)
  }
  err = ioutil.WriteFile(filename, modifiedJson, 0644)
  if err != nil{
    log.Fatal(err)
  }
}

func getDefaults(filename string) Default {
  return unmarshalSettingsJson(filename).Defaults[0] 
}


func setDefaults(filename string, defaultToChange string, newValue string) {
  unmarshaledJson := unmarshalSettingsJson(filename)//.Defaults[0]

  err := reflections.SetField(&unmarshaledJson.Defaults[0], defaultToChange, newValue)
  if err != nil{
      log.Fatal(err)
    }
  updateSettingsJson(filename, unmarshaledJson)
}


func getDirsInDir(dir string) []string{
  // Returns slice of strings containing all directories within given directory
  // param dir: string -- directory to check
  entries, err := os.ReadDir(filepath.FromSlash(dir))
  if err != nil{
    log.Fatal(err)
  }
  result := make([]string, 0) 
  //translate all vfiles to strings
  for _, e := range entries{
    if e.IsDir(){  
      result = append(result,e.Name())
    } 
  }
  return result
}

func getAllInDir(dir string) []string{
  // Returns slice of strings containing all directories within given directory
  // param dir: string -- directory to check
  entries, err := os.ReadDir(filepath.FromSlash(dir))
  if err != nil{
    log.Fatal(err)
 }
  result := make([]string, 0) 
  //translate all files to strings
  for _, e := range entries{
    result = append(result,e.Name())
  }
  return result
}


