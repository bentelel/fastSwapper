package main

import (
        "fmt"
        "os" 
        "path/filepath"
        "log"
        )
        

func main() {
  testMyFunctions()
}


func testMyFunctions() {
 // dirs := getDirsInDir(`C:\go_testing\testFolder`)
  // allFiles := getAllInDir(`C:\go_testing\testFolder`)
 
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
  fmt.Println(len(entries))
  if err != nil{
    log.Fatal(err)
    fmt.Println("error occured")
  }
  result := make([]string, 0) 
  //translate all files to strings
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
    fmt.Println("error occured")
 }
  result := make([]string, 0) 
  //translate all files to strings
  for _, e := range entries{
    result = append(result,e.Name())
  }
  return result
}


