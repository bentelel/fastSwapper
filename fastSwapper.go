package main

import (
        "fmt"
        "os" 
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
  //args := os.Args[1:]
  //setDefaults(args)
  set := getDefaults("settings.json")
  fmt.Println(set.Tgkdir)

  setDefaults("settings.json", "Tgkdir", "C:\\Tagetik\\Tagetik Excel .NET Client")
}


func typeof(v interface{}) string {
  return reflect.TypeOf(v).String()
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

func updateSettingsJson(filename string, data Default) {
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
  unmarshaledJson := unmarshalSettingsJson(filename).Defaults[0]

  err := reflections.SetField(&unmarshaledJson, defaultToChange, newValue)
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


