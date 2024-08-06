package main

import (
        "fmt"
        "os" 
        "path/filepath"
        "log"
        "encoding/json"
        "io/ioutil"
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
  // testMyFunctions()
}


func unmarshalSettinsJson(filename string) Defaults {
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

func getDefaults(filename string) Default {
  return unmarshalSettinsJson(filename).Defaults[0] 
}


func testMyFunctions() {
 // dirs := getDirsInDir(`C:\go_testing\testFolder`)
  // allFiles := getAllInDir(`C:\go_testing\testFolder`)
  fmt.Println("Inside testMyFunctions") 
  dirs := getDirsInDir(`C:\go_testing`)
  allFiles := getAllInDir(`C:\go_testing`)
  
  fmt.Println("Printing all dirs:")
  for _, e := range dirs{
    fmt.Println(e)
  }
  fmt.Println("")
  fmt.Println("Printing all dirs and files:")
  for _, e := range allFiles{
    fmt.Println(e)
  }
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


